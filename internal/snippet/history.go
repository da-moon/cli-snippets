package snippet

import (
	"bufio"
	"bytes"
	"fmt"
	stacktrace "github.com/palantir/stacktrace"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/da-moon/cli-snippets/internal/util"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
)

const (
	// ShellBash ...
	ShellBash = "bash"
	// ShellZsh ...
	ShellZsh = "zsh"
	// ShellFish ...
	ShellFish = "fish"
	// ShellUnsupported ...
	ShellUnsupported = "unsupported"
)

var shellType string

// TempHistFile ...
var TempHistFile = "/tmp/snippt.hist"

func getFishHistoryPath(homeDir string) string {
	var (
		oldHistFile = path.Join(homeDir, ".config", "fish", "fish_history")
		newHistFile = path.Join(homeDir, ".local", "share", "fish", "fish_history") // for version >= 2.3.0
	)
	// get fish version
	var buf bytes.Buffer
	err := util.Execute("fish --version", nil, &buf)
	version := strings.TrimSpace(strings.Split(buf.String(), ",")[1])[len("version "):]
	if err != nil || version == "" {
		color.Red("Couldn't read $FISH_VERSION, using %s by default", newHistFile)
		return newHistFile
	}
	// check version
	versionInfo := strings.Split(version, ".")
	major, errMajor := strconv.ParseInt(versionInfo[0], 10, 32)
	minor, errMinor := strconv.ParseInt(versionInfo[1], 10, 32)
	patch, errPatch := strconv.ParseInt(versionInfo[2], 10, 32)
	if errMajor != nil || errMinor != nil || errPatch != nil {
		color.Red("Failed to parse version: %s. Defaulting to %s. (You could reset it by setting $HISTFILE)", version, newHistFile)
		return newHistFile
	}
	if major >= 2 && minor >= 3 && patch >= 0 {
		return newHistFile
	}
	return oldHistFile
}

func getHistoryFilePath() (string, error) {
	histFilePath, suc := os.LookupEnv("HISTFILE")
	if !suc {
		color.Red("Could not find HISTFILE in env - Using default based on shell type")
	}
	homeDir, err := homedir.Dir()
	if err != nil {
		return "", err
	}
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		shellType = ShellZsh
		histFilePath = fmt.Sprintf("%s/.zsh_history", homeDir)
	} else if strings.Contains(shell, "bash") {
		shellType = ShellBash
		histFilePath = fmt.Sprintf("%s/.bash_history", homeDir)
	} else if strings.Contains(shell, "fish") {
		shellType = ShellFish
		histFilePath = getFishHistoryPath(homeDir)
	} else {
		shellType = ShellUnsupported
		err = stacktrace.NewError("%s is not supported, currently supporting Bash, Zsh and Fish", shell)
		return "", err
	}
	if _, err := os.Stat(histFilePath); os.IsNotExist(err) {
		return "", err
	}
	return histFilePath, nil
}

// ParseFileToStringArray ...
func ParseFileToStringArray(filePath string, parser CommandParser) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parsedLine := parser.Parse(line)
		lines = append(lines, parsedLine)
	}
	return lines, nil
}

// ReadShellHistory ...
func ReadShellHistory() ([]string, error) {
	histFilePath, err := getHistoryFilePath()
	if err != nil {
		return nil, err
	}
	parser, err := GetCmdParser(shellType)
	if err != nil {
		return nil, err
	}
	lines, err := ParseFileToStringArray(histFilePath, parser)
	return lines, nil
}

// SetUpHistFile ...
func SetUpHistFile(histCmds []string) error {
	// write commands to temp history file
	f, err := os.OpenFile(TempHistFile, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer f.Close()
	for _, cmd := range histCmds {
		if _, err := f.WriteString(fmt.Sprintf("%s\n", cmd)); err != nil {
			return err
		}
	}
	return nil
}

// RemoveHistFile ...
func RemoveHistFile() error {
	if err := os.Remove(TempHistFile); err != nil {
		return err
	}
	return nil
}
