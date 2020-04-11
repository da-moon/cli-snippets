package commands

// Remove ...
// @TODO maybe add a method to prevent repetitve
// operation used in loading meta
func Remove(title string) error {
	conf, snippetsMeta, err := loadConfigAndSnippetsMeta()
	if err != nil {
		return err
	}
	if len(title) == 0 {
		title, err = filterSnippetTitle(
			conf.FilterCmd,
			snippetsMeta.GetSnippetTitles(),
		)
		if err != nil {
			return err
		}
	}
	err = snippetsMeta.DeleteSnippet(title)
	if err != nil {
		return err
	}
	return nil
}
