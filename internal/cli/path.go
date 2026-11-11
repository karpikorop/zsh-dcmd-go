package cli

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Manage the folders of a project",
}

var pathAddCmd = &cobra.Command{
	Use:   "add <project>",
	Short: "Add the current folder to a project",
	Long: `Add the current folder to a project.

Use it for clones and worktrees, so they share the commands of the original project.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		if err := ensureFolderIsFree(s, dir); err != nil {
			return err
		}

		p, err := loadProjectByName(s, args[0])
		if err != nil {
			return err
		}
		if p.IsGlobal() {
			return errors.New("global commands work in every folder, they have no paths")
		}

		p.AddPath(dir)
		if err := s.Save(p); err != nil {
			return err
		}

		info("Added '%s' to project '%s'.", dir, p.Name)
		return nil
	},
}

func init() {
	pathCmd.AddCommand(pathAddCmd)
}
