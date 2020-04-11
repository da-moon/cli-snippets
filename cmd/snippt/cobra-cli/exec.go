package cobracli

import (
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	cobra "github.com/spf13/cobra"
)

var execCmd = &cobra.Command{
	Use:   "exec [title]",
	Short: "Execute a snippet",
	Args:  cobra.MaximumNArgs(1),
	RunE:  execute,
}

var (
	useDefaultParamValue bool
	stepRange            string
)

func execute(cmd *cobra.Command, args []string) error {
	// load config & snippets
	if len(args) == 0 {
		return commands.Execute(
			"",
			useDefaultParamValue,
			stepRange,
		)
	}
	return commands.Execute(
		args[0],
		useDefaultParamValue,
		stepRange,
	)
}

func init() {
	execCmd.Flags().StringVarP(&stepRange, "step", "s", "", "Select a single step to execute with \"-s <step>\" or a range of steps to execute with \"-s <start>-<end>\", end is optional")
	execCmd.Flags().BoolVar(&useDefaultParamValue, "use-default", false, "Add this flag if you would like to use the default values for your defined template fields without being asked to enter a value")
	rootCmd.AddCommand(execCmd)
}
