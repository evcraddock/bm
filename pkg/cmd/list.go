package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

// NewCmdList creates new list command
func NewCmdList(out io.Writer) *cobra.Command {
	o := NewBaseOptions(out)
	cmd := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd, args)
			o.listBookmarks()
		},
	}

	cmd.Flags().String("category", "", "category folder")

	return cmd
}

func (o *BaseOptions) listBookmarks() {
	manager := bookmarks.NewBookmarkManager(o.Config, false, o.Category)

	bookmarks, err := manager.LoadBookmarks()
	if err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
		os.Exit(1)
	}

	b, _ := yaml.Marshal(bookmarks)
	if _, err := o.Out.Write(b); err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
		os.Exit(1)
	}
}
