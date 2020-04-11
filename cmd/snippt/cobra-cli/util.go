package cobracli

import (
	"bytes"
	config "github.com/da-moon/cli-snippets/internal/config"
	snippet "github.com/da-moon/cli-snippets/internal/snippet"
	util "github.com/da-moon/cli-snippets/internal/util"
	color "github.com/fatih/color"
	stacktrace "github.com/palantir/stacktrace"
	"io"
	"os"
	"strings"
)

// ErrMissingSnippetTitle ...
var ErrMissingSnippetTitle = stacktrace.NewError("snippet title is not selected")

func loadConfigAndSnippetsMeta() (*config.Config, *snippet.SnippetsMeta, error) {
	// load config
	conf, err := config.Load()
	if err != nil {
		return nil, nil, err
	}
	// Load snippets
	snippets, err := conf.LoadSnippetsMeta()
	if err != nil {
		return nil, nil, err
	}
	return conf, snippets, nil
}

func filter(filterCmd string, candidates []string) (string, error) {
	var buf bytes.Buffer
	inputs := strings.Join(candidates, "\n")
	ws := io.MultiWriter(os.Stdout, &buf)
	if err := util.Execute(filterCmd, strings.NewReader(inputs), ws); err != nil {
		return "", err
	}
	result := strings.Trim(strings.TrimSpace(buf.String()), "\n")
	return result, nil
}

func filterSnippetTitle(filterCmd string, titles []string) (string, error) {
	if filterCmd != "" {
		title, err := filter(filterCmd, titles)
		if err != nil || title == "" {
			return "", ErrMissingSnippetTitle
		}
		return title, nil
	}
	color.Red("Install a fuzzy finder (\"fzf\" or \"peco\") to enable interactive selection")
	return "", ErrMissingSnippetTitle
}
