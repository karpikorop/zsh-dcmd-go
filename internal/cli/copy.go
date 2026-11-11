package cli

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/template"
	"github.com/spf13/cobra"
)

var copyOpts struct {
	ids []int
}

var copyCmd = &cobra.Command{
	Use:   "copy <project>",
	Short: "Copy commands from another project into the current project",
	Long: `Copy commands from another project into the current project.

  dcmd copy umcp              Copy all commands
  dcmd copy umcp --ids 1,4,7  Copy some commands
  dcmd copy global            Copy the global commands

Copied commands get new numbers after the last one. Commands that already exist are skipped.
Project variables that the copied commands use are copied too, unless the current project has them.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		source, err := loadProjectByName(s, args[0])
		if err != nil {
			return err
		}

		target, created, err := currentOrNewProject(s)
		if err != nil {
			return err
		}
		if !source.IsGlobal() && source.Name == target.Name {
			return errors.New("the current folder already uses this project")
		}

		commands, err := pickCommands(source, copyOpts.ids)
		if err != nil {
			return err
		}

		result := copyCommands(target, commands, source.Vars)
		if len(result.added) > 0 {
			if err := s.Save(target); err != nil {
				return err
			}
		}

		printCopyResult(result, source, target, created)
		return nil
	},
}

func init() {
	copyCmd.Flags().IntSliceVar(&copyOpts.ids, "ids", nil, "Numbers of the commands to copy, e.g. 1,4,7")
}

type copiedCommand struct {
	fromID  int
	command store.Command
}

type copyResult struct {
	added   []copiedCommand
	skipped []store.Command
	vars    []string
}

// pickCommands returns the commands with the given IDs, or all commands if ids is empty.
func pickCommands(p *store.Project, ids []int) ([]store.Command, error) {
	if len(ids) == 0 {
		return p.Commands, nil
	}

	var picked []store.Command
	var missing []string
	for _, id := range ids {
		c, ok := p.Command(id)
		if !ok {
			missing = append(missing, strconv.Itoa(id))
			continue
		}
		picked = append(picked, *c)
	}

	if len(missing) > 0 {
		return nil, fmt.Errorf("commands %s not found in '%s'", strings.Join(missing, ", "), p.Name)
	}
	return picked, nil
}

// copyCommands adds each command that target does not have yet, with the variables it needs.
func copyCommands(target *store.Project, commands []store.Command, sourceVars map[string]string) copyResult {
	var result copyResult
	for _, c := range commands {
		if hasCommandText(target, c.Run) {
			result.skipped = append(result.skipped, c)
			continue
		}

		added := target.AddCommand(c)
		result.added = append(result.added, copiedCommand{fromID: c.ID, command: added})
		result.vars = append(result.vars, copyMissingVars(target, c.Run, sourceVars)...)
	}
	return result
}

func hasCommandText(p *store.Project, run string) bool {
	return slices.ContainsFunc(p.Commands, func(c store.Command) bool { return c.Run == run })
}

// copyMissingVars copies the source values of the variables in run that target does not have.
func copyMissingVars(target *store.Project, run string, sourceVars map[string]string) []string {
	var copied []string
	for _, v := range template.Variables(run) {
		value, inSource := sourceVars[v.Name]
		_, inTarget := target.Vars[v.Name]
		if inSource && !inTarget {
			target.SetVar(v.Name, value)
			copied = append(copied, v.Name)
		}
	}
	return copied
}

func printCopyResult(result copyResult, source, target *store.Project, created bool) {
	if created && len(result.added) > 0 {
		info("Created project '%s' for '%s'.", target.Name, target.Paths[0])
	}

	info("Copying from '%s' to '%s':", source.Name, target.Name)
	for _, c := range result.added {
		info("Copied #%d -> #%d %s", c.fromID, c.command.ID, c.command.Description)
	}
	for _, c := range result.skipped {
		info("Skipped #%d %s (already there)", c.ID, c.Description)
	}
	for _, name := range result.vars {
		info("Copied variable '%s'", name)
	}

	if len(result.added) == 0 {
		info("Nothing copied from '%s'.", source.Name)
	}
}
