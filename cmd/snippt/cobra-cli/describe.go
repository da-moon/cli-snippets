package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var describeCmd = &cobra.Command{
	Use:   "describe [title]",
	Short: "Describe a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  describe,
}

func describe(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return commands.Describe("")
	}
	return commands.Describe(args[0])
}

func init() {
	rootCmd.AddCommand(describeCmd)
}
