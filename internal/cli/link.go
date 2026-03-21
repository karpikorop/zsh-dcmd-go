package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link orphaned configurations to the current directory",
	Run: func(cmd *cobra.Command, args []string) {
		configDir := config.GetConfigDir()
		files, _ := filepath.Glob(filepath.Join(configDir, "*.conf"))

		type orphaned struct {
			file string
			path string
		}
		var orphans []orphaned

		for _, f := range files {
			if filepath.Base(f) == "global.conf" {
				continue
			}
			cfg, err := config.ParseConfigFile(f)
			if err == nil && cfg.Path != "" {
				if _, err := os.Stat(cfg.Path); os.IsNotExist(err) {
					orphans = append(orphans, orphaned{file: f, path: cfg.Path})
				}
			}
		}

		if len(orphans) == 0 {
			fmt.Fprintf(os.Stderr, "No orphaned project configurations found.\n")
			return
		}

		fmt.Fprintf(os.Stderr, "Found the following orphaned configurations:\n")
		for i, o := range orphans {
			fmt.Fprintf(os.Stderr, "  %d. %s (last known path: %s)\n", i+1, filepath.Base(o.file), o.path)
		}

		input, err := ui.PromptInput("Which one would you like to link? (Enter number)")
		if err != nil {
			return
		}
		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(orphans) {
			fmt.Fprintf(os.Stderr, "Invalid selection.\n")
			return
		}

		chosen := orphans[idx-1]
		pwd, _ := os.Getwd()
		confirm, _ := ui.PromptConfirm(fmt.Sprintf("Link config for '%s' to '%s'?", chosen.path, pwd))
		if !confirm {
			return
		}

		cfg, _ := config.ParseConfigFile(chosen.file)
		cfg.Path = pwd
		cfg.Name = config.GetCurrentProjectName()

		// Save new, then remove old if name changed
		if err := config.SaveProjectConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
			return
		}

		if filepath.Base(chosen.file) != cfg.Name+".conf" {
			os.Remove(chosen.file)
		}
		fmt.Fprintf(os.Stderr, "Successfully linked project to '%s.conf'.\n", cfg.Name)
	},
}
