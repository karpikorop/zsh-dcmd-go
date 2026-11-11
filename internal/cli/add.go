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
	all         bool
	global      bool
	command     string
	description string
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new command to the current project",
	Long: `Add a new command to the current project.

Without flags, dcmd asks for the command and the description.
Pass both flags to add a command without prompts:
  dcmd add --command "yarn build" --description "Build"

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
		if err := addToProject(s, p, addOpts.command, addOpts.description); err != nil {
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
	addCmd.Flags().StringVar(&addOpts.command, "command", "", "The command to save")
	addCmd.Flags().StringVar(&addOpts.description, "description", "", "A short description")
	addCmd.MarkFlagsMutuallyExclusive("all", "global")
	addCmd.MarkFlagsMutuallyExclusive("all", "command")
	addCmd.MarkFlagsMutuallyExclusive("all", "description")
}

func addTarget(s *store.Store, global bool) (p *store.Project, created bool, err error) {
	if global {
		p, err = s.LoadGlobal()
		return p, false, err
	}
	return currentOrNewProject(s)
}

func addToProject(s *store.Store, p *store.Project, run, description string) error {
	run, description, err := fillNewCommand(run, description)
	if err != nil {
		return err
	}

	c := p.AddCommand(run, description)
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

	run, description, err := fillNewCommand("", "")
	if err != nil {
		return err
	}

	var saveErrors []error
	for _, p := range selected {
		c := p.AddCommand(run, description)
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
func fillNewCommand(run, description string) (string, string, error) {
	var err error
	if strings.TrimSpace(run) == "" {
		if run, err = promptRequired("Enter the command"); err != nil {
			return "", "", err
		}
	}

	if strings.TrimSpace(description) == "" {
		if description, err = promptRequired("Enter a short description"); err != nil {
			return "", "", err
		}
	}
	return run, description, nil
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
