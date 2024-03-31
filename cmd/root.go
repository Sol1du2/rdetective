package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// RootCmd provides the command line parser root.
var RootCmd = &cobra.Command{
	Run: func(cmd *cobra.Command, _ []string) {
		_ = cmd.Help()
		os.Exit(2)
	},
}
