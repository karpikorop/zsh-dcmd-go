package cli

import (
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
  dcmd global edit        Edit an existing global command

For scripts and AI agents, every command works without prompts:
  dcmd list --json
  dcmd run 4 --var FunctionName=invoca --print
  dcmd add --command "yarn build" --description "Build"`,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := currentProject(openStore())
		if err != nil {
			return err
		}

		if len(args) == 0 {
			return printProject(p, false)
		}
		return runByArg(p, args[0], runOptions{})
	},
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(showCmd)
	RootCmd.AddCommand(runCmd)
	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(editCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(renumberCmd)
	RootCmd.AddCommand(varCmd)
	RootCmd.AddCommand(globalCmd)
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(linkCmd)
	RootCmd.AddCommand(pathCmd)
	RootCmd.AddCommand(copyCmd)
	RootCmd.AddCommand(projectsCmd)
	RootCmd.AddCommand(initCmd)
}
