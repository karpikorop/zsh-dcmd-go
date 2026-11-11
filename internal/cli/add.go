package cli

import (
	"errors"
	"fmt"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var addAllFlag bool

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new command to the current project",
	Long: `Add a new command to the current project.

Use --all to add a command to multiple projects at once:
  dcmd add --all`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		if addAllFlag {
			return addToManyProjects(s)
		}

		p, err := currentOrNewProject(s)
		if err != nil {
			return err
		}
		return addToProject(s, p)
	},
}

func init() {
	addCmd.Flags().BoolVar(&addAllFlag, "all", false, "Add command to multiple projects at once")
}

func addToProject(s *store.Store, p *store.Project) error {
	run, description, err := promptNewCommand()
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

	run, description, err := promptNewCommand()
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

func promptNewCommand() (run, description string, err error) {
	run, err = promptRequired("Enter the command")
	if err != nil {
		return "", "", err
	}

	description, err = promptRequired("Enter a short description")
	if err != nil {
		return "", "", err
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
