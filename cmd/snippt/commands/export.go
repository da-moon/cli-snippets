package commands

import (
	"fmt"
	"path"
	"strings"

	snippet "github.com/da-moon/cli-snippets/internal/snippet"
)

// Execute ...
func Export(
	title string,
	outputFile string,
	fileType string,
) error {
	// find snippet corresponds to title
	s, err := loadSnippet(title)
	if err != nil {
		return err
	}
	if len(outputFile) == 0 {
		snipptFileName := path.Base(s.GetFilePath())
		if fileType == snippet.ExportTypeJSON {
			outputFile = fmt.Sprintf("./%s", snipptFileName)
		} else {
			outputFile = fmt.Sprintf("./%s", strings.Split(snipptFileName, ".")[0])
		}
	}

	return s.Export(outputFile, fileType)
}
