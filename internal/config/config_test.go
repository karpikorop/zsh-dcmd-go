package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseConfigFile(t *testing.T) {
	content := `# path: /tmp/test
1|ls -la|list files
2|echo {{name}}|greet
`
	tmpDir, err := os.MkdirTemp("", "dcmd-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	filePath := filepath.Join(tmpDir, "test-project.conf")
	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	cfg, err := ParseConfigFile(filePath)
	if err != nil {
		t.Fatalf("Failed to parse config: %v", err)
	}

	if cfg.Path != "/tmp/test" {
		t.Errorf("Expected path /tmp/test, got %s", cfg.Path)
	}

	if len(cfg.Commands) != 2 {
		t.Errorf("Expected 2 commands, got %d", len(cfg.Commands))
	}

	if cfg.Commands[0].Template != "ls -la" || cfg.Commands[0].Description != "list files" {
		t.Errorf("Command 1 mismatch: %+v", cfg.Commands[0])
	}
}

func TestSaveProjectConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "dcmd-test-save")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	os.Setenv("DCMD_CONFIG_DIR", tmpDir)
	defer os.Unsetenv("DCMD_CONFIG_DIR")

	cfg := &ProjectConfig{
		Name: "save-test",
		Path: "/test/path",
		Commands: []Command{
			{ID: 1, Template: "test-cmd", Description: "test-desc"},
		},
	}

	if err := SaveProjectConfig(cfg); err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Read back and verify
	readCfg, err := LoadProjectConfig("save-test")
	if err != nil {
		t.Fatalf("Failed to load saved config: %v", err)
	}

	if readCfg.Path != cfg.Path {
		t.Errorf("Path mismatch: %s vs %s", readCfg.Path, cfg.Path)
	}
	if len(readCfg.Commands) != 1 || readCfg.Commands[0].Template != "test-cmd" {
		t.Errorf("Commands mismatch")
	}
}
