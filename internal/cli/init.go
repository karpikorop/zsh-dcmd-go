package cli

import (
	"fmt"
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/shell"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init <shell>",
	Short: "Print the shell code that connects dcmd to your prompt",
	Long: `Print the shell code that connects dcmd to your prompt.

Add this line to your ~/.zshrc:
  eval "$(dcmd init zsh)"

Only zsh is supported, because dcmd uses 'print -z' to put commands on your prompt.`,
	Args:      cobra.ExactArgs(1),
	ValidArgs: []string{"zsh"},
	RunE: func(cmd *cobra.Command, args []string) error {
		if args[0] != "zsh" {
			return fmt.Errorf("unsupported shell '%s'. Only zsh is supported", args[0])
		}

		binPath, err := os.Executable()
		if err != nil {
			return fmt.Errorf("find the dcmd binary: %w", err)
		}

		fmt.Print(shell.Zsh(binPath))
		return nil
	},
}
