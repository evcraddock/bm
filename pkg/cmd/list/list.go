package list

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/config"
)

type ListOptions struct {
	Out      io.Writer
	Config   *config.Config
	Category string
}

func NewListOptions(out io.Writer) *ListOptions {
	return &ListOptions{
		Out: out,
	}
}

func NewCmdList(out io.Writer) *cobra.Command {
	o := NewListOptions(out)
	cmd := &cobra.Command{
		Use: "list",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd)
			o.listBookmarks()
		},
	}

	cmd.Flags().String("category", "", "category folder")

	return cmd
}

func (o *ListOptions) prepare(cmd *cobra.Command) {
	if category, err := cmd.Flags().GetString("category"); err == nil {
		o.Category = category
	}

	cfg, err := config.LoadConfigFile()
	if err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
		os.Exit(1)
	}

	o.Config = cfg

}

func (o *ListOptions) listBookmarks() {
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
