package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Update snippt configuration",
	RunE:  configure,
}

var (
	editor      string
	filterCmd   string
	snippetsDir string
)

func configure(cmd *cobra.Command, args []string) error {
	return commands.Configure(editor, filterCmd, snippetsDir)
}

func init() {
	configCmd.Flags().StringVar(&filterCmd, "filter-cmd", "", "Set the filter command to use for fuzzy searching snippet (default to fzf)")
	configCmd.Flags().StringVar(&editor, "editor", "", "Set the text editor you would like to use to edit snippet")
	configCmd.Flags().StringVar(&snippetsDir, "snippets-dir", "", "Set the path where all snippets are located")
	rootCmd.AddCommand(configCmd)
}
