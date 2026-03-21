package cli

import (
	"fmt"
	"os"
	"strconv"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var globalCmd = &cobra.Command{
	Use:   "global",
	Short: "Manage global commands",
	Args:  cobra.ArbitraryArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showCommands("global")
			return
		}

		if num, err := strconv.Atoi(args[0]); err == nil {
			executeCommand("global", num)
			return
		}

		cmd.Help()
	},
}

var globalAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new global command",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadProjectConfig("global")
		if err != nil {
			cfg = &config.ProjectConfig{Name: "global"}
		}

		newCmd, err := ui.PromptInputNoPipe("Enter the GLOBAL command")
		if err != nil || newCmd == "" {
			return
		}
		newDesc, err := ui.PromptInputNoPipe("Enter a short description")
		if err != nil || newDesc == "" {
			return
		}

		lastID := 0
		for _, c := range cfg.Commands {
			if c.ID > lastID {
				lastID = c.ID
			}
		}

		cfg.Commands = append(cfg.Commands, config.Command{
			ID:          lastID + 1,
			Template:    newCmd,
			Description: newDesc,
		})

		if err := config.SaveProjectConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
		} else {
			fmt.Fprintf(os.Stderr, "Global command added.\n")
		}
	},
}

var globalEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit an existing global command",
	Run: func(cmd *cobra.Command, args []string) {
		editCommandInProject("global")
	},
}

func init() {
	globalCmd.AddCommand(globalAddCmd)
	globalCmd.AddCommand(globalEditCmd)
}
