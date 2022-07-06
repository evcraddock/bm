package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	category = "readlater"
)

var rootCmd = &cobra.Command{
	Use:   "bm",
	Short: "bm",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list",
}

func init() {
	rootCmd.AddCommand(listCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
