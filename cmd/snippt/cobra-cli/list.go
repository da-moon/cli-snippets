package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all snippets",
	RunE:  list,
}

func list(cmd *cobra.Command, args []string) error {
	return commands.List()
}

func init() {
	rootCmd.AddCommand(listCmd)
}
