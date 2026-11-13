package cli

import (
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/spf13/cobra"
)

var showOpts struct {
	global bool
	json   bool
}

var showCmd = &cobra.Command{
	Use:   "show <number>",
	Short: "Show one command with its template and variables",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := targetProject(openStore(), showOpts.global)
		if err != nil {
			return err
		}

		c, err := commandByArg(p, args[0])
		if err != nil {
			return err
		}

		if showOpts.json {
			return writeJSON(newCommandView(*c, p.Vars))
		}
		printCommand(os.Stdout, c, p.Vars)
		return nil
	},
}

func init() {
	addGlobalFlag(showCmd, &showOpts.global)
	addJSONFlag(showCmd, &showOpts.json)
}

func printProject(p *store.Project, asJSON bool) error {
	if asJSON {
		return writeJSON(newProjectView(p))
	}

	printCommands(os.Stdout, p)
	return nil
}
