package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/AIXion-Team/github-kanri/internal/config"
	"github.com/AIXion-Team/github-kanri/internal/fsutil"
	"github.com/AIXion-Team/github-kanri/internal/gitutil"
)

func (a App) runSkillsClone(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills clone", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	remote := fs.String("remote", "", "remote url")
	force := fs.Bool("force", false, "force")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	cfg, _, err := loadConfig()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	url := strings.TrimSpace(*remote)
	if url == "" {
		url = strings.TrimSpace(cfg.SkillsRemote)
	}
	if url == "" {
		a.Out.Err("skillsRemote or --remote required", nil)
		return 1
	}
	if _, err := os.Stat(cfg.SkillsRoot); err == nil {
		if fsutil.IsGitRepo(cfg.SkillsRoot) {
			runner := buildRunner(cfg, false)
			if err := gitutil.Pull(ctx, runner, cfg.SkillsRoot); err != nil {
				a.Out.Err(err.Error(), nil)
				return 1
			}
			a.Out.OK("skills updated", nil)
			return 0
		}
		if !*force {
			a.Out.Err("skillsRoot exists (use --force)", nil)
			return 1
		}
		if err := os.RemoveAll(cfg.SkillsRoot); err != nil {
			a.Out.Err(err.Error(), nil)
			return 1
		}
	}
	runner := buildRunner(cfg, false)
	if err := gitutil.Clone(ctx, runner, url, cfg.SkillsRoot); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	a.Out.OK("skills cloned", nil)
	return 0
}

func (a App) runSkillsPin(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills pin", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	target := fs.String("target", "", "target")
	ref := fs.String("ref", "", "ref")
	force := fs.Bool("force", false, "force")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if strings.TrimSpace(*target) == "" || strings.TrimSpace(*ref) == "" {
		a.Out.Err("--target and --ref required", nil)
		return 1
	}
	cfg, _, err := loadConfig()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if _, err := selectTargets(cfg, *target); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	runner := buildRunner(cfg, false)
	clean, err := gitutil.IsClean(ctx, runner, cfg.SkillsRoot)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if !clean && !*force {
		a.Out.Err("skillsRoot is dirty (use --force)", nil)
		return 1
	}
	if err := gitutil.Fetch(ctx, runner, cfg.SkillsRoot); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := gitutil.Checkout(ctx, runner, cfg.SkillsRoot, *ref); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	a.Out.OK(fmt.Sprintf("pinned %s", *target), nil)
	return 0
}

func (a App) runSkillsWatch(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills watch", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	target := fs.String("target", "", "target")
	interval := fs.Int("interval", 5, "interval seconds")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if *interval <= 0 {
		*interval = 5
	}
	cfg, _, err := loadConfig()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	targets, err := selectTargets(cfg, *target)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	state, err := skillsState(ctx, cfg)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	a.Out.OK("watch started", nil)
	for {
		select {
		case <-ctx.Done():
			return 0
		case <-time.After(time.Duration(*interval) * time.Second):
			newState, err := skillsState(ctx, cfg)
			if err != nil {
				a.Out.Err(err.Error(), nil)
				return 1
			}
			if newState != state {
				state = newState
				err := a.syncTargets(ctx, cfg, targets, "", false, false, nil, nil)
				if err != 0 {
					return err
				}
			}
		}
	}
}

func skillsState(ctx context.Context, cfg config.Config) (string, error) {
	if _, err := os.Stat(cfg.SkillsRoot); err != nil {
		return "", err
	}
	runner := buildRunner(cfg, false)
	head, err := runner.Run(ctx, cfg.SkillsRoot, "git", "rev-parse", "HEAD")
	if err != nil {
		return "", err
	}
	status, err := runner.Run(ctx, cfg.SkillsRoot, "git", "status", "--porcelain")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(head.Stdout) + "|" + strings.TrimSpace(status.Stdout), nil
}
