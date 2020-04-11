package cobracli

import (
	"fmt"
	commands "github.com/da-moon/cli-snippets/cmd/snippt/commands"
	snippet "github.com/da-moon/cli-snippets/internal/snippet"
	cobra "github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [title]",
	Short: "Export a snippet to json file",
	Args:  cobra.MaximumNArgs(1),
	RunE:  export,
}

var (
	outputFile string
	fileType   string
)

func export(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return commands.Export(
			"",
			outputFile,
			fileType,
		)
	}
	return commands.Export(
		args[0],
		outputFile,
		fileType,
	)
}

func init() {
	exportCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Specify the output path of the snippet")
	exportCmd.Flags().StringVarP(&fileType, "type", "t", snippet.ExportTypeJSON, fmt.Sprintf("Choose export file type. Allowed values are: \"%s\", \"%s\".", snippet.ExportTypeJSON, snippet.ExportTypeShell))
	rootCmd.AddCommand(exportCmd)
}
