package app

import (
	"context"
	"errors"
	"fmt"

	"github.com/TT-AIXion/github-kanri/internal/output"
)

type App struct {
	Version string
	Out     output.Writer
}

type ExitError struct {
	Code int
	Err  error
}

func (e ExitError) Error() string {
	return e.Err.Error()
}

var ErrMultiMatch = errors.New("multiple matches")

func (a App) Run(ctx context.Context, args []string) int {
	if len(args) == 0 {
		a.printUsage()
		return 0
	}
	switch args[0] {
	case "help", "-h", "--help":
		a.printUsage()
		return 0
	case "repo":
		return a.runRepo(ctx, args[1:])
	case "clone":
		return a.runRepoClone(ctx, args[1:])
	case "skills":
		return a.runSkills(ctx, args[1:])
	case "config":
		return a.runConfig(ctx, args[1:])
	case "doctor":
		return a.runDoctor(ctx, args[1:])
	case "version":
		return a.runVersion()
	default:
		a.Out.Err(fmt.Sprintf("unknown command: %s", args[0]), nil)
		a.printUsage()
		return 1
	}
}

func (a App) printUsage() {
	a.Out.Raw(`gkn <command>

Commands:
  clone <url> [--name repo]
  repo      repo operations
  skills    skills sync operations
  config    config operations
  doctor    environment checks
  version   show version

Run "gkn <command> --help" for details.`)
}
