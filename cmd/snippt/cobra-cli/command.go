package cobracli

import (
	"os"

	cobra "github.com/spf13/cobra"
)

var appVersion = "v0.2.5"

var rootCmd = &cobra.Command{
	Use:          "snippt",
	Short:        "snippt helps with managin cli snippets and workflow",
	Version:      appVersion,
	SilenceUsage: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute ...
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
