package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/AIXion-Team/github-kanri/internal/gitutil"
	"github.com/AIXion-Team/github-kanri/internal/repo"
)

func (a App) runRepoExec(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("repo exec", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	cmd := fs.String("cmd", "", "command")
	parallel := fs.Int("parallel", 1, "parallelism")
	timeout := fs.Int("timeout", 0, "timeout seconds")
	requireClean := fs.Bool("require-clean", false, "require clean repo")
	dryRun := fs.Bool("dry-run", false, "dry run")
	var only multiFlag
	var exclude multiFlag
	fs.Var(&only, "only", "only patterns")
	fs.Var(&exclude, "exclude", "exclude patterns")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if strings.TrimSpace(*cmd) == "" {
		a.Out.Err("--cmd required", nil)
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
	runner := buildRunner(cfg, *dryRun)

	results := make([]execResult, len(repos))
	if *parallel < 1 {
		*parallel = 1
	}
	jobs := make(chan int)
	var wg sync.WaitGroup
	worker := func() {
		defer wg.Done()
		for idx := range jobs {
			r := repos[idx]
			if *requireClean {
				clean, err := gitutil.IsClean(ctx, runner, r.Path)
				if err != nil {
					results[idx] = execResult{Name: r.Name, Path: r.Path, Error: err.Error(), ExitCode: 1}
					continue
				}
				if !clean {
					results[idx] = execResult{Name: r.Name, Path: r.Path, Error: "dirty", ExitCode: 0, Skipped: true}
					continue
				}
			}
			cmdCtx := ctx
			var cancel context.CancelFunc
			if *timeout > 0 {
				cmdCtx, cancel = context.WithTimeout(ctx, time.Duration(*timeout)*time.Second)
			}
			res, err := runner.RunShell(cmdCtx, r.Path, *cmd)
			if cancel != nil {
				cancel()
			}
			result := execResult{
				Name:     r.Name,
				Path:     r.Path,
				ExitCode: res.ExitCode,
				Duration: res.Duration.String(),
				Stdout:   strings.TrimSpace(res.Stdout),
				Stderr:   strings.TrimSpace(res.Stderr),
			}
			if err != nil {
				result.Error = err.Error()
			}
			results[idx] = result
		}
	}
	for i := 0; i < *parallel; i++ {
		wg.Add(1)
		go worker()
	}
	for i := range repos {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	if a.Out.JSON {
		a.Out.OK("repo exec", results)
		return 0
	}
	for _, r := range results {
		if r.Skipped {
			a.Out.Warn(fmt.Sprintf("%s skipped (dirty)", r.Name), nil)
			continue
		}
		if r.Error != "" {
			a.Out.Err(fmt.Sprintf("%s %s", r.Name, r.Error), nil)
			continue
		}
		a.Out.OK(fmt.Sprintf("%s exit=%d", r.Name, r.ExitCode), nil)
		if r.Stdout != "" {
			a.Out.Raw(r.Stdout)
		}
		if r.Stderr != "" {
			a.Out.Raw(r.Stderr)
		}
	}
	return 0
}
