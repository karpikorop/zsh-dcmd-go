package store

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoadKeepsCommandsExactly(t *testing.T) {
	s := New(t.TempDir())
	original := &Project{
		Name:  "app",
		Paths: []string{"/work/app/"},
		Commands: []Command{
			{ID: 1, Description: "Deploy", Run: `yarn deploy --only functions:{{Name}} --project "dev"`},
			{ID: 4, Description: "Pipe | and quotes", Run: `cat a.log | grep 'x' \ y`},
		},
	}

	if err := s.Save(original); err != nil {
		t.Fatal(err)
	}
	loaded, err := s.Load("app")
	if err != nil {
		t.Fatal(err)
	}

	if loaded.Paths[0] != "/work/app" {
		t.Errorf("path not cleaned: %q", loaded.Paths[0])
	}
	for i, want := range original.Commands {
		if loaded.Commands[i] != want {
			t.Errorf("command %d: got %+v, want %+v", i, loaded.Commands[i], want)
		}
	}

	if _, err := s.Load("missing"); !errors.Is(err, ErrNotFound) {
		t.Errorf("missing project: got %v, want ErrNotFound", err)
	}
	if _, err := s.Load("../escape"); err == nil {
		t.Error("expected an error for a name with a path separator")
	}
}

func TestFindByDirSearchesParentFolders(t *testing.T) {
	root := t.TempDir()
	dirs := mustMkdirs(t, root, "work/app/functions", "work/other", "clones/app/src/deep", "elsewhere")
	if err := os.Symlink(dirs["work/app/functions"], filepath.Join(root, "link")); err != nil {
		t.Fatal(err)
	}

	s := New(t.TempDir())
	mustSave(t, s, &Project{Name: "outer", Paths: []string{filepath.Join(root, "work")}})
	mustSave(t, s, &Project{Name: "inner", Paths: []string{filepath.Join(root, "work/app"), filepath.Join(root, "clones/app")}})

	cases := map[string]string{
		"work/app":            "inner",
		"work/app/functions":  "inner",
		"clones/app/src/deep": "inner",
		"work/other":          "outer",
		"link":                "inner",
	}
	for dir, want := range cases {
		p, err := s.FindByDir(filepath.Join(root, dir))
		if err != nil || p.Name != want {
			t.Errorf("FindByDir(%q) = %v, %v; want %q", dir, p, err, want)
		}
	}

	if _, err := s.FindByDir(dirs["elsewhere"]); !errors.Is(err, ErrNotFound) {
		t.Errorf("unknown folder: got %v, want ErrNotFound", err)
	}
	if _, err := s.OwnerOf(dirs["work/app/functions"]); !errors.Is(err, ErrNotFound) {
		t.Errorf("OwnerOf a subfolder: got %v, want ErrNotFound", err)
	}
}

func TestReadRejectsMistakesInFile(t *testing.T) {
	cases := map[string]string{
		"unknown key":   "[[commands]]\nid = 1\ndescription = \"A\"\ncommand = \"ls\"\n",
		"relative path": "paths = [\"~/app\"]\n",
		"duplicate id":  "[[commands]]\nid = 1\nrun = \"a\"\n[[commands]]\nid = 1\nrun = \"b\"\n",
		"zero id":       "[[commands]]\nid = 0\nrun = \"a\"\n",
	}

	for name, content := range cases {
		dir := t.TempDir()
		mustMkdirs(t, dir, "projects")
		if err := os.WriteFile(filepath.Join(dir, "projects", "app.toml"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}

		if _, err := New(dir).Load("app"); err == nil {
			t.Errorf("%s: expected an error", name)
		}
	}
}

func TestNewProjectPicksFreeName(t *testing.T) {
	s := New(t.TempDir())
	mustSave(t, s, &Project{Name: "app"})

	p, err := s.NewProject("/other/app")
	if err != nil {
		t.Fatal(err)
	}
	if p.Name != "app-2" || p.Paths[0] != "/other/app" {
		t.Errorf("got name %q paths %v", p.Name, p.Paths)
	}
}

func TestLoadGlobal(t *testing.T) {
	dir := t.TempDir()
	s := New(dir)

	empty, err := s.LoadGlobal()
	if err != nil || !empty.IsGlobal() || len(empty.Commands) != 0 {
		t.Fatalf("missing global file: got %+v, %v", empty, err)
	}

	empty.AddCommand("echo hi", "Say hi")
	mustSave(t, s, empty)

	if _, err := os.Stat(filepath.Join(dir, "global.toml")); err != nil {
		t.Errorf("global file not written: %v", err)
	}
}

func TestProjectCommandIDs(t *testing.T) {
	p := &Project{}
	p.AddCommand("a", "A")
	p.AddCommand("b", "B")
	p.AddCommand("c", "C")

	if !p.RemoveCommand(2) || p.RemoveCommand(2) {
		t.Fatal("RemoveCommand should succeed once")
	}
	if added := p.AddCommand("d", "D"); added.ID != 4 {
		t.Errorf("new ID after a gap: got %d, want 4", added.ID)
	}
	if _, ok := p.Command(3); !ok {
		t.Error("ID 3 should keep its number after deleting ID 2")
	}

	p.Renumber()
	for i, c := range p.Commands {
		if c.ID != i+1 {
			t.Errorf("renumber: position %d has ID %d", i, c.ID)
		}
	}
}

func mustMkdirs(t *testing.T, root string, names ...string) map[string]string {
	t.Helper()
	dirs := make(map[string]string, len(names))
	for _, name := range names {
		dirs[name] = filepath.Join(root, name)
		if err := os.MkdirAll(dirs[name], 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return dirs
}

func mustSave(t *testing.T, s *Store, p *Project) {
	t.Helper()
	if err := s.Save(p); err != nil {
		t.Fatal(err)
	}
}
