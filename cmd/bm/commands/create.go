package commands

import (
	"fmt"
	"net/url"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

var createBookmarkCmd = &cobra.Command{
	Use:   "bookmark",
	Short: "create bookmark",
	Run:   cmdCreateBookmark,
}

func init() {
	createCmd.AddCommand(createBookmarkCmd)

	createBookmarkCmd.Flags().StringVarP(&category, "category", "c", "readlater", "category")
}

func cmdCreateBookmark(cmd *cobra.Command, args []string) {
	var title string
	if len(args) < 1 {
		fmt.Println("bm create bookmark {url} ({title})")
		os.Exit(1)
	}

	link := args[0]
	if len(args) == 2 {
		title = args[1]
	}

	_, err := url.ParseRequestURI(link)
	if err != nil {
		fmt.Println("invalid url")
		os.Exit(1)
	}

	if category == "" {
		category = viper.GetString("DefaultCategory")
	}

	manager := bookmarks.NewBookmarkManager(false, category)
	if ok, err := manager.Create(link, title); ok {
		fmt.Printf("Saved Bookmark %s\n", title)
	} else if err != nil {
		// TODO: add a proper log and log error in debug mode
		fmt.Println("unable to create bookmark")
		os.Exit(1)
	}
}
