package cli

import (
	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/spf13/cobra"
)

var globalCmd = &cobra.Command{
	Use:   "global [number]",
	Short: "Manage global commands",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := openStore().LoadGlobal()
		if err != nil {
			return err
		}
		return listOrRun(p, args)
	},
}

var globalAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new global command",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := s.LoadGlobal()
		if err != nil {
			return err
		}
		return addToProject(s, p, store.Command{})
	},
}

var globalEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing global command",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := s.LoadGlobal()
		if err != nil {
			return err
		}
		return editInProject(s, p, nil, commandChange{})
	},
}

func init() {
	addMenuFlags(globalCmd)
	globalCmd.AddCommand(globalAddCmd)
	globalCmd.AddCommand(globalEditCmd)
}
