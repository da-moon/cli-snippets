package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new snippt snippet",
	Long: `Create a new snippt snippet from your command line history interactively

Note: if you plan to use other snippt command that takes a snippet title (for example: snippt exec -t <title>), make sure you don't put double quotes around <title>, otherwise weird failure will happen`,
	RunE: create,
}
var (
	lastCmds int
	title    string
)

func create(cmd *cobra.Command, args []string) error {
	return commands.Create(lastCmds, title)
}

func init() {
	newCmd.Flags().IntVarP(&lastCmds, "last", "l", 0, "The number of history commands to look back, they'll be the default for each step. If 0 or unspecified, each step will not have a default.")
	newCmd.Flags().StringVarP(&title, "title", "t", "", "Title of the snippet, do not put any whitespace if you plan to use this snippet for composition")
	rootCmd.AddCommand(newCmd)
}
