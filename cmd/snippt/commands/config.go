package commands

// Configure ...
func Configure(
	editor string,
	filterCmd string,
	snippetsDir string,
) error {
	conf, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	if editor != "" {
		conf.Editor = editor
		if err := conf.Save(); err != nil {
			return err
		}
	}
	if filterCmd != "" {
		conf.FilterCmd = filterCmd
		if err := conf.Save(); err != nil {
			return err
		}
	}
	if snippetsDir != "" {
		conf.SnippetsDir = snippetsDir
		if err := conf.Save(); err != nil {
			return err
		}
		snippetsMeta.IsMetaDirty = true
		if err := snippetsMeta.Save(); err != nil {
			return err
		}
	}
	return nil
}
