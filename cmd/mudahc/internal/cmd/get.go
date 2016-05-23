package cmd

import (
	"fmt"
	"io"
	"os"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
	"github.com/uluyol/mudahkv/lib/client"
)

var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get a value",
	Long: `Get returns the value for the given key.

Example Usage: mudahc get key`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			dief("need exactly 1 argument, got %d", len(args))
		}
		c, err := client.Dial(serverAddr)
		if err != nil {
			dief("unable to connect to server: %v", err)
		}
		ctx, _ := context.WithTimeout(context.Background(), reqTimeout)
		r, err := c.GetStream(ctx, args[0])
		if err != nil {
			c.Close()
			dief("unable to get value: %v", err)
		}
		_, err = io.Copy(os.Stdout, r)
		fmt.Println()
		if err != nil {
			dief("unable to get full value: %v", err)
		}
		c.Close()
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
