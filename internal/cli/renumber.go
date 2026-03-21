package cli

import (
	"fmt"
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
	"github.com/spf13/cobra"
)

var renumberCmd = &cobra.Command{
	Use:   "renumber",
	Short: "Renumber commands in the current project",
	Run: func(cmd *cobra.Command, args []string) {
		projectName := config.GetCurrentProjectName()
		cfg, err := config.LoadProjectConfig(projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "No config found for '%s'.\n", projectName)
			return
		}

		for i := range cfg.Commands {
			cfg.Commands[i].ID = i + 1
		}

		if err := config.SaveProjectConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Commands renumbered.\n")
		}
	},
}
