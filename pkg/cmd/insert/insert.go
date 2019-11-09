package insert

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/config"
)

type InsertOptions struct {
	Out         io.Writer
	Config      *config.Config
	Interactive bool
	Category    string
	Title       string
	URL         string
}

func NewInsertOptions(out io.Writer) *InsertOptions {
	return &InsertOptions{
		Out: out,
	}
}

func NewCmdInsert(out io.Writer) *cobra.Command {
	o := NewInsertOptions(out)
	cmd := &cobra.Command{
		Use: "insert",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd, args)
			o.validate()
			o.addBookmark()
		},
	}

	cmd.Flags().String("interactive", "false", "interactive mode")
	cmd.Flags().String("category", "readlater", "category folder")

	return cmd
}

func (o *InsertOptions) prepare(cmd *cobra.Command, args []string) {
	o.Interactive = false
	if interactive, err := cmd.Flags().GetString("interactive"); err == nil {
		boolval, err := strconv.ParseBool(interactive)
		if err == nil {
			o.Interactive = boolval
		}
	}

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
	o.URL = args[1]
}

func (o *InsertOptions) validate() {
	//TODO: validate inputs
}

func (o *InsertOptions) addBookmark() {
	bookmark := bookmarks.Bookmark{
		Title: o.Title,
		URL:   o.URL,
	}

	manager := bookmarks.NewBookmarkManager(o.Config, o.Interactive)
	err := manager.Save(bookmark, o.Category)

	if err != nil {
		fmt.Fprintf(o.Out, "%s", err.Error())
	}
}
