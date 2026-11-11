package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
	"github.com/karpikorop/zsh-dcmd-go/internal/ui"
	"github.com/spf13/cobra"
)

var linkCmd = &cobra.Command{
	Use:   "link",
	Short: "Link an orphaned project to the current folder",
	Long: `Link an orphaned project to the current folder.

A project is orphaned when none of its stored paths exist, for example after you moved or renamed the folder.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		s := openStore()
		dir, err := os.Getwd()
		if err != nil {
			return err
		}
		if err := ensureFolderIsFree(s, dir); err != nil {
			return err
		}

		orphans, err := orphanedProjects(s)
		if err != nil {
			return err
		}
		if len(orphans) == 0 {
			info("No orphaned projects found.")
			return nil
		}

		chosen, err := promptOrphan(orphans)
		if err != nil {
			return err
		}

		confirmed, err := ui.Confirm(fmt.Sprintf("Link project '%s' to '%s'", chosen.Name, dir))
		if err != nil {
			return err
		}
		if !confirmed {
			info("Link cancelled.")
			return nil
		}

		chosen.Paths = []string{dir}
		if err := s.Save(chosen); err != nil {
			return err
		}

		info("Linked project '%s' to '%s'.", chosen.Name, dir)
		return nil
	},
}

func orphanedProjects(s *store.Store) ([]*store.Project, error) {
	projects, err := s.Projects()
	if err != nil {
		return nil, err
	}

	var orphans []*store.Project
	for _, p := range projects {
		if isOrphaned(p) {
			orphans = append(orphans, p)
		}
	}
	return orphans, nil
}

func isOrphaned(p *store.Project) bool {
	if len(p.Paths) == 0 {
		return false
	}

	for _, path := range p.Paths {
		if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
			return false
		}
	}
	return true
}

func promptOrphan(orphans []*store.Project) (*store.Project, error) {
	info("Found the following orphaned projects:")
	for i, p := range orphans {
		info("  %d. %s (last known paths: %v)", i+1, p.Name, p.Paths)
	}

	input, err := ui.Input("Which one would you like to link? (Enter number)")
	if err != nil {
		return nil, err
	}

	index, err := parseNumber(input)
	if err != nil || index > len(orphans) {
		return nil, errors.New("invalid selection")
	}
	return orphans[index-1], nil
}
