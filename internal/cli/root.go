package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "dcmd",
	Short: "Directory-Specific Command Manager",
	Long: `dcmd — Directory-Specific Command Manager

Run without arguments to list saved commands for the current project.
Run with a number to execute that command (e.g. dcmd 3).

Global commands are available everywhere:
  dcmd global             List global commands
  dcmd global <number>    Run a global command
  dcmd global add         Add a new global command
  dcmd global edit        Edit an existing global command`,
	Args: cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := config.GetCurrentProjectName()
		if len(args) == 0 {
			showCommands(projectName)
			return
		}

		if num, err := strconv.Atoi(args[0]); err == nil {
			executeCommand(projectName, num)
			return
		}

		cmd.Help()
	},
}

func showCommands(projectName string) {
	cfg, err := config.LoadProjectConfig(projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No commands configured for project: '%s'\n", projectName)
		if projectName == "global" {
			fmt.Fprintf(os.Stderr, "Run 'dcmd global add' to create the first command.\n")
		} else {
			fmt.Fprintf(os.Stderr, "Run 'dcmd add' to create the first command.\n")
		}
		return
	}

	if len(cfg.Commands) == 0 {
		fmt.Fprintf(os.Stderr, "No commands found for '%s'.\n", projectName)
		return
	}

	fmt.Fprintf(os.Stderr, "Available commands for '%s':\n", projectName)
	for _, c := range cfg.Commands {
		fmt.Fprintf(os.Stderr, "  %d. %s\n", c.ID, c.Description)
	}
	fmt.Fprintf(os.Stderr, "\nUsage: dcmd <number>\n")
}

func executeCommand(projectName string, num int) {
	cfg, err := config.LoadProjectConfig(projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No commands configured for project: '%s'\n", projectName)
		return
	}

	var target *config.Command
	for _, c := range cfg.Commands {
		if c.ID == num {
			target = &c
			break
		}
	}

	if target == nil {
		fmt.Fprintf(os.Stderr, "Command %d not found.\n", num)
		return
	}

	fmt.Fprintf(os.Stderr, "Selected: %s\n", target.Description)
	fmt.Fprintf(os.Stderr, "Template: %s\n\n", target.Template)

	finalCmd, err := ui.PromptForVariables(target.Template)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Execution cancelled.\n")
		return
	}

	// Output the final command to Stdout so the shell wrapper can capture it
	fmt.Print(finalCmd)
}

func init() {
	RootCmd.SetOut(os.Stderr)
	RootCmd.SetErr(os.Stderr)
	RootCmd.CompletionOptions.DisableDefaultCmd = true

	RootCmd.AddCommand(addCmd)
	RootCmd.AddCommand(editCmd)
	RootCmd.AddCommand(configCmd)
	RootCmd.AddCommand(globalCmd)
	RootCmd.AddCommand(deleteCmd)
	RootCmd.AddCommand(linkCmd)
	RootCmd.AddCommand(renumberCmd)
}
