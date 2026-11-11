package cli

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"slices"
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

func printCommand(w io.Writer, c *store.Command, projectVars map[string]string) {
	fmt.Fprintf(w, "#%d %s\n", c.ID, c.Description)
	fmt.Fprintf(w, "  %s\n", c.Run)
	if c.Confirm {
		fmt.Fprintf(w, "  Needs confirmation\n")
	}

	variables := template.Variables(c.Run)
	if len(variables) == 0 {
		return
	}

	descriptions := make([]string, 0, len(variables))
	for _, v := range variables {
		descriptions = append(descriptions, describeVariable(v, projectVars))
	}
	fmt.Fprintf(w, "  Variables: %s\n", strings.Join(descriptions, ", "))
}

func describeVariable(v template.Variable, projectVars map[string]string) string {
	if value, ok := projectVars[v.Name]; ok {
		return fmt.Sprintf("%s = %s (project)", v.Name, value)
	}
	if v.HasDefault {
		return fmt.Sprintf("%s (default: %s)", v.Name, v.Default)
	}
	return v.Name
}

func printVars(w io.Writer, p *store.Project) {
	if len(p.Vars) == 0 {
		fmt.Fprintf(w, "No variables set for '%s'.\n", p.Name)
		return
	}

	fmt.Fprintf(w, "Variables of '%s':\n", p.Name)
	for _, name := range slices.Sorted(maps.Keys(p.Vars)) {
		fmt.Fprintf(w, "  %s = %s\n", name, p.Vars[name])
	}
}

func printProjects(w io.Writer, projects []*store.Project) {
	if len(projects) == 0 {
		fmt.Fprintln(w, "No projects found. Run 'dcmd add' in a project folder to create one.")
		return
	}

	fmt.Fprintln(w, "Projects:")
	for _, p := range projects {
		fmt.Fprintf(w, "  %s (%d commands): %s\n", p.Name, len(p.Commands), strings.Join(p.Paths, ", "))
	}
}
