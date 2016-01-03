package cmd

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
	"github.com/uluyol/mudahkv/client"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists the values that share a prefix",
	Long: `List returns the values that have the given prefix.

Example Usage: mudahc list prefix`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			dief("need exactly 1 argument, got %d", len(args))
		}
		c, err := client.Dial(serverAddr)
		if err != nil {
			dief("unable to connect to server: %v", err)
		}
		ctx, _ := context.WithTimeout(context.Background(), reqTimeout)
		kvs, err := c.List(ctx, args[0])
		if err != nil {
			c.Close()
			dief("unable to list values: %v", err)
		}
		for _, kv := range kvs {
			fmt.Printf("%s:\n", kv.Key)
			fmt.Println(kv.Value)
		}
		c.Close()
	},
}

func init() {
	RootCmd.AddCommand(listCmd)
}
