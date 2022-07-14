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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "create",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete",
}

func init() {
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
