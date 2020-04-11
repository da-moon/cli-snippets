package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [file1] [file2...]",
	Short: "Import a snippet from one or multiple json files",
	Args:  cobra.MinimumNArgs(1),
	RunE:  importSnippet,
}

func importSnippet(cmd *cobra.Command, args []string) error {
	return commands.ImportSnippet(args)
}

func init() {
	rootCmd.AddCommand(importCmd)
}
