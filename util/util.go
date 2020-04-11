package util

import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/chzyer/readline"
)

const (
	JSON_MARSHAL_PREFIX = ""
	JSON_MARSHAL_INDENT = "  "
	STEP_RANGE_SEP      = "-"
	NEXT_LINE_SUFFIX    = "\\"
)

const (
	SHELL_RED      = "\033[0;31m"
	SHELL_GREEN    = "\033[0;32m"
	SHELL_YELLOW   = "\033[1;33m"
	SHELL_NO_COLOR = "\033[0m"
)

func LoadJsonDataFromFile(filePath string, object interface{}) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(data, object); len(data) > 0 && err != nil {
		return err
	}
	return nil
}

func Scan(prompt string, defaultInp string, historyFile string) (string, error) {
	// create config
	config := &readline.Config{
		Prompt:            prompt,
		HistoryFile:       historyFile,
		HistorySearchFold: true,
		InterruptPrompt:   "^C",
		EOFPrompt:         "exit",
	}
	rl, err := readline.NewEx(config)
	if err != nil {
		return "", err
	}
	defer rl.Close()

	var cmds []string
	for {
		line, err := rl.ReadlineWithDefault(defaultInp)
		if err != nil {
			break
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if strings.HasSuffix(line, NEXT_LINE_SUFFIX) {
			cmds = append(cmds, strings.TrimRight(line, NEXT_LINE_SUFFIX))
			rl.SetPrompt("> ")
			continue
		} else {
			cmds = append(cmds, line)
		}
		cmd := strings.Join(cmds, " ")
		cmds = cmds[:0]
		rl.SetPrompt(prompt)
		rl.SaveHistory(cmd)
		return cmd, nil
	}
	return "", errors.New("cancelled")
}

func Execute(command string, r io.Reader, w io.Writer) error {
	// Retrieve the default shell. Otherwise, fallback to `sh`
	defaultShell, isPresent := os.LookupEnv("SHELL")
	if !isPresent {
		defaultShell = "sh"
	}
	cmd := exec.Command(defaultShell, "-c", strings.TrimSpace(command))
	cmd.Stdin = r
	cmd.Stdout = w
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func GetOrCreatePath(loc string, perm os.FileMode, isDir bool) error {
	dirPath := path.Dir(loc)
	if isDir {
		dirPath = loc
	}
	if _, err := os.Stat(loc); os.IsNotExist(err) {
		if err = os.MkdirAll(dirPath, perm); err != nil {
			return err
		}
		if !isDir {
			f, err := os.Create(loc)
			if err != nil {
				return err
			}
			defer f.Close()
		}
	}
	return nil
}
