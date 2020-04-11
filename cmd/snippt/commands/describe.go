package commands

// Describe ...
func Describe(title string) error {
	// find snippet
	s, err := loadSnippet(title)
	if err != nil {
		return err
	}
	// @TODO should it return the error
	s.Describe()
	return nil
}
