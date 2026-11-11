package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/template"
)

// info prints a message for the user on stderr, so stdout stays clean for lists and JSON.
func info(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
}

func writeJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
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

func printCommand(w io.Writer, c *store.Command) {
	fmt.Fprintf(w, "#%d %s\n", c.ID, c.Description)
	fmt.Fprintf(w, "  %s\n", c.Run)

	variables := template.Variables(c.Run)
	if len(variables) == 0 {
		return
	}

	names := make([]string, 0, len(variables))
	for _, v := range variables {
		names = append(names, v.Name)
	}
	fmt.Fprintf(w, "  Variables: %s\n", strings.Join(names, ", "))
}
