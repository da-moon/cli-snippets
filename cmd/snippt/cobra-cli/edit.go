package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit [title]",
	Short: "Edit a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  edit,
}

func edit(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return commands.Edit("")
	}
	return commands.Edit(args[0])
}

func init() {
	rootCmd.AddCommand(editCmd)
}
