package app

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/TT-AIXion/github-kanri/internal/fsutil"
	"github.com/TT-AIXion/github-kanri/internal/repo"
)

func (a App) runSkillsDiff(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills diff", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	target := fs.String("target", "", "target")
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
	repos, err := repo.Scan(cfg.ReposRoot)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	repos = repo.Filter(repos, only, exclude)
	var results []diffResult
	for _, r := range repos {
		for _, t := range targets {
			for _, dest := range t.Dest {
				destPath := fsutil.ResolvePath(r.Path, dest)
				added, removed, changed, err := fsutil.DiffDir(t.Src, destPath, t.Include, t.Exclude)
				if err != nil {
					a.Out.Err(fmt.Sprintf("%s %s: %v", r.Name, t.Name, err), nil)
					return 1
				}
				results = append(results, diffResult{Repo: r.Name, Target: t.Name, Dest: destPath, Added: added, Removed: removed, Changed: changed})
			}
		}
	}
	if a.Out.JSON {
		a.Out.OK("skills diff", results)
		return 0
	}
	for _, r := range results {
		if len(r.Added)+len(r.Removed)+len(r.Changed) == 0 {
			a.Out.OK(fmt.Sprintf("%s %s clean", r.Repo, r.Target), nil)
			continue
		}
		if len(r.Added) > 0 {
			a.Out.Warn(fmt.Sprintf("%s %s added=%d", r.Repo, r.Target, len(r.Added)), nil)
		}
		if len(r.Removed) > 0 {
			a.Out.Warn(fmt.Sprintf("%s %s removed=%d", r.Repo, r.Target, len(r.Removed)), nil)
		}
		if len(r.Changed) > 0 {
			a.Out.Warn(fmt.Sprintf("%s %s changed=%d", r.Repo, r.Target, len(r.Changed)), nil)
		}
	}
	return 0
}

func (a App) runSkillsVerify(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills verify", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	target := fs.String("target", "", "target")
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
	repos, err := repo.Scan(cfg.ReposRoot)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	repos = repo.Filter(repos, only, exclude)
	var results []verifyResult
	ok := true
	for _, r := range repos {
		for _, t := range targets {
			for _, dest := range t.Dest {
				destPath := fsutil.ResolvePath(r.Path, dest)
				added, removed, changed, err := fsutil.DiffDir(t.Src, destPath, t.Include, t.Exclude)
				if err != nil {
					a.Out.Err(fmt.Sprintf("%s %s: %v", r.Name, t.Name, err), nil)
					return 1
				}
				match := len(added)+len(removed)+len(changed) == 0
				if !match {
					ok = false
				}
				results = append(results, verifyResult{Repo: r.Name, Target: t.Name, Dest: destPath, Match: match})
			}
		}
	}
	if a.Out.JSON {
		a.Out.OK("skills verify", results)
		if ok {
			return 0
		}
		return 2
	}
	for _, r := range results {
		if r.Match {
			a.Out.OK(fmt.Sprintf("%s %s ok", r.Repo, r.Target), nil)
			continue
		}
		a.Out.Err(fmt.Sprintf("%s %s mismatch", r.Repo, r.Target), nil)
	}
	if ok {
		return 0
	}
	return 2
}

func (a App) runSkillsStatus(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills status", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	target := fs.String("target", "", "target")
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
	repos, err := repo.Scan(cfg.ReposRoot)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	repos = repo.Filter(repos, only, exclude)
	var results []diffResult
	for _, r := range repos {
		for _, t := range targets {
			for _, dest := range t.Dest {
				destPath := fsutil.ResolvePath(r.Path, dest)
				added, removed, changed, err := fsutil.DiffDir(t.Src, destPath, t.Include, t.Exclude)
				if err != nil {
					a.Out.Err(fmt.Sprintf("%s %s: %v", r.Name, t.Name, err), nil)
					return 1
				}
				results = append(results, diffResult{Repo: r.Name, Target: t.Name, Dest: destPath, Added: added, Removed: removed, Changed: changed})
			}
		}
	}
	if a.Out.JSON {
		a.Out.OK("skills status", results)
		return 0
	}
	for _, r := range results {
		if len(r.Added)+len(r.Removed)+len(r.Changed) == 0 {
			a.Out.OK(fmt.Sprintf("%s %s clean", r.Repo, r.Target), nil)
			continue
		}
		a.Out.Warn(fmt.Sprintf("%s %s drift", r.Repo, r.Target), nil)
	}
	return 0
}

func (a App) runSkillsClean(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("skills clean", flag.ContinueOnError)
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
	if !*force && cfg.ConflictPolicy != string(fsutil.ConflictOverwrite) {
		a.Out.Err("clean requires --force or conflictPolicy=overwrite", nil)
		return 1
	}
	targets, err := selectTargets(cfg, *target)
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
	guard := guardFromConfig(cfg)
	for _, r := range repos {
		for _, t := range targets {
			files, err := fsutil.ListFiles(t.Src, t.Include, t.Exclude)
			if err != nil {
				a.Out.Err(fmt.Sprintf("%s %s: %v", r.Name, t.Name, err), nil)
				return 1
			}
			for _, dest := range t.Dest {
				destPath := fsutil.ResolvePath(r.Path, dest)
				if err := guard.CheckPath(destPath); err != nil {
					a.Out.Err(err.Error(), nil)
					return 1
				}
				if err := fsutil.CleanDir(destPath, files, *dryRun); err != nil {
					a.Out.Err(fmt.Sprintf("%s %s: %v", r.Name, t.Name, err), nil)
					return 1
				}
			}
		}
	}
	a.Out.OK("skills clean done", nil)
	return 0
}
