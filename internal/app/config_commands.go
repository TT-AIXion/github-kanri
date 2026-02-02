package app

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/TT-AIXion/github-kanri/internal/config"
)

var defaultConfig = config.DefaultConfig

func (a App) runConfig(ctx context.Context, args []string) int {
	if len(args) == 0 || isHelp(args[0]) {
		a.Out.Raw(`gkn config <command>

Commands:
  show       show config
  init       create default config
  validate   validate config`)
		return 0
	}
	switch args[0] {
	case "show":
		return a.runConfigShow(ctx, args[1:])
	case "init":
		return a.runConfigInit(ctx, args[1:])
	case "validate":
		return a.runConfigValidate(ctx, args[1:])
	default:
		a.Out.Err(fmt.Sprintf("unknown config command: %s", args[0]), nil)
		return 1
	}
}

func (a App) runConfigShow(_ context.Context, args []string) int {
	fs := flag.NewFlagSet("config show", flag.ContinueOnError)
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
		a.Out.Err(err.Error(), nil)
		return 1
	}
	a.Out.OK(path, cfg)
	return 0
}

func (a App) runConfigInit(_ context.Context, args []string) int {
	fs := flag.NewFlagSet("config init", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	force := fs.Bool("force", false, "overwrite existing config")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	path, err := config.DefaultConfigPath()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if _, err := os.Stat(path); err == nil && !*force {
		a.Out.Err("config already exists (use --force)", nil)
		return 1
	}
	cfg, err := defaultConfig()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := config.Save(path, cfg); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	a.Out.OK("config initialized", map[string]string{"path": path})
	return 0
}

func (a App) runConfigValidate(_ context.Context, args []string) int {
	fs := flag.NewFlagSet("config validate", flag.ContinueOnError)
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
		a.Out.Err(err.Error(), nil)
		return 1
	}
	cfg, err = config.ExpandConfigPaths(cfg)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	errs := config.Validate(cfg)
	if len(errs) == 0 {
		a.Out.OK("config valid", map[string]string{"path": path})
		return 0
	}
	for _, e := range errs {
		a.Out.Err(e.Error(), nil)
	}
	return 1
}

func isHelp(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}
