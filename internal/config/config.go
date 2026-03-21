package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Command struct {
	ID          int
	Template    string
	Description string
}

type ProjectConfig struct {
	Name     string
	Path     string
	Commands []Command
}

func GetConfigDir() string {
	if dir := os.Getenv("DCMD_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		xdgConfig = filepath.Join(home, ".config")
	}
	return filepath.Join(xdgConfig, "dcmd")
}

func LoadProjectConfig(name string) (*ProjectConfig, error) {
	configDir := GetConfigDir()
	filePath := filepath.Join(configDir, name+".conf")
	return ParseConfigFile(filePath)
}

func ParseConfigFile(filePath string) (*ProjectConfig, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	config := &ProjectConfig{
		Name: strings.TrimSuffix(filepath.Base(filePath), ".conf"),
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "# path:") {
			config.Path = strings.TrimSpace(strings.TrimPrefix(line, "# path:"))
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 2 {
			continue
		}

		id, _ := strconv.Atoi(parts[0])
		cmd := Command{
			ID:       id,
			Template: parts[1],
		}
		if len(parts) > 2 && parts[2] != "" {
			cmd.Description = parts[2]
		} else {
			cmd.Description = parts[1]
		}
		config.Commands = append(config.Commands, cmd)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return config, nil
}

func SaveProjectConfig(config *ProjectConfig) error {
	configDir := GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}

	filePath := filepath.Join(configDir, config.Name+".conf")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	if config.Path != "" {
		fmt.Fprintf(writer, "# path: %s\n\n", config.Path)
	}

	for _, cmd := range config.Commands {
		fmt.Fprintf(writer, "%d|%s|%s\n", cmd.ID, cmd.Template, cmd.Description)
	}

	return writer.Flush()
}

func GetCurrentProjectName() string {
	pwd, _ := os.Getwd()
	return filepath.Base(pwd)
}

func ListProjectNames() []string {
	configDir := GetConfigDir()
	files, _ := filepath.Glob(filepath.Join(configDir, "*.conf"))

	var names []string
	for _, f := range files {
		name := strings.TrimSuffix(filepath.Base(f), ".conf")
		if name != "global" {
			names = append(names, name)
		}
	}
	return names
}
