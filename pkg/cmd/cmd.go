package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/app"
	"github.com/evcraddock/bm/pkg/bookmarks"
	"github.com/evcraddock/bm/pkg/config"
)

// NewDefaultCommand creates a new cobra.Command
func NewDefaultCommand() *cobra.Command {
	return NewBookmarkCommand(os.Stdin, os.Stdout, os.Stderr)
}

// NewBookmarkCommand creates a new bookmark command
func NewBookmarkCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	o := NewBaseOptions(out)
	cmds := &cobra.Command{
		Use:   "bm",
		Short: "bm manages bookmark lists",
		Long:  "bm manages bookmark lists",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd, args)
			o.startApp()
		},
	}

	cmds.Flags().String("category", "readlater", "category folder")

	cmds.AddCommand(NewCmdInsert(out))
	cmds.AddCommand(NewCmdRemove(out))
	cmds.AddCommand(NewCmdUpdate(out))
	cmds.AddCommand(NewCmdList(out))
	cmds.AddCommand(NewCmdShow(out))

	return cmds
}

// BaseOptions default options for a command
type BaseOptions struct {
	Out         io.Writer
	Config      *config.Config
	Category    string
	Interactive bool
	Title       string
	URL         string
}

// NewBaseOptions creates a new baseOptions
func NewBaseOptions(out io.Writer) *BaseOptions {
	return &BaseOptions{
		Out: out,
	}
}

func (o *BaseOptions) prepare(cmd *cobra.Command, args []string) {
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

	if len(args) > 0 {
		o.Title = args[0]
	}

	if len(args) > 1 {
		o.URL = args[1]
	}
}

func (o *BaseOptions) startApp() {
	manager := bookmarks.NewBookmarkManager(o.Config, false, o.Category)

	bmApp := app.NewBookmarkApp(manager)
	bmApp.Load()
}
