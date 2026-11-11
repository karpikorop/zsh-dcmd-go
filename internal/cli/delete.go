package cli

import (
	"fmt"

	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a command from the current project",
	Long: `Delete a command from the current project.

The other commands keep their numbers. Run 'dcmd renumber' to close the gaps.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := currentProject(s)
		if err != nil {
			return err
		}

		c, err := selectCommand(p, "delete")
		if err != nil {
			return err
		}

		confirmed, err := ui.Confirm(fmt.Sprintf("Delete command #%d \"%s\"", c.ID, c.Description))
		if err != nil {
			return err
		}
		if !confirmed {
			info("Deletion cancelled.")
			return nil
		}

		id := c.ID
		p.RemoveCommand(id)
		if err := s.Save(p); err != nil {
			return err
		}

		info("Command #%d deleted.", id)
		return nil
	},
}
