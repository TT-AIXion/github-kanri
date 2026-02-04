package app

import (
	"context"
	"errors"
	"flag"
	"os"

	"github.com/TT-AIXion/github-kanri/internal/repo"
)

func (a App) runRepoCd(_ context.Context, args []string) int {
	fs := flag.NewFlagSet("repo cd", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	pick := fs.Int("pick", 0, "pick index")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if fs.NArg() == 0 {
		a.Out.Err("pattern required", nil)
		return 1
	}
	pattern := fs.Arg(0)
	cfg, _, err := loadConfig()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	repos, err := repo.Scan(cfg.ReposRoot)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	result := repo.Find(repos, pattern)
	selected, err := repo.Pick(result, *pick)
	if err != nil {
		if errors.Is(err, repo.ErrMultipleMatches) {
			return a.handleMultiMatch(result)
		}
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if a.Out.JSON {
		a.Out.OK("repo cd", selected)
		return 0
	}
	a.Out.Raw(selected.Path)
	return 0
}
