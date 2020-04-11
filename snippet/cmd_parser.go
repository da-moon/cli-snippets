package snippet

import (
	"fmt"
	"strings"
)

// CommandParser ...
type CommandParser interface {
	Parse(string) string
}

// BashCmdParser ...
type BashCmdParser struct{}

// ZshCmdParser ...
type ZshCmdParser struct{}

// FishCmdParser ...
type FishCmdParser struct{}

// Parse ...
func (z ZshCmdParser) Parse(line string) string {
	parts := strings.Split(line, ";")
	return strings.Join(parts[1:], ";")
}

// Parse ...
func (b BashCmdParser) Parse(line string) string {
	return line
}

// Parse ...
func (f FishCmdParser) Parse(line string) string {
	fishCmdPrefix := "- cmd: "
	if strings.HasPrefix(line, fishCmdPrefix) {
		return line[len(fishCmdPrefix):]
	}
	return ""
}

// GetCmdParser ...
func GetCmdParser(shellType string) (CommandParser, error) {
	if shellType == SHELL_ZSH {
		return ZshCmdParser{}, nil
	} else if shellType == SHELL_BASH {
		return BashCmdParser{}, nil
	} else if shellType == SHELL_FISH {
		return FishCmdParser{}, nil
	}
	return nil, fmt.Errorf("unsupported shell type: %s", shellType)
}
