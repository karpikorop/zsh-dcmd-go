package cli

import (
	"fmt"
	"io"
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
)

// info prints a message for the user. It never goes to stdout, because stdout carries the command.
func info(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func printCommands(w io.Writer, p *store.Project) {
	if len(p.Commands) == 0 {
		fmt.Fprintf(w, "No commands found for '%s'.\n", p.Name)
		return
	}

	fmt.Fprintf(w, "Available commands for '%s':\n", p.Name)
	for _, c := range p.Commands {
		fmt.Fprintf(w, "  %d. %s\n", c.ID, c.Description)
	}

	usage := "dcmd <number>"
	if p.IsGlobal() {
		usage = "dcmd global <number>"
	}
	fmt.Fprintf(w, "\nUsage: %s\n", usage)
}
