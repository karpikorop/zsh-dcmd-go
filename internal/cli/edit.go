package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing command in the current project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName := config.GetCurrentProjectName()
		editCommandInProject(projectName)
	},
}

func editCommandInProject(projectName string) {
	cfg, err := config.LoadProjectConfig(projectName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "No commands found for '%s'.\n", projectName)
		return
	}

	if len(cfg.Commands) == 0 {
		fmt.Fprintf(os.Stderr, "No commands found for '%s'.\n", projectName)
		return
	}

	showCommands(projectName)

	input, err := ui.PromptInput("Enter the number of the command to edit")
	if err != nil {
		return
	}

	num, err := strconv.Atoi(input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid number.\n")
		return
	}

	index := -1
	for i, c := range cfg.Commands {
		if c.ID == num {
			index = i
			break
		}
	}

	if index == -1 {
		fmt.Fprintf(os.Stderr, "Command %d not found.\n", num)
		return
	}

	current := cfg.Commands[index]

	newTemplate, err := ui.PromptInputWithDefaultNoPipe("Command", current.Template)
	if err != nil {
		return
	}
	newDesc, err := ui.PromptInputWithDefaultNoPipe("Description", current.Description)
	if err != nil {
		return
	}

	if newTemplate == current.Template && newDesc == current.Description {
		fmt.Fprintf(os.Stderr, "No changes made.\n")
		return
	}

	cfg.Commands[index].Template = newTemplate
	cfg.Commands[index].Description = newDesc

	if err := config.SaveProjectConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Command #%d updated.\n", num)
	}
}
