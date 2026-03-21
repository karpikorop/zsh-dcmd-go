package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [project]",
	Short: "Open a config file in your default GUI editor",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := config.GetCurrentProjectName()
		if len(args) > 0 {
			projectName = args[0]
		}

		configDir := config.GetConfigDir()
		filePath := filepath.Join(configDir, projectName+".conf")

		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v\n", err)
			return
		}

		// Ensure the file exists so the editor doesn't complain
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			os.WriteFile(filePath, []byte(""), 0644)
		}

		var err error
		switch runtime.GOOS {
		case "darwin":
			err = exec.Command("open", filePath).Run()
		case "linux":
			err = exec.Command("xdg-open", filePath).Run()
		case "windows":
			err = exec.Command("cmd", "/c", "start", filePath).Run()
		default:
			err = fmt.Errorf("unsupported platform: %s", runtime.GOOS)
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening editor: %v\n", err)
		}
	},
}
