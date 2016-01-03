package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var serverAddr string
var reqTimeout time.Duration

var RootCmd = &cobra.Command{
	Use:   "mudahc",
	Short: "A client for a MudahKV server",
	Long: `This is a simple client application that can be used
to set and retrieve values from a MudahKV server.`,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func dief(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
	os.Exit(1)
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports Persistent Flags, which, if defined here,
	// will be global for your application.

	RootCmd.PersistentFlags().StringVarP(&serverAddr, "addr", "a", "localhost:6070", "address of server")
	RootCmd.PersistentFlags().DurationVarP(&reqTimeout, "timeout", "t", 3*time.Second, "request timeout length")
}
