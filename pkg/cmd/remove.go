package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

// NewCmdRemove creates new remove command
func NewCmdRemove(out io.Writer) *cobra.Command {
	o := NewBaseOptions(out)
	cmd := &cobra.Command{
		Use: "remove",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd, args)
			o.removeBookmark()
		},
	}

	cmd.Flags().String("category", "readlater", "category folder")

	return cmd
}

func (o *BaseOptions) removeBookmark() {
	manager := bookmarks.NewBookmarkManager(o.Config, false, o.Category)
	err := manager.Remove(o.Title)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(o.Out, "bookmark %s/%s does not exist \n", o.Category, o.Title)
			return
		}

		fmt.Fprintf(o.Out, "unable to remove bookmark %s/%s \n", o.Category, o.Title)
		return
	}

	fmt.Fprintf(o.Out, "Bookmark %s/%s has been removed \n", o.Category, o.Title)
}
