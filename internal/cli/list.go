package cli

import (
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/spf13/cobra"
)

var listOpts struct {
	global bool
	json   bool
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List the commands of the current project",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := targetProject(openStore(), listOpts.global)
		if err != nil {
			return err
		}
		return printProject(p, listOpts.json)
	},
}

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
			return writeJSON(newCommandView(*c))
		}
		printCommand(os.Stdout, c)
		return nil
	},
}

func init() {
	addGlobalFlag(listCmd, &listOpts.global)
	addJSONFlag(listCmd, &listOpts.json)

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
