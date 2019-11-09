package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/evcraddock/bm/pkg/cmd/insert"
	"github.com/evcraddock/bm/pkg/cmd/remove"
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

	cmds.AddCommand(insert.NewCmdInsert(out))
	cmds.AddCommand(remove.NewCmdRemove(out))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
