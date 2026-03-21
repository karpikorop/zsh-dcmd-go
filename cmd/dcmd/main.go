package main

import (
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/cli"
)

func main() {
	if err := cli.RootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
