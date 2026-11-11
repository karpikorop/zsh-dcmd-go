package cli

import (
	"fmt"
	"maps"
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

var runOpts struct {
	global bool
	print  bool
	vars   []string
}

var runCmd = &cobra.Command{
	Use:   "run <number>",
	Short: "Fill the variables of a command and put it on your prompt",
	Long: `Fill the variables of a command and put it on your prompt.

Pass variable values with --var to skip the prompts:
  dcmd run 4 --var FunctionName=invoca

Use --print to write the final command to stdout instead of your prompt.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := targetProject(openStore(), runOpts.global)
		if err != nil {
			return err
		}

		values, err := parseVarFlags(runOpts.vars)
		if err != nil {
			return err
		}
		return runByArg(p, args[0], runOptions{values: values, print: runOpts.print})
	},
}

func init() {
	addGlobalFlag(runCmd, &runOpts.global)
	runCmd.Flags().BoolVar(&runOpts.print, "print", false, "Print the final command to stdout")
	runCmd.Flags().StringArrayVar(&runOpts.vars, "var", nil, "Value for a variable as NAME=VALUE (repeatable)")
}

func runByArg(p *store.Project, arg string, opts runOptions) error {
	c, err := commandByArg(p, arg)
	if err != nil {
		return err
	}
	return runCommand(c, opts)
}

func runCommand(c *store.Command, opts runOptions) error {
	if !opts.print {
		info("Selected: %s", c.Description)
		info("Template: %s\n", c.Run)
	}

	values, err := resolveValues(c.Run, opts.values)
	if err != nil {
		return err
	}
	return emitCommand(template.Render(c.Run, values), opts.print)
}

// resolveValues checks the given values and asks the user for the missing ones.
func resolveValues(command string, given map[string]string) (map[string]string, error) {
	variables := template.Variables(command)
	if err := checkUnknownVariables(variables, given); err != nil {
		return nil, err
	}

	values := make(map[string]string, len(variables))
	maps.Copy(values, given)

	missing := missingVariables(variables, values)
	if len(missing) > 0 && !ui.IsInteractive() {
		return nil, fmt.Errorf("missing values for %s. Pass each one with 'dcmd run <number> --var NAME=VALUE'", quoteNames(missing))
	}

	for _, v := range missing {
		value, err := ui.Input(fmt.Sprintf("Enter value for '%s'", v.Name))
		if err != nil {
			return nil, err
		}
		values[v.Name] = value
	}
	return values, nil
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

func missingVariables(variables []template.Variable, values map[string]string) []template.Variable {
	var missing []template.Variable
	for _, v := range variables {
		if _, ok := values[v.Name]; !ok {
			missing = append(missing, v)
		}
	}
	return missing
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
