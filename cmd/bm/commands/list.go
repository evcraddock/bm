package commands

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/categories"
)

var listBookmarksCmd = &cobra.Command{
	Use:   "bookmarks",
	Short: "list bookmarks",
	Run:   cmdListBookmarks,
}

var listCategoriesCmd = &cobra.Command{
	Use:   "categories",
	Short: "list categories",
	Run:   cmdListCategories,
}

func init() {
	listCmd.AddCommand(listBookmarksCmd)
	listCmd.AddCommand(listCategoriesCmd)

	listBookmarksCmd.Flags().StringVarP(&category, "category", "c", "readlater", "category")
}

func cmdListCategories(cmd *cobra.Command, args []string) {
	manager := categories.NewCategoryManager()
	clist, err := manager.GetCategoryList()
	if err != nil {
		fmt.Println("there was an error listing categories")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(clist) == 0 {
		fmt.Println("no categories were found")
		os.Exit(1)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	header := fmt.Sprintf("Category\tDescription")
	fmt.Fprintln(writer, header)
	for _, c := range clist {
		row := fmt.Sprintf("%s\t%s", c.Name, c.Description)
		fmt.Fprintln(writer, row)
	}

	writer.Flush()
}

func cmdListBookmarks(cmd *cobra.Command, args []string) {
	manager := bookmarks.NewBookmarkManager(false, category)
	bookmarks, err := manager.LoadBookmarks()
	if err != nil {
		fmt.Println("there was an error listing bookmarks")
		fmt.Println(err)
		os.Exit(1)
	}

	if len(bookmarks) == 0 {
		fmt.Println("no bookmarks were found")
		os.Exit(1)
	}

	writer := tabwriter.NewWriter(os.Stdout, 0, 8, 1, '\t', tabwriter.AlignRight)
	header := fmt.Sprintf("Title\tLink")
	fmt.Fprintln(writer, header)
	for _, bookmark := range bookmarks {
		row := fmt.Sprintf("%s\t%s", bookmark.Title, bookmark.URL)
		fmt.Fprintln(writer, row)
	}

	writer.Flush()
}
