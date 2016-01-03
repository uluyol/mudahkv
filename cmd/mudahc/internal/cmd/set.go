package cmd

import (
	"fmt"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
	"github.com/uluyol/mudahkv/client"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a value",
	Long: `Set sets the value for the given key.

Example Usage: mudahc set key value`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			dief("need exactly 2 arguments, got %d", len(args))
		}
		c, err := client.Dial(serverAddr)
		if err != nil {
			dief("unable to connect to server: %v", err)
		}
		ctx, _ := context.WithTimeout(context.Background(), reqTimeout)
		err = c.Set(ctx, args[0], args[1])
		if err != nil {
			c.Close()
			dief("unable to set value: %v", err)
		}
		fmt.Println(args[1])
		c.Close()
	},
}

func init() {
	RootCmd.AddCommand(setCmd)
}
