package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AIXion-Team/github-kanri/internal/gitutil"
	"github.com/AIXion-Team/github-kanri/internal/repo"
)

func (a App) runRepoOpen(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repo open", flag.ContinueOnError)
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
	runner := buildRunner(cfg, false)
	_, err = runner.Run(ctx, "", "code", selected.Path)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	a.Out.OK(fmt.Sprintf("opened %s", selected.Name), nil)
	return 0
}

func (a App) runRepoPath(_ context.Context, args []string) int {
	fs := flag.NewFlagSet("repo path", flag.ContinueOnError)
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
		a.Out.OK("repo path", selected)
		return 0
	}
	a.Out.OK(selected.Path, nil)
	return 0
}

func (a App) runRepoInfo(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repo info", flag.ContinueOnError)
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
	runner := buildRunner(cfg, false)
	origin, _ := gitutil.OriginURL(ctx, runner, selected.Path)
	current, _ := gitutil.CurrentBranch(ctx, runner, selected.Path)
	def, _ := gitutil.DefaultBranch(ctx, runner, selected.Path)
	clean, _ := gitutil.IsClean(ctx, runner, selected.Path)
	info := repoInfo{
		Name:          selected.Name,
		Path:          selected.Path,
		Origin:        origin,
		CurrentBranch: current,
		DefaultBranch: def,
		Dirty:         !clean,
	}
	if a.Out.JSON {
		a.Out.OK("repo info", info)
		return 0
	}
	line := fmt.Sprintf("%s origin=%s current=%s default=%s dirty=%v", info.Name, info.Origin, info.CurrentBranch, info.DefaultBranch, info.Dirty)
	a.Out.OK(line, nil)
	return 0
}

func (a App) runRepoGraph(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repo graph", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	pick := fs.Int("pick", 0, "pick index")
	limit := fs.Int("limit", 20, "limit")
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
	runner := buildRunner(cfg, false)
	log, err := gitutil.LogOneline(ctx, runner, selected.Path, *limit)
	if err != nil {
		a.Out.Warn("no commits", nil)
		return 0
	}
	if a.Out.JSON {
		a.Out.OK("repo graph", map[string]string{"name": selected.Name, "log": log})
		return 0
	}
	if log == "" {
		a.Out.Warn("no commits", nil)
		return 0
	}
	for _, line := range strings.Split(log, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		a.Out.OK(line, nil)
	}
	return 0
}

func (a App) runRepoClone(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repo clone", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	name := fs.String("name", "", "repo name")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if fs.NArg() == 0 {
		a.Out.Err("url required", nil)
		return 1
	}
	url := fs.Arg(0)
	cfg, _, err := loadConfig()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	repoName := *name
	if repoName == "" {
		repoName = strings.TrimSuffix(filepath.Base(url), ".git")
	}
	dest := filepath.Join(cfg.ReposRoot, repoName)
	if _, err := os.Stat(dest); err == nil {
		a.Out.Err("destination exists", nil)
		return 1
	}
	runner := buildRunner(cfg, false)
	if err := gitutil.Clone(ctx, runner, url, dest); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	a.Out.OK(fmt.Sprintf("cloned %s", repoName), nil)
	return 0
}
