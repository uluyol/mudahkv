package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"golang.org/x/net/context"

	"github.com/uluyol/mudahkv/client"
)

const timeout = 10 * time.Second

func main() {
	if len(os.Args) < 4 {
		fmt.Fprintf(os.Stderr, "%s mudahAddr outputkey command to run with flags\n", os.Args[0])
		os.Exit(1)
	}
	failureOccured := false
	addr := os.Args[1]
	outputKey := os.Args[2]
	cmdName := os.Args[3]
	cmdArgs := os.Args[4:]

	cmd := exec.Command(cmdName, cmdArgs...)
	out, outErr := cmd.CombinedOutput()
	errString := "success"
	if outErr != nil {
		errString = fmt.Sprintf("failure: %v", outErr)
	}

	c, err := client.Dial(addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect to server: %v\n", err)
		os.Exit(1)
	}

	// try to send whatever we can to the DB, fail at the end
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	if err := c.Set(ctx, outputKey, string(out)); err != nil {
		fmt.Fprintf(os.Stderr, "error sending output: %v", err)
		failureOccured = true
	}
	ctx, _ = context.WithTimeout(context.Background(), timeout)
	if err := c.Set(ctx, outputKey+"-error", errString); err != nil {
		fmt.Fprintf(os.Stderr, "error sending command exit code: %v", err)
		failureOccured = true
	}
	c.Close()
	if failureOccured {
		os.Exit(1)
	}
}
