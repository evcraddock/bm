package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

var deleteBookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "delete bookmark",
	Run:   cmdDeleteBookmark,
}

func init() {
	deleteCmd.AddCommand(deleteBookmarkCmd)
}

func cmdDeleteBookmark(cmd *cobra.Command, args []string) {
	if len(args) != 1 {
		fmt.Println("bm delete bookmark {title}")
		os.Exit(1)
	}

	title := args[0]

	manager := bookmarks.NewBookmarkManager(false, category)
	if err := manager.Remove(title); err != nil {
		fmt.Println("unable to delete bookmark")
		os.Exit(1)
	}

	fmt.Printf("deleted bookmark %s\n", title)
	os.Exit(1)
}
