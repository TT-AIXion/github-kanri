package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/AIXion-Team/github-kanri/internal/config"
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
