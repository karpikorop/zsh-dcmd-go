package store

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

var ErrNotFound = errors.New("project not found")

// GlobalName is the name of the global project. No folder project can use it.
const GlobalName = "global"

const (
	fileExt     = ".toml"
	projectsDir = "projects"
	globalFile  = GlobalName + fileExt
)

// Store reads and writes project files inside one config folder.
type Store struct {
	dir string
}

func New(dir string) *Store {
	return &Store{dir: dir}
}

// DefaultDir returns $DCMD_CONFIG_DIR, or the dcmd folder inside the XDG config folder.
func DefaultDir() string {
	if dir := os.Getenv("DCMD_CONFIG_DIR"); dir != "" {
		return dir
	}

	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		home, _ := os.UserHomeDir()
		configHome = filepath.Join(home, ".config")
	}
	return filepath.Join(configHome, "dcmd")
}

func (s *Store) Load(name string) (*Project, error) {
	if err := validateName(name); err != nil {
		return nil, err
	}

	p, err := s.read(s.projectFile(name))
	if err != nil {
		return nil, err
	}

	p.Name = name
	return p, nil
}

// LoadGlobal returns the global project. It returns an empty project if the file does not exist.
func (s *Store) LoadGlobal() (*Project, error) {
	p, err := s.read(s.globalFile())
	if errors.Is(err, ErrNotFound) {
		p, err = &Project{}, nil
	}
	if err != nil {
		return nil, err
	}

	p.Name = GlobalName
	p.global = true
	return p, nil
}

func (s *Store) Projects() ([]*Project, error) {
	files, err := filepath.Glob(filepath.Join(s.dir, projectsDir, "*"+fileExt))
	if err != nil {
		return nil, err
	}

	projects := make([]*Project, 0, len(files))
	for _, file := range files {
		p, err := s.Load(strings.TrimSuffix(filepath.Base(file), fileExt))
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, nil
}

// FindByDir returns the project that owns dir or the closest parent folder of dir.
func (s *Store) FindByDir(dir string) (*Project, error) {
	known, err := s.knownFolders()
	if err != nil {
		return nil, err
	}

	current, err := filepath.EvalSymlinks(dir)
	if err != nil {
		current = filepath.Clean(dir)
	}

	for {
		if p := ownerOf(known, current); p != nil {
			return p, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return nil, ErrNotFound
		}
		current = parent
	}
}

// OwnerOf returns the project that lists dir itself as one of its paths.
func (s *Store) OwnerOf(dir string) (*Project, error) {
	known, err := s.knownFolders()
	if err != nil {
		return nil, err
	}

	if p := ownerOf(known, filepath.Clean(dir)); p != nil {
		return p, nil
	}
	return nil, ErrNotFound
}

// NewProject returns an unsaved project for dir. The name is the folder name, with a number added if it is taken.
func (s *Store) NewProject(dir string) (*Project, error) {
	base := filepath.Base(filepath.Clean(dir))
	name := base

	for i := 2; name == GlobalName || s.exists(name); i++ {
		name = base + "-" + strconv.Itoa(i)
	}

	if err := validateName(name); err != nil {
		return nil, err
	}
	return &Project{Name: name, Paths: []string{filepath.Clean(dir)}}, nil
}

func (s *Store) Save(p *Project) error {
	var buf bytes.Buffer
	encoder := toml.NewEncoder(&buf)
	encoder.Indent = ""
	if err := encoder.Encode(p); err != nil {
		return fmt.Errorf("encode project '%s': %w", p.Name, err)
	}

	return writeFileAtomic(s.FilePath(p), buf.Bytes())
}

func (s *Store) FilePath(p *Project) string {
	if p.global {
		return s.globalFile()
	}
	return s.projectFile(p.Name)
}

func (s *Store) read(file string) (*Project, error) {
	var p Project
	meta, err := toml.DecodeFile(file, &p)
	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", file, err)
	}

	if undecoded := meta.Undecoded(); len(undecoded) > 0 {
		return nil, fmt.Errorf("read %s: unknown keys %v", file, undecoded)
	}
	if err := p.validate(); err != nil {
		return nil, fmt.Errorf("read %s: %w", file, err)
	}

	p.normalizePaths()
	return &p, nil
}

// knownFolders returns every stored path that exists on disk, with its file info.
func (s *Store) knownFolders() ([]knownFolder, error) {
	projects, err := s.Projects()
	if err != nil {
		return nil, err
	}

	var known []knownFolder
	for _, p := range projects {
		for _, path := range p.Paths {
			if stat, err := os.Stat(path); err == nil {
				known = append(known, knownFolder{project: p, stat: stat})
			}
		}
	}
	return known, nil
}

func (s *Store) exists(name string) bool {
	_, err := os.Stat(s.projectFile(name))
	return err == nil
}

func (s *Store) projectFile(name string) string {
	return filepath.Join(s.dir, projectsDir, name+fileExt)
}

func (s *Store) globalFile() string {
	return filepath.Join(s.dir, globalFile)
}

func validateName(name string) error {
	if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
		return fmt.Errorf("invalid project name '%s'", name)
	}
	return nil
}

// knownFolder is a stored project path that exists on disk.
type knownFolder struct {
	project *Project
	stat    os.FileInfo
}

// ownerOf compares folders by file identity, so a different letter case or a symlink still matches.
func ownerOf(known []knownFolder, dir string) *Project {
	stat, err := os.Stat(dir)
	if err != nil {
		return nil
	}

	for _, k := range known {
		if os.SameFile(k.stat, stat) {
			return k.project
		}
	}
	return nil
}

// writeFileAtomic writes to a temp file and renames it, so a half-written file never replaces the old one.
func writeFileAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	tmp, err := os.CreateTemp(filepath.Dir(path), "."+filepath.Base(path)+".tmp-*")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	if err := os.Chmod(tmp.Name(), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp.Name(), path)
}
