package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [title]",
	Short: "Remove a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  remove,
}

func remove(cmd *cobra.Command, args []string) error {

	if len(args) == 0 {
		return commands.Remove("")
	}
	return commands.Remove(args[0])

}

func init() {
	rootCmd.AddCommand(removeCmd)
}
