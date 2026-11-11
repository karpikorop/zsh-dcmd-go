package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var addOpts struct {
	all    bool
	global bool
	input  store.Command
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new command to the current project",
	Long: `Add a new command to the current project.

Without flags, dcmd asks for the command and the description.
Pass both flags to add a command without prompts:
  dcmd add --command "yarn build" --description "Build"

Use --confirm for commands that change a remote service, so agents ask first.

Use --all to add a command to multiple projects at once:
  dcmd add --all`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		if addOpts.all {
			return addToManyProjects(s)
		}

		p, created, err := addTarget(s, addOpts.global)
		if err != nil {
			return err
		}
		if err := addToProject(s, p, addOpts.input); err != nil {
			return err
		}

		if created {
			info("Created project '%s' for '%s'.", p.Name, p.Paths[0])
		}
		return nil
	},
}

func init() {
	addCmd.Flags().BoolVar(&addOpts.all, "all", false, "Add command to multiple projects at once")
	addGlobalFlag(addCmd, &addOpts.global)
	addCmd.Flags().StringVar(&addOpts.input.Run, "command", "", "The command to save")
	addCmd.Flags().StringVar(&addOpts.input.Description, "description", "", "A short description")
	addCmd.Flags().BoolVar(&addOpts.input.Confirm, "confirm", false, "Mark the command as one that needs confirmation")
	for _, flag := range []string{"global", "command", "description", "confirm"} {
		addCmd.MarkFlagsMutuallyExclusive("all", flag)
	}
}

func addTarget(s *store.Store, global bool) (p *store.Project, created bool, err error) {
	if global {
		p, err = s.LoadGlobal()
		return p, false, err
	}
	return currentOrNewProject(s)
}

func addToProject(s *store.Store, p *store.Project, input store.Command) error {
	input, err := fillNewCommand(input)
	if err != nil {
		return err
	}

	c := p.AddCommand(input)
	if err := s.Save(p); err != nil {
		return err
	}

	info("Command #%d added to '%s'.", c.ID, p.Name)
	return nil
}

func addToManyProjects(s *store.Store) error {
	projects, err := s.Projects()
	if err != nil {
		return err
	}
	if len(projects) == 0 {
		return errors.New("no projects found to add a command to")
	}

	selected, err := promptProjectSelection(projects)
	if err != nil {
		return err
	}

	input, err := fillNewCommand(store.Command{})
	if err != nil {
		return err
	}

	var saveErrors []error
	for _, p := range selected {
		c := p.AddCommand(input)
		if err := s.Save(p); err != nil {
			saveErrors = append(saveErrors, fmt.Errorf("save '%s': %w", p.Name, err))
			continue
		}
		info("Command #%d added to '%s'.", c.ID, p.Name)
	}
	return errors.Join(saveErrors...)
}

func promptProjectSelection(projects []*store.Project) ([]*store.Project, error) {
	info("Projects available to add a command to:")
	for i, p := range projects {
		info("  %d. %s", i+1, p.Name)
	}

	input, err := ui.Input("Enter project numbers (e.g. 1 3), or press ENTER for ALL")
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(input) == "" {
		info("Selecting all %d projects.", len(projects))
		return projects, nil
	}

	var selected []*store.Project
	chosen := make(map[int]bool)
	for _, field := range strings.Fields(input) {
		index, err := parseNumber(field)
		if err != nil || index > len(projects) {
			info("Skipping invalid selection: %s", field)
			continue
		}

		if !chosen[index] {
			chosen[index] = true
			selected = append(selected, projects[index-1])
		}
	}

	if len(selected) == 0 {
		return nil, errors.New("no valid projects selected")
	}
	return selected, nil
}

// fillNewCommand asks for the command and the description that were not passed as flags.
func fillNewCommand(input store.Command) (store.Command, error) {
	var err error
	if strings.TrimSpace(input.Run) == "" {
		if input.Run, err = promptRequired("Enter the command"); err != nil {
			return store.Command{}, err
		}
	}

	if strings.TrimSpace(input.Description) == "" {
		if input.Description, err = promptRequired("Enter a short description"); err != nil {
			return store.Command{}, err
		}
	}
	return input, nil
}

// promptRequired treats an empty answer as a cancel.
func promptRequired(label string) (string, error) {
	value, err := ui.Input(label)
	if err != nil {
		return "", err
	}

	if strings.TrimSpace(value) == "" {
		return "", ui.ErrCancelled
	}
	return value, nil
}
