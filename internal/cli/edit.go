package cli

import (
	"errors"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing command in the current project",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := currentProject(s)
		if err != nil {
			return err
		}
		return editInProject(s, p)
	},
}

func editInProject(s *store.Store, p *store.Project) error {
	c, err := selectCommand(p, "edit")
	if err != nil {
		return err
	}

	run, err := ui.InputWithDefault("Command", c.Run)
	if err != nil {
		return err
	}
	description, err := ui.InputWithDefault("Description", c.Description)
	if err != nil {
		return err
	}

	if strings.TrimSpace(run) == "" || strings.TrimSpace(description) == "" {
		return errors.New("the command and the description cannot be empty")
	}
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
