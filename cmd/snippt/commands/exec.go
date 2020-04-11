package commands

// Execute ...
func Execute(
	title string,
	useDefaultParamValue bool,
	stepRange string,
) error {
	// load config & snippets
	s, err := loadSnippet(title)
	if err != nil {
		return err
	}
	// @TODO should i return this error
	s.Execute(useDefaultParamValue, stepRange)
	return nil
}
