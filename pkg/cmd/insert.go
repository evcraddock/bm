package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

func NewCmdInsert(out io.Writer) *cobra.Command {
	o := NewBaseOptions(out)
	cmd := &cobra.Command{
		Use: "insert",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd, args)
			o.addBookmark()
		},
	}

	cmd.Flags().String("interactive", "false", "interactive mode")
	cmd.Flags().String("category", "readlater", "category folder")

	return cmd
}

func (o *baseOptions) addBookmark() {
	manager := bookmarks.NewBookmarkManager(o.Config, o.Interactive, o.Category)
	err := manager.Create(o.Title, o.URL)
	if err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
	}
}
