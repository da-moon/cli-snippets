package commands

import (
	"fmt"
	color "github.com/fatih/color"
)

// List ...
func List() error {
	// load config & snippets
	_, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	// display
	fmt.Println("Here is the list of snippt snippets saved on your system:")
	for _, s := range snippetsMeta.Snippets {
		color.Yellow("- %s", s.Title)
	}
	return nil
}
