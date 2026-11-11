package cli

import "github.com/spf13/cobra"

func addGlobalFlag(cmd *cobra.Command, target *bool) {
	cmd.Flags().BoolVar(target, "global", false, "Use the global commands instead of the current project")
}

func addJSONFlag(cmd *cobra.Command, target *bool) {
	cmd.Flags().BoolVar(target, "json", false, "Print the result as JSON")
}
