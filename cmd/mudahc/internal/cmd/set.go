package cmd

import (
	"fmt"
	"io/ioutil"
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
		if len(args) != 2 && filename == "" {
			dief("need exactly 2 arguments or filename, got %d args", len(args))
		}
		var data string
		if len(args) == 2 {
			data = args[1]
		} else {
			f, err := os.Open(filename)
			if err != nil {
				dief("unable to open file %s: %v", filename, err)
			}
			defer f.Close()
			if b, err := ioutil.ReadAll(f); err != nil {
				dief("error while reading file: %v", err)
			} else {
				data = string(b)
			}
		}
		c, err := client.Dial(serverAddr)
		if err != nil {
			dief("unable to connect to server: %v", err)
		}
		ctx, _ := context.WithTimeout(context.Background(), reqTimeout)
		err = c.Set(ctx, args[0], data)
		if err != nil {
			c.Close()
			dief("unable to set value: %v", err)
		}
		fmt.Println(data)
		c.Close()
	},
}

func init() {
	RootCmd.AddCommand(setCmd)

	setCmd.Flags().StringVarP(&filename, "filename", "f", "", "path to file to load data from")
}
