package commands

import (
	snippet "github.com/da-moon/cli-snippets/internal/snippet"
)

// ImportSnippet ...
func ImportSnippet(snippetJSONFiles []string) error {
	// load config and snippets
	_, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// import snippet
	for _, f := range snippetJSONFiles {
		s, err := snippet.LoadSnippet(f)
		if err != nil {
			return err
		}
		err = snippetsMeta.SaveNewSnippet(s)
		if err != nil {
			return err
		}
	}
	return nil
}
