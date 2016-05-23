package cmd

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/net/context"

	"github.com/spf13/cobra"
	"github.com/uluyol/mudahkv/client"
)

var filename string

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a value",
	Long: `Set sets the value for the given key.

Example Usage: mudahc set key value`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := doSet(args); err != nil {
			dief("%v", err)
		}
	},
}

func init() {
	RootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&filename, "filename", "f", "", "path to file to load data from")
}

func doSet(args []string) error {
	if len(args) != 2 && filename == "" {
		return errors.New("need exactly 2 args or filename")
	}

	c, err := client.Dial(serverAddr)
	if err != nil {
		return fmt.Errorf("unable to connect to server: %v", err)
	}
	defer c.Close()

	if len(args) == 2 {
		ctx, _ := context.WithTimeout(context.Background(), reqTimeout)
		if err := c.Set(ctx, args[0], []byte(args[1])); err != nil {
			return fmt.Errorf("unable to set value: %v", err)
		}
		return nil
	}

	f, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("unable to open file %s: %v", filename, err)
	}
	defer f.Close()

	ctx, _ := context.WithTimeout(context.Background(), reqTimeout)
	if err := c.SetStream(ctx, args[0], f); err != nil {
		return fmt.Errorf("unable to set value: %v", err)
	}
	return nil
}
