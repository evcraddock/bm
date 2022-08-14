package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/cmd/bm/tui"
)

var (
	category = "readlater"
)

var rootCmd = &cobra.Command{
	Use:   "bm",
	Short: "bm",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			p := tea.NewProgram(tui.New(category))
			if err := p.Start(); err != nil {
				fmt.Println("unable to start app")
				os.Exit(1)
			}
		}
	},
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

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update",
}

func init() {
	rootCmd.Flags().StringVarP(&category, "category", "c", "readlater", "category")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(updateCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
