package cli

import (
	"errors"
	"fmt"
	"os"

	"github.com/karpikorop/zsh-dcmd-go/internal/store"
)

func openStore() *store.Store {
	return store.New(store.DefaultDir())
}

func currentProject(s *store.Store) (*store.Project, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	p, err := s.FindByDir(dir)
	if errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("no commands configured for '%s'. Run 'dcmd add' to create the first command", dir)
	}
	return p, err
}

// currentOrNewProject returns the project of the current folder, or a new unsaved project for it.
func currentOrNewProject(s *store.Store) (*store.Project, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	p, err := s.FindByDir(dir)
	if errors.Is(err, store.ErrNotFound) {
		return s.NewProject(dir)
	}
	return p, err
}

// loadProjectByName loads a project by its name. The name "global" loads the global project.
func loadProjectByName(s *store.Store, name string) (*store.Project, error) {
	if name == store.GlobalName {
		return s.LoadGlobal()
	}

	p, err := s.Load(name)
	if errors.Is(err, store.ErrNotFound) {
		return nil, fmt.Errorf("project '%s' not found", name)
	}
	return p, err
}

// ensureFolderIsFree fails if a project already lists dir as one of its paths.
func ensureFolderIsFree(s *store.Store, dir string) error {
	owner, err := s.OwnerOf(dir)
	if errors.Is(err, store.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return fmt.Errorf("'%s' already belongs to project '%s'", dir, owner.Name)
}
