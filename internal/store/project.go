package store

import (
	"fmt"
	"path/filepath"
	"slices"
)

type Command struct {
	ID          int    `toml:"id"`
	Description string `toml:"description"`
	Run         string `toml:"run"`
}

// Project is a named list of commands that belongs to one or more folders.
type Project struct {
	Name     string    `toml:"-"`
	Paths    []string  `toml:"paths,omitempty"`
	Commands []Command `toml:"commands"`

	global bool
}

func (p *Project) IsGlobal() bool {
	return p.global
}

func (p *Project) Command(id int) (*Command, bool) {
	for i := range p.Commands {
		if p.Commands[i].ID == id {
			return &p.Commands[i], true
		}
	}
	return nil, false
}

// AddCommand appends a command with an ID one higher than the highest ID in use.
func (p *Project) AddCommand(run, description string) Command {
	cmd := Command{ID: p.nextID(), Description: description, Run: run}
	p.Commands = append(p.Commands, cmd)
	return cmd
}

// RemoveCommand deletes a command and keeps the IDs of all other commands.
func (p *Project) RemoveCommand(id int) bool {
	index := slices.IndexFunc(p.Commands, func(c Command) bool { return c.ID == id })
	if index == -1 {
		return false
	}

	p.Commands = slices.Delete(p.Commands, index, index+1)
	return true
}

func (p *Project) Renumber() {
	for i := range p.Commands {
		p.Commands[i].ID = i + 1
	}
}

func (p *Project) AddPath(dir string) {
	p.Paths = append(p.Paths, filepath.Clean(dir))
}

func (p *Project) nextID() int {
	maxID := 0
	for _, c := range p.Commands {
		maxID = max(maxID, c.ID)
	}
	return maxID + 1
}

// validate catches mistakes in a hand-edited file.
func (p *Project) validate() error {
	for _, path := range p.Paths {
		if !filepath.IsAbs(path) {
			return fmt.Errorf("path '%s' must be absolute", path)
		}
	}

	seen := make(map[int]bool)
	for _, c := range p.Commands {
		if c.ID < 1 {
			return fmt.Errorf("command '%s' has id %d, ids must be 1 or higher", c.Description, c.ID)
		}
		if seen[c.ID] {
			return fmt.Errorf("id %d is used by more than one command", c.ID)
		}
		seen[c.ID] = true
	}
	return nil
}

func (p *Project) normalizePaths() {
	for i, path := range p.Paths {
		p.Paths[i] = filepath.Clean(path)
	}
}
