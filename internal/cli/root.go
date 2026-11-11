package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "dcmd [number]",
	Short: "Directory-Specific Command Manager",
	Long: `dcmd — Directory-Specific Command Manager

Run without arguments to list saved commands for the current project.
Run with a number to put that command on your prompt (e.g. dcmd 3).

Global commands are available everywhere:
  dcmd global             List global commands
  dcmd global <number>    Run a global command
  dcmd global add         Add a new global command
  dcmd global edit        Edit an existing global command`,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := currentProject(openStore())
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

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.SetOut(os.Stderr)
	RootCmd.SetErr(os.Stderr)
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(editCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(renumberCmd)
	RootCmd.AddCommand(globalCmd)
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(linkCmd)
	RootCmd.AddCommand(pathCmd)
}
