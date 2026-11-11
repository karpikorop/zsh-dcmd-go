package cli

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config [project]",
	Short: "Open a project file in your default GUI editor",
	Long: `Open the file of the current project in your default GUI editor.

Pass a project name to open another project. Use "global" for the global commands.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		p, err := projectForConfig(s, args)
		if err != nil {
			return err
		}

		path := s.FilePath(p)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			if err := s.Save(p); err != nil {
				return err
			}
		}
		return openInEditor(path)
	},
}

func projectForConfig(s *store.Store, args []string) (*store.Project, error) {
	if len(args) == 0 {
		return currentProject(s)
	}
	return loadProjectByName(s, args[0])
}

func openInEditor(path string) error {
	var opener *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		opener = exec.Command("open", path)
	case "linux":
		opener = exec.Command("xdg-open", path)
	case "windows":
		opener = exec.Command("cmd", "/c", "start", path)
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	if err := opener.Run(); err != nil {
		return fmt.Errorf("open editor: %w", err)
	}
	return nil
}
