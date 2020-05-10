package cmd

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/physcat/tcpheader"
	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Dialer, place a call and send a test message",
	Long: `Dialer, place a call and send a test message
	
The client will dial a given address and attempt to send a test message.
	`,
	Run: clientMain,
}

func init() {
	rootCmd.AddCommand(clientCmd)

	clientCmd.PersistentFlags().String("host", "localhost", "Host to dial")
	clientCmd.PersistentFlags().String("message", "Test message sent from client", "Message to send")
	clientCmd.PersistentFlags().Bool("listen", false, "listen first and then respond with the message")
}

func clientMain(cmd *cobra.Command, args []string) {
	header := GetHeader(cmd)
	if header == tcpheader.Unknown {
		PrintHeaderError(cmd)
		return
	}

	var d net.Dialer

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	conn, err := d.DialContext(ctx, "tcp", "localhost:8080")
	if err != nil {
		fmt.Printf("Failed to dial: %+v", err)
		return
	}

	defer conn.Close()

	echo, _ := cmd.Flags().GetBool("echo")

	netMsg := ReadWithHeaderC(conn, header)
	stdinMsg := ReadPlainC(os.Stdin)

	for {
		select {
		case msg, ok := <-netMsg:
			fmt.Printf("net>%q\n", msg)
			if !ok {
				return
			}

		case msg, ok := <-stdinMsg:
			if !ok {
				return
			}
			if echo {
				fmt.Printf("stdio>%q\n", msg)
			}
			if err := tcpheader.WriteMessage(conn, []byte(msg), header); err != nil {
				return
			}
		}
	}
}
