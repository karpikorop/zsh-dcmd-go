package cli

import (
	"fmt"
	"maps"
	"os"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/template"
	"github.com/spf13/cobra"
)

var varOpts struct {
	global bool
	json   bool
}

var varCmd = &cobra.Command{
	Use:   "var",
	Short: "Manage the variables of the current project",
	Long: `List the variables of the current project.

A project variable fills every {{name}} with the same name, so dcmd does not ask for it:
  dcmd var set project my-app-dev`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := targetProject(openStore(), varOpts.global)
		if err != nil {
			return err
		}

		if varOpts.json {
			vars := make(map[string]string, len(p.Vars))
			maps.Copy(vars, p.Vars)
			return writeJSON(vars)
		}
		printVars(os.Stdout, p)
		return nil
	},
}

var varSetCmd = &cobra.Command{
	Use:   "set <name> <value>",
	Short: "Set a project variable",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, value := strings.TrimSpace(args[0]), args[1]
		if !template.ValidName(name) {
			return fmt.Errorf("'%s' is not a valid variable name. It cannot be empty or contain ':', '{' or '}'", name)
		}

		p, err := updateVars(func(p *store.Project) error {
			p.SetVar(name, value)
			return nil
		})
		if err != nil {
			return err
		}

		info("Set '%s' in '%s'.", name, p.Name)
		return nil
	},
}

var varUnsetCmd = &cobra.Command{
	Use:   "unset <name>",
	Short: "Remove a project variable",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := strings.TrimSpace(args[0])
		p, err := updateVars(func(p *store.Project) error {
			if !p.UnsetVar(name) {
				return fmt.Errorf("variable '%s' is not set in '%s'", name, p.Name)
			}
			return nil
		})
		if err != nil {
			return err
		}

		info("Removed '%s' from '%s'.", name, p.Name)
		return nil
	},
}

func init() {
	varCmd.PersistentFlags().BoolVar(&varOpts.global, "global", false, "Use the global variables instead of the current project")
	addJSONFlag(varCmd, &varOpts.json)

	varCmd.AddCommand(varSetCmd)
	varCmd.AddCommand(varUnsetCmd)
}

// updateVars loads the target project, applies change, and saves the project.
func updateVars(change func(p *store.Project) error) (*store.Project, error) {
	s := openStore()
	p, err := targetProject(s, varOpts.global)
	if err != nil {
		return nil, err
	}

	if err := change(p); err != nil {
		return nil, err
	}
	return p, s.Save(p)
}
