package remove

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/config"
)

type RemoveOptions struct {
	Out      io.Writer
	Config   *config.Config
	Category string
	Title    string
}

func NewRemoveOptions(out io.Writer) *RemoveOptions {
	return &RemoveOptions{
		Out: out,
	}
}

func NewCmdRemove(out io.Writer) *cobra.Command {
	o := NewRemoveOptions(out)
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

func (o *RemoveOptions) prepare(cmd *cobra.Command, args []string) {
	if category, err := cmd.Flags().GetString("category"); err == nil {
		o.Category = category
	}

	cfg, err := config.LoadConfigFile()
	if err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
		os.Exit(1)
	}

	o.Config = cfg
	o.Title = args[0]
}

func (o *RemoveOptions) removeBookmark() {
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
