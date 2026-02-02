package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/TT-AIXion/github-kanri/internal/config"
	"github.com/TT-AIXion/github-kanri/internal/output"
)

func TestConfigCommands(t *testing.T) {
	app, cfg := newTestApp(t)
	if code := app.runConfig(context.Background(), []string{"show"}); code != 0 {
		t.Fatalf("show failed")
	}
	if code := app.runConfig(context.Background(), []string{"validate"}); code != 0 {
		t.Fatalf("validate failed")
	}
	if code := app.runConfig(context.Background(), []string{"unknown"}); code == 0 {
		t.Fatalf("expected error")
	}
	if code := app.runConfig(context.Background(), []string{"init"}); code == 0 {
		t.Fatalf("expected init error")
	}
	cfg.SyncMode = "bad"
	writeConfig(t, cfg)
	if code := app.runConfig(context.Background(), []string{"validate"}); code == 0 {
		t.Fatalf("expected validate error")
	}
	if code := app.runConfig(context.Background(), []string{"init", "--force"}); code != 0 {
		t.Fatalf("init force failed")
	}
	if code := app.runConfig(context.Background(), []string{"show", "--bad"}); code == 0 {
		t.Fatalf("expected parse error")
	}
	if code := app.runConfig(context.Background(), []string{}); code != 0 {
		t.Fatalf("expected help")
	}
	if code := app.runConfig(context.Background(), []string{"--help"}); code != 0 {
		t.Fatalf("expected help")
	}
}

func TestConfigShowErrors(t *testing.T) {
	app, _ := newTestApp(t)
	config.SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	if code := app.runConfig(context.Background(), []string{"show"}); code == 0 {
		t.Fatalf("expected show error")
	}
	config.ResetUserHomeDirForTest()
	path, _ := config.DefaultConfigPath()
	_ = os.Remove(path)
	if code := app.runConfig(context.Background(), []string{"show"}); code == 0 {
		t.Fatalf("expected load error")
	}
}

func TestConfigInitErrors(t *testing.T) {
	app, _ := newTestApp(t)
	config.SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	if code := app.runConfig(context.Background(), []string{"init"}); code == 0 {
		t.Fatalf("expected init error")
	}
	config.ResetUserHomeDirForTest()

	// force save error by read-only home
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, ".config"), 0o555)
	config.SetUserHomeDirForTest(func() (string, error) { return root, nil })
	if code := app.runConfig(context.Background(), []string{"init", "--force"}); code == 0 {
		t.Fatalf("expected save error")
	}
	config.ResetUserHomeDirForTest()
}

func TestConfigValidateErrors(t *testing.T) {
	app, _ := newTestApp(t)
	config.SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	if code := app.runConfig(context.Background(), []string{"validate"}); code == 0 {
		t.Fatalf("expected validate error")
	}
	config.ResetUserHomeDirForTest()

	path, _ := config.DefaultConfigPath()
	_ = os.WriteFile(path, []byte("{bad"), 0o644)
	if code := app.runConfig(context.Background(), []string{"validate"}); code == 0 {
		t.Fatalf("expected load error")
	}
}

func TestConfigInitParseError(t *testing.T) {
	app, _ := newTestApp(t)
	if code := app.runConfigInit(context.Background(), []string{"--bad"}); code == 0 {
		t.Fatalf("expected init parse error")
	}
}

func TestConfigValidateParseError(t *testing.T) {
	app, _ := newTestApp(t)
	if code := app.runConfigValidate(context.Background(), []string{"--bad"}); code == 0 {
		t.Fatalf("expected validate parse error")
	}
}

func TestConfigInitDefaultConfigError(t *testing.T) {
	app, _ := newTestApp(t)
	orig := defaultConfig
	defaultConfig = func() (config.Config, error) { return config.Config{}, os.ErrInvalid }
	defer func() { defaultConfig = orig }()
	if code := app.runConfigInit(context.Background(), []string{"--force"}); code == 0 {
		t.Fatalf("expected default config error")
	}
}

func TestConfigValidateExpandError(t *testing.T) {
	home := t.TempDir()
	config.SetUserHomeDirForTest(func() (string, error) { return home, nil })
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("path error: %v", err)
	}
	cfg := config.Config{
		ProjectsRoot: "~/Projects",
		ReposRoot:    "~/Projects/repos",
		SkillsRoot:   "~/Projects/skills",
		SkillTargets: []string{".codex/skills"},
		SyncTargets: []config.SyncTarget{{
			Name: "skills",
			Src:  "~/Projects/skills",
			Dest: []string{".codex/skills"},
		}},
		SyncMode:       "copy",
		ConflictPolicy: "fail",
		DenyCommands:   []string{"rm -rf*"},
	}
	if err := config.Save(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}
	calls := 0
	config.SetUserHomeDirForTest(func() (string, error) {
		calls++
		if calls == 1 {
			return home, nil
		}
		return "", os.ErrPermission
	})
	defer config.ResetUserHomeDirForTest()
	app := App{Out: output.Writer{Out: os.Stdout, ErrW: os.Stderr}}
	if code := app.runConfigValidate(context.Background(), nil); code == 0 {
		t.Fatalf("expected expand error")
	}
}
