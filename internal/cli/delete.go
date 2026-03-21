package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a command from the current project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName := config.GetCurrentProjectName()
		cfg, err := config.LoadProjectConfig(projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "No commands found for '%s'.\n", projectName)
			return
		}

		showCommands(projectName)
		input, err := ui.PromptInput("Enter the number of the command to delete")
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

		confirm, _ := ui.PromptConfirm(fmt.Sprintf("Are you sure you want to delete command #%d: \"%s\"?", num, cfg.Commands[index].Description))
		if !confirm {
			fmt.Fprintf(os.Stderr, "Deletion cancelled.\n")
			return
		}

		cfg.Commands = append(cfg.Commands[:index], cfg.Commands[index+1:]...)

		// Renumber
		for i := range cfg.Commands {
			cfg.Commands[i].ID = i + 1
		}

		if err := config.SaveProjectConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Command deleted and renumbered.\n")
		}
	},
}
