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

Variables in a command:
  {{name}}            dcmd asks for the value
  {{name:default}}    the prompt starts with the default value
  dcmd var set ...    a project variable fills {{name}} without a prompt

Global commands are available everywhere:
  dcmd global             List global commands
  dcmd global <number>    Run a global command
  dcmd global add         Add a new global command
  dcmd global edit        Edit an existing global command

For scripts and AI agents, every command works without prompts:
  dcmd --json                                  List the commands as JSON
  dcmd 4 --var FunctionName=invoca --print     Print the final command to stdout
  dcmd add --command "yarn build" --description "Build"`,
	Args:          cobra.MaximumNArgs(1),
	SilenceUsage:  true,
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := currentProject(openStore())
		if err != nil {
			return err
		}
		return listOrRun(p, args)
	},
}

func Execute() error {
	return RootCmd.Execute()
}

func init() {
	RootCmd.CompletionOptions.DisableDefaultCmd = true
	addMenuFlags(RootCmd)

	RootCmd.AddCommand(showCmd)
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
