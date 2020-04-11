package commands

import (
	"fmt"
	"os"

	config "github.com/da-moon/cli-snippets/internal/config"
	util "github.com/da-moon/cli-snippets/internal/util"
)

// Edit ...
func Edit(title string) error {

	// find snippet
	s, err := loadSnippet(title)
	if err != nil {
		return err
	}
	conf, _ := config.Load()
	snippetsMeta, _ := conf.LoadSnippetsMeta()
	command := fmt.Sprintf("%s %s", conf.Editor, s.GetFilePath())
	err = util.Execute(command, os.Stdin, os.Stdout)
	if err != nil {
		return err
	}
	// mark snippetsMeta dirty
	snippetsMeta.IsMetaDirty = true
	err = snippetsMeta.Save()
	if err != nil {
		return err
	}
	return nil
}
