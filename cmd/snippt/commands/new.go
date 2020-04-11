package commands

import (
	snippet "github.com/da-moon/cli-snippets/internal/snippet"
)

// Create ...
func Create(
	lastCmds int,
	title string,

) error {
	// set up history
	histCmds, err := snippet.ReadShellHistory()
	if err != nil {
		return err
	}
	if err = snippet.SetUpHistFile(histCmds); err != nil {
		return err
	}
	defer snippet.RemoveHistFile()
	// load config and snippets
	_, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// create snippet
	initialDefaultCmds := histCmds[len(histCmds)-(lastCmds+1) : len(histCmds)-1]
	newSnippet, err := snippet.NewSnippet(title, initialDefaultCmds)
	if err != nil {
		return err
	}
	// add new sninppet to snippets meta and save
	return snippetsMeta.SaveNewSnippet(newSnippet)
}
