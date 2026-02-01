package safety

import (
	"fmt"
	"path/filepath"

	"github.com/AIXion-Team/github-kanri/internal/match"
)

type Guard struct {
	AllowCommands []string
	DenyCommands  []string
	AllowPaths    []string
	DenyPaths     []string
}

func (g Guard) CheckCommand(command string) error {
	if match.AnyCommand(g.DenyCommands, command) {
		return fmt.Errorf("deny command: %s", command)
	}
	if len(g.AllowCommands) == 0 {
		return nil
	}
	if match.AnyCommand(g.AllowCommands, command) {
		return nil
	}
	return fmt.Errorf("command not allowed: %s", command)
}

func (g Guard) CheckPath(path string) error {
	clean := filepath.Clean(path)
	clean = filepath.ToSlash(clean)
	if match.Any(g.DenyPaths, clean) {
		return fmt.Errorf("deny path: %s", clean)
	}
	if len(g.AllowPaths) == 0 {
		return nil
	}
	if match.Any(g.AllowPaths, clean) {
		return nil
	}
	return fmt.Errorf("path not allowed: %s", clean)
}
