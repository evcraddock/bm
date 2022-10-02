package commands

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/evcraddock/bm/internal/config"
	"github.com/evcraddock/bm/internal/tui"
)

var (
	category string
)

var rootCmd = &cobra.Command{
	Use:   "bm",
	Short: "bm is a command line tool for managing bookmarks",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		config.LoadConfig()
		viper.BindPFlag("category", cmd.PersistentFlags().Lookup("category"))
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if category == "" {
			category = viper.GetString("DefaultCategory")
		}

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
	rootCmd.PersistentFlags().StringVarP(&category, "category", "c", "", "category")
	// viper.BindPFlag("category", rootCmd.PersistentFlags().Lookup("category"))

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
