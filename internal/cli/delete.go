package cli

import (
	"fmt"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var deleteOpts struct {
	global bool
	yes    bool
}

var deleteCmd = &cobra.Command{
	Use:   "delete [number]",
	Short: "Delete a command from the current project",
	Long: `Delete a command from the current project.

The other commands keep their numbers. Run 'dcmd renumber' to close the gaps.
Pass --yes to skip the confirmation.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := targetProject(s, deleteOpts.global)
		if err != nil {
			return err
		}

		c, err := commandFromArgs(p, args, "delete")
		if err != nil {
			return err
		}

		confirmed, err := confirmDelete(c, deleteOpts.yes)
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

func init() {
	addGlobalFlag(deleteCmd, &deleteOpts.global)
	deleteCmd.Flags().BoolVarP(&deleteOpts.yes, "yes", "y", false, "Delete without asking for confirmation")
}

func confirmDelete(c *store.Command, skip bool) (bool, error) {
	if skip {
		return true, nil
	}
	return ui.Confirm(fmt.Sprintf("Delete command #%d \"%s\"", c.ID, c.Description))
}
