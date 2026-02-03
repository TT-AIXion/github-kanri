package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestRepoCommands(t *testing.T) {
	app, cfg := newTestApp(t)
	repo1 := initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	repo2 := initGitRepo(t, filepath.Join(cfg.ReposRoot, "beta"), false)
	_ = repo2
	if code := app.runRepo(context.Background(), []string{"list"}); code != 0 {
		t.Fatalf("list failed")
	}
	if code := app.runRepo(context.Background(), []string{"status"}); code != 0 {
		t.Fatalf("status failed")
	}
	_ = os.WriteFile(filepath.Join(repo1, "a.txt"), []byte("change"), 0o644)
	if code := app.runRepo(context.Background(), []string{"status"}); code != 0 {
		t.Fatalf("status dirty failed")
	}

	bin := createStubCommand(t, "code")
	prependPath(t, bin)
	if code := app.runRepo(context.Background(), []string{"open", "alpha"}); code != 0 {
		t.Fatalf("open failed")
	}
	if code := app.runRepo(context.Background(), []string{"open", "a"}); code == 0 {
		t.Fatalf("expected multi match")
	}
	cfg.AllowCommands = []string{"git*"}
	writeConfig(t, cfg)
	if code := app.runRepo(context.Background(), []string{"open", "alpha"}); code == 0 {
		t.Fatalf("expected open error")
	}
	cfg.AllowCommands = []string{"*"}
	writeConfig(t, cfg)

	if code := app.runRepo(context.Background(), []string{"path", "alpha"}); code != 0 {
		t.Fatalf("path failed")
	}
	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runRepo(context.Background(), []string{"path", "a"}); code == 0 {
		t.Fatalf("expected multi match json")
	}

	if code := app.runRepo(context.Background(), []string{"recent", "--limit", "1"}); code != 0 {
		t.Fatalf("recent failed")
	}
	if code := app.runRepo(context.Background(), []string{"info", "alpha"}); code != 0 {
		t.Fatalf("info failed")
	}
	if code := app.runRepo(context.Background(), []string{"graph", "alpha", "--limit", "1"}); code != 0 {
		t.Fatalf("graph failed")
	}
	if code := app.runRepo(context.Background(), []string{"graph", "beta"}); code != 0 {
		t.Fatalf("graph no commits failed")
	}
	if code := app.runRepo(context.Background(), []string{"unknown"}); code == 0 {
		t.Fatalf("expected repo error")
	}
}

func TestRepoErrorsAndJSON(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	if code := app.runRepo(context.Background(), []string{"list", "--bad"}); code == 0 {
		t.Fatalf("expected list parse error")
	}
	if code := app.runRepo(context.Background(), []string{"open"}); code == 0 {
		t.Fatalf("expected open error")
	}
	if code := app.runRepo(context.Background(), []string{"open", "missing"}); code == 0 {
		t.Fatalf("expected open no match")
	}
	if code := app.runRepo(context.Background(), []string{"path"}); code == 0 {
		t.Fatalf("expected path error")
	}
	if code := app.runRepo(context.Background(), []string{"info"}); code == 0 {
		t.Fatalf("expected info error")
	}
	if code := app.runRepo(context.Background(), []string{"graph"}); code == 0 {
		t.Fatalf("expected graph error")
	}
	if code := app.runRepo(context.Background(), []string{"clone"}); code == 0 {
		t.Fatalf("expected clone error")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--bad"}); code == 0 {
		t.Fatalf("expected exec parse error")
	}
	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runRepo(context.Background(), []string{"list"}); code != 0 {
		t.Fatalf("json list failed")
	}
	if code := appJSON.runRepo(context.Background(), []string{"status"}); code != 0 {
		t.Fatalf("json status failed")
	}
	if code := appJSON.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok"}); code != 0 {
		t.Fatalf("json exec failed")
	}
}

func TestRepoCloneAndExec(t *testing.T) {
	app, cfg := newTestApp(t)
	repo1 := initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	_ = repo1
	bare := initBareRepo(t, filepath.Join(t.TempDir(), "remote.git"))
	seedBareRepo(t, bare)
	if code := app.runRepo(context.Background(), []string{"clone", bare, "--name", "cloned"}); code != 0 {
		t.Fatalf("clone failed")
	}
	if code := app.runRepo(context.Background(), []string{"clone", bare, "--name", "cloned"}); code == 0 {
		t.Fatalf("expected clone error")
	}
	if code := app.Run(context.Background(), []string{"clone", bare, "--name", "cloned2"}); code != 0 {
		t.Fatalf("clone alias failed")
	}
	if code := app.runRepo(context.Background(), []string{"exec"}); code == 0 {
		t.Fatalf("expected exec error")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok"}); code != 0 {
		t.Fatalf("exec failed")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "exit 1", "--parallel", "2"}); code == 0 {
		t.Fatalf("expected exec error")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok", "--require-clean"}); code != 0 {
		t.Fatalf("exec require-clean failed")
	}
	_ = os.WriteFile(filepath.Join(cfg.ReposRoot, "alpha", "a.txt"), []byte("dirty"), 0o644)
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok", "--require-clean"}); code != 0 {
		t.Fatalf("exec dirty failed")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok", "--timeout", "1"}); code != 0 {
		t.Fatalf("exec timeout path failed")
	}
	cfg.AllowCommands = []string{"git*"}
	writeConfig(t, cfg)
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok", "--require-clean"}); code == 0 {
		t.Fatalf("expected require-clean error")
	}
}
