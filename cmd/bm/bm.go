package main

import (
	"fmt"
	"os"

	"github.com/evcraddock/bm/cmd/bm/commands"
	"github.com/evcraddock/bm/internal/config"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	commands.Execute()
}
