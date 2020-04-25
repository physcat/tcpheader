package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/physcat/tcpheader"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "broker",
	Short: "Broker - a test client for tcp connection using headers",
	Long: `Broker - a test client for tcp connection using headers

TCP sockets often use headers indicating the length of the message.
This sample program is useful for testing such connection as either a
server or a client application.

Currently supported are the following headers:
- Two Byte unsigned length indicator
- Four Byte unsigned length indicator
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.broker.yaml)")
	rootCmd.PersistentFlags().String("port", "8080", "port to use")

	rootCmd.PersistentFlags().Int("header", 2, "Header length")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		viper.AddConfigPath(home)
		viper.SetConfigName(".broker")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func GetHeader(cmd *cobra.Command) tcpheader.HeaderType {
	h, _ := cmd.Flags().GetInt("header")

	switch h {
	case 2:
		return tcpheader.TwoByteUnsigned
	case 4:
		return tcpheader.FourByteUnsigned
	default:
		return tcpheader.Unknown
	}
}

func PrintHeaderError(cmd *cobra.Command) {
	h, _ := cmd.Flags().GetInt("header")
	fmt.Printf(`Unknown header type (%d)

Currently only 2 and 4 are supported.
`, h)
}

func ReadAndWrite(r io.ReadWriter, listen bool, msg string, header tcpheader.HeaderType) {
	switch listen {
	case true:
		// Read response
		l, err := tcpheader.ReadLen(r, header)
		if err != nil {
			fmt.Printf("Error reading form conn: %+v\n", err)
			return
		}

		buf := make([]byte, l)
		if err = tcpheader.ReadMessage(r, buf); err != nil {
			fmt.Printf("Error reading form conn: %+v\n", err)
			return
		}

		fmt.Printf("Received message: %s\n", buf)

		//Send Message
		fmt.Printf("Sending message: %s\n", msg)

		if err := tcpheader.WriteMessage(r, []byte(msg), header); err != nil {
			fmt.Printf("Failed to write: %+v", err)
			return
		}

	case false:
		//Send Message
		fmt.Printf("Sending message: %s\n", msg)

		if err := tcpheader.WriteMessage(r, []byte(msg), header); err != nil {
			fmt.Printf("Failed to write: %+v", err)
			return
		}

		// Read response
		l, err := tcpheader.ReadLen(r, header)
		if err != nil {
			fmt.Printf("Error reading form conn: %+v\n", err)
			return
		}

		buf := make([]byte, l)
		if err = tcpheader.ReadMessage(r, buf); err != nil {
			fmt.Printf("Error reading form conn: %+v\n", err)
			return
		}

		fmt.Printf("Received message: %s\n", buf)
	}
}
