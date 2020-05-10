package cmd

import (
	"fmt"
	"net"
	"os"

	"github.com/physcat/tcpheader"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Listener, wait for a call read and respond",
	Long: `Listener, wait for a call read and respond
	
The server will bind to a given port and wait for a message, 
display the message and send a test response.`,
	Run: serverMain,
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().String("ip", "", "optional network interface to connect to")
	serverCmd.PersistentFlags().String("message", "Got it!", "Response message to send to client")
	serverCmd.PersistentFlags().Bool("listen", true, "listen first and then respond with the message")
}

func serverMain(cmd *cobra.Command, args []string) {
	ip, _ := cmd.Flags().GetString("ip")
	port, _ := cmd.Flags().GetString("port")

	connStr := fmt.Sprintf("%s:%s", ip, port)

	ln, err := net.Listen("tcp", connStr)
	if err != nil {
		fmt.Println("Failed to listen on port 8080")
		return
	}

	fmt.Printf("Listening on: %s\n", ln.Addr())

	conn, err := ln.Accept()
	if err != nil {
		fmt.Printf("Failed to accept: %+v\n", err)
		return
	}

	handleServerConnection(conn, cmd)
}

func handleServerConnection(conn net.Conn, cmd *cobra.Command) {
	defer conn.Close()
	fmt.Printf("Got connection from: %s\n", conn.RemoteAddr())

	header := GetHeader(cmd)
	if header == tcpheader.Unknown {
		PrintHeaderError(cmd)
		return
	}
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
