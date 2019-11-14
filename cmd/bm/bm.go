package main

import (
	"fmt"
	"os"

	"github.com/evcraddock/bm/pkg/cmd"
	"github.com/evcraddock/bm/pkg/config"
)

func main() {
	if ok := config.FileExists(); !ok {
		err := config.CreateDefaultConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}

	command := cmd.NewDefaultCommand()
	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
