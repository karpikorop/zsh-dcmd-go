package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var projectsOpts struct {
	json bool
}

var projectsCmd = &cobra.Command{
	Use:   "projects",
	Short: "List all projects with their folders",
	Long: `List all projects with their folders.

Use a project name with 'dcmd copy <project>' or 'dcmd path add <project>'.
Global commands are not listed here. See 'dcmd global'.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		projects, err := openStore().Projects()
		if err != nil {
			return err
		}

		if projectsOpts.json {
			return writeJSON(newProjectSummaries(projects))
		}
		printProjects(os.Stdout, projects)
		return nil
	},
}

func init() {
	addJSONFlag(projectsCmd, &projectsOpts.json)
}
