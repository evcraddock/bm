package cmd

import (
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/config"
)

func NewDefaultCommand() *cobra.Command {
	return NewBookmarkCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewBookmarkCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "bm",
		Short: "bm manages bookmark lists",
		Long:  "bm manages bookmark lists",
		Run:   runHelp,
	}

	cmds.AddCommand(NewCmdInsert(out))
	cmds.AddCommand(NewCmdRemove(out))
	cmds.AddCommand(NewCmdUpdate(out))
	cmds.AddCommand(NewCmdList(out))
	cmds.AddCommand(NewCmdShow(out))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}

type baseOptions struct {
	Out         io.Writer
	Config      *config.Config
	Category    string
	Interactive bool
	Title       string
	URL         string
}

func NewBaseOptions(out io.Writer) *baseOptions {
	return &baseOptions{
		Out: out,
	}
}

func (o *baseOptions) prepare(cmd *cobra.Command, args []string) {
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
