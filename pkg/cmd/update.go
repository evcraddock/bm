package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

// NewCmdUpdate creates new update command
func NewCmdUpdate(out io.Writer) *cobra.Command {
	o := NewBaseOptions(out)
	cmd := &cobra.Command{
		Use: "update",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd, args)
			o.updateBookmark()
		},
	}

	cmd.Flags().String("category", "readlater", "category folder")

	return cmd
}

func (o *BaseOptions) updateBookmark() {
	manager := bookmarks.NewBookmarkManager(o.Config, true, o.Category)
	err := manager.Update(o.Title)
	if err != nil {
		fmt.Fprintf(o.Out, "%s\n", err.Error())
	}
}
