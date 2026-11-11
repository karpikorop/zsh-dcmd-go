package cli

import (
	"os"

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

		if len(args) == 0 {
			printCommands(os.Stderr, p)
			return nil
		}
		return runByArg(p, args[0])
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
		return addToProject(s, p)
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
		return editInProject(s, p)
	},
}

func init() {
	globalCmd.AddCommand(globalAddCmd)
	globalCmd.AddCommand(globalEditCmd)
}
