package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/TT-AIXion/github-kanri/internal/config"
	"github.com/TT-AIXion/github-kanri/internal/fsutil"
	"github.com/TT-AIXion/github-kanri/internal/repo"
	"github.com/TT-AIXion/github-kanri/internal/safety"
)

func (a App) runSkillsSync(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills sync", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	target := fs.String("target", "", "target")
	mode := fs.String("mode", "", "mode")
	force := fs.Bool("force", false, "force")
	dryRun := fs.Bool("dry-run", false, "dry run")
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
	targets, err := selectTargets(cfg, *target)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	return a.syncTargets(ctx, cfg, targets, *mode, *force, *dryRun, only, exclude)
}

func (a App) runSkillsLink(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills link", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	target := fs.String("target", "", "target")
	force := fs.Bool("force", false, "force")
	dryRun := fs.Bool("dry-run", false, "dry run")
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
	targets, err := selectTargets(cfg, *target)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	return a.syncTargets(ctx, cfg, targets, string(fsutil.ModeLink), *force, *dryRun, only, exclude)
}

func selectTargets(cfg config.Config, name string) ([]config.SyncTarget, error) {
	if strings.TrimSpace(name) == "" {
		return cfg.SyncTargets, nil
	}
	for _, t := range cfg.SyncTargets {
		if t.Name == name {
			return []config.SyncTarget{t}, nil
		}
	}
	return nil, fmt.Errorf("target not found: %s", name)
}

func (a App) syncTargets(ctx context.Context, cfg config.Config, targets []config.SyncTarget, mode string, force bool, dryRun bool, only []string, exclude []string) int {
	repos, err := repo.Scan(cfg.ReposRoot)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	repos = repo.Filter(repos, only, exclude)
	guard := guardFromConfig(cfg)
	policy := cfg.ConflictPolicy
	if force {
		policy = string(fsutil.ConflictOverwrite)
	}
	syncMode := cfg.SyncMode
	if strings.TrimSpace(mode) != "" {
		syncMode = mode
	}
	if syncMode != string(fsutil.ModeCopy) && syncMode != string(fsutil.ModeMirror) && syncMode != string(fsutil.ModeLink) {
		a.Out.Err(fmt.Sprintf("invalid mode: %s", syncMode), nil)
		return 1
	}
	opts := fsutil.SyncOptions{
		Mode:           fsutil.SyncMode(syncMode),
		ConflictPolicy: fsutil.ConflictPolicy(policy),
		DryRun:         dryRun,
	}
	var results []syncResult
	for _, r := range repos {
		for _, t := range targets {
			if err := guard.CheckPath(t.Src); err != nil {
				a.Out.Err(err.Error(), nil)
				return 1
			}
			opts.Include = t.Include
			opts.Exclude = t.Exclude
			for _, dest := range t.Dest {
				destPath := fsutil.ResolvePath(r.Path, dest)
				if err := guard.CheckPath(destPath); err != nil {
					a.Out.Err(err.Error(), nil)
					return 1
				}
				if err := fsutil.SyncDir(t.Src, destPath, opts); err != nil {
					a.Out.Err(fmt.Sprintf("%s %s: %v", r.Name, t.Name, err), nil)
					return 1
				}
				results = append(results, syncResult{Repo: r.Name, Target: t.Name, Dest: destPath})
			}
		}
	}
	if a.Out.JSON {
		a.Out.OK("skills sync", results)
		return 0
	}
	for _, r := range results {
		a.Out.OK(fmt.Sprintf("%s %s %s", r.Repo, r.Target, r.Dest), nil)
	}
	return 0
}

func guardFromConfig(cfg config.Config) safety.Guard {
	return safety.Guard{
		AllowCommands: cfg.AllowCommands,
		DenyCommands:  cfg.DenyCommands,
		AllowPaths:    cfg.AllowPaths,
		DenyPaths:     cfg.DenyPaths,
	}
}
