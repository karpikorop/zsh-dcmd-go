package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/karpikorop/zsh-dcmd-go/internal/config"
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
	Run: func(cmd *cobra.Command, args []string) {
		if addAllFlag {
			addToMultipleProjects()
			return
		}
		addToCurrentProject()
	},
}

func init() {
	addCmd.Flags().BoolVar(&addAllFlag, "all", false, "Add command to multiple projects at once")
}

func addToCurrentProject() {
	projectName := config.GetCurrentProjectName()
	pwd, _ := os.Getwd()

	cfg, err := config.LoadProjectConfig(projectName)
	if err != nil {
		cfg = &config.ProjectConfig{
			Name: projectName,
			Path: pwd,
		}
		fmt.Fprintf(os.Stderr, "Created new config for project: '%s'\n", projectName)
	}

	newCmd, err := ui.PromptInputNoPipe("Enter the command")
	if err != nil || newCmd == "" {
		return
	}
	newDesc, err := ui.PromptInputNoPipe("Enter a short description")
	if err != nil || newDesc == "" {
		return
	}

	addCommandToConfig(cfg, newCmd, newDesc)

	if err := config.SaveProjectConfig(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving config: %v\n", err)
	} else {
		fmt.Fprintf(os.Stderr, "Command added to '%s'.\n", projectName)
	}
}

func addToMultipleProjects() {
	projects := config.ListProjectNames()
	if len(projects) == 0 {
		fmt.Fprintf(os.Stderr, "No projects found to add a command to.\n")
		return
	}

	fmt.Fprintf(os.Stderr, "Projects available to add a command to:\n")
	for i, name := range projects {
		fmt.Fprintf(os.Stderr, "  %d. %s\n", i+1, name)
	}

	input, err := ui.PromptInput("Enter project numbers (e.g. 1 3), or press ENTER for ALL")
	if err != nil {
		return
	}

	var selected []string
	if strings.TrimSpace(input) == "" {
		selected = projects
		fmt.Fprintf(os.Stderr, "Selecting all %d projects.\n", len(projects))
	} else {
		for _, s := range strings.Fields(input) {
			idx, err := strconv.Atoi(s)
			if err != nil || idx < 1 || idx > len(projects) {
				fmt.Fprintf(os.Stderr, "Skipping invalid selection: %s\n", s)
				continue
			}
			selected = append(selected, projects[idx-1])
		}
	}

	if len(selected) == 0 {
		fmt.Fprintf(os.Stderr, "No valid projects selected.\n")
		return
	}

	newCmd, err := ui.PromptInputNoPipe("Enter the command")
	if err != nil || newCmd == "" {
		return
	}
	newDesc, err := ui.PromptInputNoPipe("Enter a short description")
	if err != nil || newDesc == "" {
		return
	}

	for _, name := range selected {
		cfg, err := config.LoadProjectConfig(name)
		if err != nil {
			continue
		}
		addCommandToConfig(cfg, newCmd, newDesc)
		if err := config.SaveProjectConfig(cfg); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving to '%s': %v\n", name, err)
		} else {
			fmt.Fprintf(os.Stderr, "Command added to '%s'.\n", name)
		}
	}
}

func addCommandToConfig(cfg *config.ProjectConfig, template, description string) {
	lastID := 0
	for _, c := range cfg.Commands {
		if c.ID > lastID {
			lastID = c.ID
		}
	}
	cfg.Commands = append(cfg.Commands, config.Command{
		ID:          lastID + 1,
		Template:    template,
		Description: description,
	})
}
