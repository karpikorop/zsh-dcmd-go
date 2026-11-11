package cli

import (
	"errors"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

// commandChange holds new values for a command. An empty field keeps the current value.
type commandChange struct {
	run         string
	description string
}

var editOpts struct {
	global bool
	change commandChange
}

var editCmd = &cobra.Command{
	Use:   "edit [number]",
	Short: "Edit an existing command in the current project",
	Long: `Edit an existing command in the current project.

Without flags, dcmd shows the current values for you to edit.
Pass --command and/or --description to change a command without prompts:
  dcmd edit 4 --description "Deploy one function"`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := targetProject(s, editOpts.global)
		if err != nil {
			return err
		}
		return editInProject(s, p, args, editOpts.change)
	},
}

func init() {
	addGlobalFlag(editCmd, &editOpts.global)
	editCmd.Flags().StringVar(&editOpts.change.run, "command", "", "The new command")
	editCmd.Flags().StringVar(&editOpts.change.description, "description", "", "The new description")
}

func editInProject(s *store.Store, p *store.Project, args []string, change commandChange) error {
	c, err := commandFromArgs(p, args, "edit")
	if err != nil {
		return err
	}

	if change == (commandChange{}) {
		if change, err = promptChange(c); err != nil {
			return err
		}
	}

	run := valueOr(change.run, c.Run)
	description := valueOr(change.description, c.Description)
	if run == c.Run && description == c.Description {
		info("No changes made.")
		return nil
	}

	c.Run = run
	c.Description = description
	if err := s.Save(p); err != nil {
		return err
	}

	info("Command #%d updated.", c.ID)
	return nil
}

func promptChange(c *store.Command) (commandChange, error) {
	run, err := ui.InputWithDefault("Command", c.Run)
	if err != nil {
		return commandChange{}, err
	}

	description, err := ui.InputWithDefault("Description", c.Description)
	if err != nil {
		return commandChange{}, err
	}

	if strings.TrimSpace(run) == "" || strings.TrimSpace(description) == "" {
		return commandChange{}, errors.New("the command and the description cannot be empty")
	}
	return commandChange{run: run, description: description}, nil
}

func valueOr(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
