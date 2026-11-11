package main

import (
	"fmt"
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "dcmd:", err)
		os.Exit(1)
	}
}
