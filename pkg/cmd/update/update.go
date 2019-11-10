package update

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/config"
)

type UpdateOptions struct {
	Out      io.Writer
	Config   *config.Config
	Category string
	Title    string
}

func NewUpdateOptions(out io.Writer) *UpdateOptions {
	return &UpdateOptions{
		Out: out,
	}
}

func NewCmdUpdate(out io.Writer) *cobra.Command {
	o := NewUpdateOptions(out)
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

func (o *UpdateOptions) prepare(cmd *cobra.Command, args []string) {
	if category, err := cmd.Flags().GetString("category"); err == nil {
		o.Category = category
	}

	cfg, err := config.LoadConfigFile()
	if err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
		os.Exit(1)
	}

	if len(args) < 1 {
		fmt.Fprintf(o.Out, "wrong number of arguments\n")
		os.Exit(1)
	}

	o.Config = cfg
	o.Title = args[0]
}

func (o *UpdateOptions) updateBookmark() {
	manager := bookmarks.NewBookmarkManager(o.Config, true, o.Category)

	err := manager.Update(o.Title)
	if err != nil {
		fmt.Fprintf(o.Out, "%s\n", err.Error())
	}
}
