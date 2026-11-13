package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/template"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

// bufferFileEnv names the file where the shell wrapper expects the command for the prompt.
const bufferFileEnv = "DCMD_BUFFER_FILE"

type runOptions struct {
	values map[string]string
	print  bool
}

var menuOpts struct {
	json  bool
	print bool
	vars  []string
}

func addMenuFlags(cmd *cobra.Command) {
	addJSONFlag(cmd, &menuOpts.json)
	cmd.Flags().BoolVar(&menuOpts.print, "print", false, "Print the final command to stdout instead of your prompt")
	cmd.Flags().StringArrayVar(&menuOpts.vars, "var", nil, "Value for a variable as NAME=VALUE (repeatable)")
}

// listOrRun lists the commands of p, or runs the command whose number is in args.
func listOrRun(p *store.Project, args []string) error {
	if len(args) == 0 {
		if menuOpts.print || len(menuOpts.vars) > 0 {
			return errors.New("--print and --var need a command number")
		}
		return printProject(p, menuOpts.json)
	}

	if menuOpts.json {
		return fmt.Errorf("--json lists all commands. For one command, use 'dcmd show %s --json'", args[0])
	}
	values, err := parseVarFlags(menuOpts.vars)
	if err != nil {
		return err
	}
	return runByArg(p, args[0], runOptions{values: values, print: menuOpts.print})
}

func runByArg(p *store.Project, arg string, opts runOptions) error {
	c, err := commandByArg(p, arg)
	if err != nil {
		return err
	}
	return runCommand(p, c, opts)
}

func runCommand(p *store.Project, c *store.Command, opts runOptions) error {
	if !opts.print {
		info("Selected: %s", c.Description)
		info("Template: %s\n", c.Run)
	}

	values, err := resolveValues(c.Run, p.Vars, opts.values)
	if err != nil {
		return err
	}
	return emitCommand(template.Render(c.Run, values), opts.print)
}

// resolveValues picks a value for each variable. The order is: given values, project variables,
// then a prompt (or the default when there is no terminal).
func resolveValues(command string, projectVars, given map[string]string) (map[string]string, error) {
	variables := template.Variables(command)
	err := checkUnknownVariables(variables, given)
	if err != nil {
		return nil, err
	}

	values := make(map[string]string, len(variables))
	var missing []template.Variable
	for _, v := range variables {
		if value, ok := given[v.Name]; ok {
			values[v.Name] = value
		} else if value, ok := projectVars[v.Name]; ok {
			values[v.Name] = value
		} else {
			missing = append(missing, v)
		}
	}

	if ui.IsInteractive() {
		err = promptForValues(values, missing)
	} else {
		err = fillDefaults(values, missing)
	}

	if err != nil {
		return nil, err
	}
	return values, nil
}

// fillDefaults sets each missing variable to its default. It fails if a variable has no default.
func fillDefaults(values map[string]string, missing []template.Variable) error {
	var noDefault []template.Variable
	for _, v := range missing {
		if v.HasDefault {
			values[v.Name] = v.Default
		} else {
			noDefault = append(noDefault, v)
		}
	}

	if len(noDefault) > 0 {
		return fmt.Errorf("missing values for %s. Pass each one with 'dcmd <number> --var NAME=VALUE'", quoteNames(noDefault))
	}
	return nil
}

func promptForValues(values map[string]string, missing []template.Variable) error {
	for _, v := range missing {
		label := fmt.Sprintf("Enter value for '%s'", v.Name)

		var value string
		var err error
		if v.HasDefault {
			value, err = ui.InputWithDefault(label, v.Default)
		} else {
			value, err = ui.Input(label)
		}
		if err != nil {
			return err
		}
		values[v.Name] = value
	}
	return nil
}

func checkUnknownVariables(variables []template.Variable, given map[string]string) error {
	known := make(map[string]bool, len(variables))
	for _, v := range variables {
		known[v.Name] = true
	}

	for name := range given {
		if !known[name] {
			return fmt.Errorf("the command has no variable '%s'", name)
		}
	}
	return nil
}

func quoteNames(variables []template.Variable) string {
	quoted := make([]string, 0, len(variables))
	for _, v := range variables {
		quoted = append(quoted, "'"+v.Name+"'")
	}
	return strings.Join(quoted, ", ")
}

func parseVarFlags(flags []string) (map[string]string, error) {
	values := make(map[string]string, len(flags))
	for _, flag := range flags {
		name, value, ok := strings.Cut(flag, "=")
		name = strings.TrimSpace(name)
		if !ok || name == "" {
			return nil, fmt.Errorf("--var '%s' must look like NAME=VALUE", flag)
		}
		values[name] = value
	}
	return values, nil
}

// emitCommand hands the command to the shell wrapper. It prints to stdout when asked, or when no wrapper runs dcmd.
func emitCommand(command string, print bool) error {
	if file := os.Getenv(bufferFileEnv); file != "" && !print {
		return os.WriteFile(file, []byte(command), 0o600)
	}

	_, err := fmt.Println(command)
	return err
}
