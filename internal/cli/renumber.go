package cli

import (
	"github.com/spf13/cobra"
)

var renumberCmd = &cobra.Command{
	Use:   "renumber",
	Short: "Renumber commands in the current project from 1",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := currentProject(s)
		if err != nil {
			return err
		}

		p.Renumber()
		if err := s.Save(p); err != nil {
			return err
		}

		info("Commands renumbered.")
		return nil
	},
}
