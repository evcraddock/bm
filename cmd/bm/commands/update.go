package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var updateBookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "update bookmark",
	Run:   cmdUpdateBookmark,
}

func init() {
	updateCmd.AddCommand(updateBookmarkCmd)
}

func cmdUpdateBookmark(cmd *cobra.Command, args []string) {
	// var title string

	// if len(args) > 0 {
	// 	title = args[0]
	// }

	// manager := bookmarks.NewBookmarkManager(true, category)
	// _, err := manager.Update(title)
	// if err != nil {
	// 	fmt.Println("unable to update bookmark")
	// 	os.Exit(1)
	// }

	fmt.Println("not implemented")
}
