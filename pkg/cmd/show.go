package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/evcraddock/bm/pkg/bookmarks"
)

// NewCmdShow creates new show command
func NewCmdShow(out io.Writer) *cobra.Command {
	o := NewBaseOptions(out)
	cmd := &cobra.Command{
		Use: "show",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd, args)
			o.showBookmark()
		},
	}

	cmd.Flags().String("category", "readlater", "category folder")

	return cmd
}

func (o *BaseOptions) showBookmark() {
	manager := bookmarks.NewBookmarkManager(o.Config, false, o.Category)
	bookmarkLocation := manager.GetBookmarkLocation(o.Title)
	bookmark, err := manager.Load(bookmarkLocation)
	if err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
		os.Exit(1)
	}

	b, _ := yaml.Marshal(bookmark)
	if _, err := o.Out.Write(b); err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
		os.Exit(1)
	}
}
