package app

import (
	"fmt"

	"github.com/AIXion-Team/github-kanri/internal/config"
	"github.com/AIXion-Team/github-kanri/internal/executil"
	"github.com/AIXion-Team/github-kanri/internal/safety"
)

func loadConfig() (config.Config, string, error) {
	path, err := config.DefaultConfigPath()
	if err != nil {
		return config.Config{}, "", err
	}
	cfg, err := config.Load(path)
	if err != nil {
		return config.Config{}, "", fmt.Errorf("config load failed: %w", err)
	}
	cfg, err = config.ExpandConfigPaths(cfg)
	if err != nil {
		return config.Config{}, "", err
	}
	return cfg, path, nil
}

func buildRunner(cfg config.Config, dryRun bool) executil.Runner {
	return executil.Runner{
		Guard: safety.Guard{
			AllowCommands: cfg.AllowCommands,
			DenyCommands:  cfg.DenyCommands,
			AllowPaths:    cfg.AllowPaths,
			DenyPaths:     cfg.DenyPaths,
		},
		DryRun: dryRun,
	}
}
