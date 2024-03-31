package main

import (
	"fmt"
	"os"

	"github.com/sol1du2/rdetective/cmd"
	"github.com/sol1du2/rdetective/cmd/rdetective/diff"
)

func main() {
	cmd.RootCmd.Use = "rdetective"

	cmd.RootCmd.AddCommand(cmd.CommandVersion())
	cmd.RootCmd.AddCommand(diff.CommandDiff())

	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
