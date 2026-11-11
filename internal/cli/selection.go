package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
)

// parseNumber parses a number of 1 or higher, as used for command IDs and list choices.
func parseNumber(value string) (int, error) {
	number, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || number < 1 {
		return 0, fmt.Errorf("'%s' is not a valid number", value)
	}
	return number, nil
}

func findCommand(p *store.Project, id int) (*store.Command, error) {
	c, ok := p.Command(id)
	if !ok {
		return nil, fmt.Errorf("command %d not found in '%s'", id, p.Name)
	}
	return c, nil
}

// selectCommand shows the command list and asks the user for a command number.
func selectCommand(p *store.Project, action string) (*store.Command, error) {
	if len(p.Commands) == 0 {
		return nil, fmt.Errorf("no commands found for '%s'", p.Name)
	}

	printCommands(os.Stderr, p)
	input, err := ui.Input(fmt.Sprintf("Enter the number of the command to %s", action))
	if err != nil {
		return nil, err
	}

	id, err := parseNumber(input)
	if err != nil {
		return nil, err
	}
	return findCommand(p, id)
}
