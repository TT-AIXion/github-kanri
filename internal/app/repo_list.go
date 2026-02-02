package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/TT-AIXion/github-kanri/internal/gitutil"
	"github.com/TT-AIXion/github-kanri/internal/repo"
)

func (a App) runRepoList(_ context.Context, args []string) int {
	fs := flag.NewFlagSet("repo list", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var only multiFlag
	var exclude multiFlag
	fs.Var(&only, "only", "only patterns")
	fs.Var(&exclude, "exclude", "exclude patterns")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
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
	repos = repo.Filter(repos, only, exclude)
	if a.Out.JSON {
		a.Out.OK("repo list", repos)
		return 0
	}
	for _, r := range repos {
		a.Out.OK(fmt.Sprintf("%s %s", r.Name, r.Path), nil)
	}
	return 0
}

func (a App) runRepoStatus(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repo status", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	var only multiFlag
	var exclude multiFlag
	fs.Var(&only, "only", "only patterns")
	fs.Var(&exclude, "exclude", "exclude patterns")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
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
	repos = repo.Filter(repos, only, exclude)
	runner := buildRunner(cfg, false)
	var out []repoStatus
	for _, r := range repos {
		dirty := false
		status, err := gitutil.StatusPorcelain(ctx, runner, r.Path)
		if err != nil {
			a.Out.Err(fmt.Sprintf("%s: %v", r.Name, err), nil)
			return 1
		}
		if status != "" {
			dirty = true
		}
		out = append(out, repoStatus{Name: r.Name, Path: r.Path, Dirty: dirty})
	}
	if a.Out.JSON {
		a.Out.OK("repo status", out)
		return 0
	}
	for _, r := range out {
		if r.Dirty {
			a.Out.Warn(fmt.Sprintf("%s dirty", r.Name), nil)
			continue
		}
		a.Out.OK(fmt.Sprintf("%s clean", r.Name), nil)
	}
	return 0
}

func (a App) runRepoRecent(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repo recent", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	limit := fs.Int("limit", 20, "limit")
	var only multiFlag
	var exclude multiFlag
	fs.Var(&only, "only", "only patterns")
	fs.Var(&exclude, "exclude", "exclude patterns")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
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
	repos = repo.Filter(repos, only, exclude)
	runner := buildRunner(cfg, false)
	var recents []repoRecent
	for _, r := range repos {
		unix, err := gitutil.LastCommitUnix(ctx, runner, r.Path)
		if err != nil {
			recents = append(recents, repoRecent{Name: r.Name, Path: r.Path, HasCommits: false})
			continue
		}
		recents = append(recents, repoRecent{Name: r.Name, Path: r.Path, Unix: unix, Timestamp: time.Unix(unix, 0), HasCommits: true})
	}
	sort.Slice(recents, func(i, j int) bool { return recents[i].Unix > recents[j].Unix })
	if *limit > 0 && len(recents) > *limit {
		recents = recents[:*limit]
	}
	if a.Out.JSON {
		a.Out.OK("repo recent", recents)
		return 0
	}
	for _, r := range recents {
		if !r.HasCommits {
			a.Out.Warn(fmt.Sprintf("%s no-commits", r.Name), nil)
			continue
		}
		a.Out.OK(fmt.Sprintf("%s %s", r.Name, r.Timestamp.Format(time.RFC3339)), nil)
	}
	return 0
}
