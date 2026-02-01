package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"

	"github.com/AIXion-Team/github-kanri/internal/config"
)

func (a App) runDoctor(_ context.Context, args []string) int {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	path, err := config.DefaultConfigPath()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	cfg, err := config.Load(path)
	if err != nil {
		a.Out.Err(fmt.Sprintf("config load failed: %v", err), nil)
		return 1
	}
	cfg, err = config.ExpandConfigPaths(cfg)
	if err != nil {
		a.Out.Err(fmt.Sprintf("config expand failed: %v", err), nil)
		return 1
	}
	if errs := config.Validate(cfg); len(errs) > 0 {
		for _, e := range errs {
			a.Out.Err(e.Error(), nil)
		}
		return 1
	}
	if _, err := exec.LookPath("git"); err != nil {
		a.Out.Err("git not found", nil)
		return 1
	}
	if _, err := os.Stat(cfg.ReposRoot); err != nil {
		a.Out.Warn(fmt.Sprintf("reposRoot not found: %s", cfg.ReposRoot), nil)
	}
	if _, err := os.Stat(cfg.SkillsRoot); err != nil {
		a.Out.Warn(fmt.Sprintf("skillsRoot not found: %s", cfg.SkillsRoot), nil)
	}
	a.Out.OK("doctor ok", map[string]string{"config": path})
	return 0
}
