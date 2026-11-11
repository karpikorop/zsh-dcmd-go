package cli

import (
	"fmt"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
)

func runByArg(p *store.Project, arg string) error {
	id, err := parseNumber(arg)
	if err != nil {
		return err
	}
	return runCommand(p, id)
}

// runCommand fills the variables of a command and prints the result to stdout for the shell wrapper.
func runCommand(p *store.Project, id int) error {
	c, err := findCommand(p, id)
	if err != nil {
		return err
	}

	info("Selected: %s", c.Description)
	info("Template: %s\n", c.Run)

	finalCmd, err := ui.PromptForVariables(c.Run)
	if err != nil {
		return err
	}

	fmt.Print(finalCmd)
	return nil
}
