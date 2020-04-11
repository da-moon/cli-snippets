package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/da-moon/cli-snippets/snippet"
	"github.com/da-moon/cli-snippets/util"
)

// Config ...
type Config struct {
	SnippetsFile string `json:"snippets_file"`
	SnippetsDir  string `json:"snippets_dir"`
	Editor       string `json:"editor"`
	FilterCmd    string `json:"filter_cmd"`
}

const (
	// DEFAULT_CONFIG_FILE ...
	DEFAULT_CONFIG_FILE = "corgi_conf.json"
	// DEFAULT_SNIPPETS_DIR ...
	DEFAULT_SNIPPETS_DIR = "snippets"
	// DEFAULT_SNIPPETS_FILE ...
	DEFAULT_SNIPPETS_FILE = "snippets.json"
	// DEFAULT_EDITOR ...
	DEFAULT_EDITOR = "vim"
	// DEFAULT_FILTER_CMD_FZF ...
	DEFAULT_FILTER_CMD_FZF = "fzf"
	// DEFAULT_FILTER_CMD_PECO ...
	DEFAULT_FILTER_CMD_PECO = "peco"
)

// MissingDefaultFilterCmdError ...
var MissingDefaultFilterCmdError = errors.New("missing default filter cmd")

// GetDefaultConfigHome ...
func GetDefaultConfigHome() string {
	var configHome string
	var isPresent bool

	configHome, isPresent = os.LookupEnv("XDG_CONFIG_HOME")
	if isPresent {
		configHome = path.Join(configHome, "corgi")
	} else {
		configHome = path.Join(os.Getenv("HOME"), ".corgi")
	}

	return configHome
}

// GetDefaultConfigFile ...
func GetDefaultConfigFile(configHome string) (string, error) {
	var defaultConfigFileLoc = path.Join(configHome, DEFAULT_CONFIG_FILE)
	if err := util.GetOrCreatePath(defaultConfigFileLoc, 0755, false); err != nil {
		return "", err
	}
	return defaultConfigFileLoc, nil
}

// GetDefaultSnippetsDir ...
func GetDefaultSnippetsDir(configHome string) (string, error) {
	var defaultSnippetsDir = path.Join(configHome, DEFAULT_SNIPPETS_DIR)
	if err := util.GetOrCreatePath(defaultSnippetsDir, 0755, true); err != nil {
		return "", err
	}
	return defaultSnippetsDir, nil
}

// GetDefaultSnippetsFile ...
func GetDefaultSnippetsFile(configHome string) (string, error) {
	var defaultSnippetsFile = path.Join(configHome, DEFAULT_SNIPPETS_FILE)
	if err := util.GetOrCreatePath(defaultSnippetsFile, 0755, false); err != nil {
		return "", err
	}
	return defaultSnippetsFile, nil
}

// GetDefaultEditor ...
func GetDefaultEditor() (string, error) {
	editorPath, suc := os.LookupEnv("EDITOR")
	if !suc {
		editorPath, err := exec.LookPath(DEFAULT_EDITOR)
		if err != nil {
			return "", fmt.Errorf("could not find %s (default) in $PATH, update your editor choice with \"corgi configure --editor <path to your editor>\"", DEFAULT_EDITOR)
		}
		return editorPath, nil
	}
	return editorPath, nil
}

// GetDefaultFilterCmd ...
func GetDefaultFilterCmd() (string, error) {
	filterCmdPath, err := exec.LookPath(DEFAULT_FILTER_CMD_PECO)
	if err != nil {
		filterCmdPath = ""
	}
	filterCmdPath, err = exec.LookPath(DEFAULT_FILTER_CMD_FZF)
	if err != nil {
		filterCmdPath = ""
	}
	if filterCmdPath == "" {
		return "", MissingDefaultFilterCmdError
	}
	return filterCmdPath, nil
}

// Load ...
func Load() (*Config, error) {
	// find config dir location
	configHome := GetDefaultConfigHome()
	// loading other config files
	configFile, err := GetDefaultConfigFile(configHome)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	if err = util.LoadJsonDataFromFile(configFile, config); err != nil {
		return nil, err
	}
	// if config file has no content, initialize it with default
	if config.IsNew() {
		// set default snippets file
		snippetsFile, err := GetDefaultSnippetsFile(configHome)
		if err != nil {
			return nil, err
		}
		config.SnippetsFile = snippetsFile
		// set default snippets dir
		snippetsDir, err := GetDefaultSnippetsDir(configHome)
		if err != nil {
			return nil, err
		}
		config.SnippetsDir = snippetsDir
		// set default editor
		editor, err := GetDefaultEditor()
		if err != nil {
			return nil, err
		}
		config.Editor = editor
		// set default filter cmd
		filterCmd, err := GetDefaultFilterCmd()
		if err != nil && err != MissingDefaultFilterCmdError {
			return nil, err
		}
		config.FilterCmd = filterCmd
		// save
		config.Save()
	}
	return config, nil
}

// LoadSnippetsMeta ...
func (c *Config) LoadSnippetsMeta() (*snippet.SnippetsMeta, error) {
	if _, err := os.Stat(c.SnippetsFile); os.IsNotExist(err) {
		return nil, err
	}
	snippetsMeta := &snippet.SnippetsMeta{}
	if err := util.LoadJsonDataFromFile(c.SnippetsFile, snippetsMeta); err != nil {
		return nil, err
	}
	snippetsMeta.SetFileLoc(c.SnippetsFile)
	snippetsMeta.SetSnippetsDir(c.SnippetsDir)
	if snippetsMeta.IsMetaDirty {
		if err := snippetsMeta.SyncWithSnippets(); err != nil {
			return nil, err
		}
	}
	return snippetsMeta, nil
}

// Save ...
func (c *Config) Save() error {
	configHome := GetDefaultConfigHome()
	// get config file
	confFile, err := GetDefaultConfigFile(configHome)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, util.JSON_MARSHAL_PREFIX, util.JSON_MARSHAL_INDENT)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(confFile, data, 0644)
	return err
}

// IsNew ...
func (c *Config) IsNew() bool {
	return c.SnippetsFile == "" && c.SnippetsDir == "" && c.Editor == "" && c.FilterCmd == ""
}
