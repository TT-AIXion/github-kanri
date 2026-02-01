package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/AIXion-Team/github-kanri/internal/config"
)

func TestSkillsParseErrors(t *testing.T) {
	app, _ := newTestApp(t)
	cases := [][]string{
		{"sync", "--bad"},
		{"link", "--bad"},
		{"diff", "--bad"},
		{"verify", "--bad"},
		{"status", "--bad"},
		{"clean", "--bad"},
		{"pin", "--bad"},
		{"clone", "--bad"},
		{"watch", "--bad"},
	}
	for _, args := range cases {
		if code := app.runSkills(context.Background(), args); code == 0 {
			t.Fatalf("expected parse error: %v", args)
		}
	}
}

func TestSkillsLoadConfigErrors(t *testing.T) {
	app, _ := newTestApp(t)
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("path: %v", err)
	}
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	cases := [][]string{
		{"sync"},
		{"link"},
		{"diff"},
		{"verify"},
		{"status"},
		{"clean"},
		{"watch"},
		{"clone"},
		{"pin", "--target", "skills", "--ref", "HEAD"},
	}
	for _, args := range cases {
		if code := app.runSkills(context.Background(), args); code == 0 {
			t.Fatalf("expected load error: %v", args)
		}
	}
}

func TestSkillsRepoScanErrors(t *testing.T) {
	app, cfg := newTestApp(t)
	cfg.ReposRoot = filepath.Join(t.TempDir(), "missing")
	writeConfig(t, cfg)
	cases := []func() int{
		func() int { return app.runSkillsSync(context.Background(), []string{"--target", "skills"}) },
		func() int { return app.runSkillsDiff(context.Background(), []string{"--target", "skills"}) },
		func() int { return app.runSkillsVerify(context.Background(), []string{"--target", "skills"}) },
		func() int { return app.runSkillsStatus(context.Background(), []string{"--target", "skills"}) },
		func() int { return app.runSkillsClean(context.Background(), []string{"--target", "skills", "--force"}) },
	}
	for _, fn := range cases {
		if code := fn(); code == 0 {
			t.Fatalf("expected scan error")
		}
	}
}

func TestSkillsSelectTargetsErrors(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	if code := app.runSkillsLink(context.Background(), []string{"--target", "missing"}); code == 0 {
		t.Fatalf("expected link target error")
	}
	if code := app.runSkillsWatch(context.Background(), []string{"--target", "missing"}); code == 0 {
		t.Fatalf("expected watch target error")
	}
	if code := app.runSkillsDiff(context.Background(), []string{"--target", "missing"}); code == 0 {
		t.Fatalf("expected diff target error")
	}
	if code := app.runSkillsVerify(context.Background(), []string{"--target", "missing"}); code == 0 {
		t.Fatalf("expected verify target error")
	}
	if code := app.runSkillsStatus(context.Background(), []string{"--target", "missing"}); code == 0 {
		t.Fatalf("expected status target error")
	}
	if code := app.runSkillsClean(context.Background(), []string{"--target", "missing", "--force"}); code == 0 {
		t.Fatalf("expected clean target error")
	}
}

func TestRepoParseErrors(t *testing.T) {
	app, _ := newTestApp(t)
	cases := [][]string{
		{"list", "--bad"},
		{"status", "--bad"},
		{"recent", "--bad"},
		{"open", "--bad"},
		{"path", "--bad"},
		{"info", "--bad"},
		{"graph", "--bad"},
		{"clone", "--bad"},
		{"exec", "--bad"},
	}
	for _, args := range cases {
		if code := app.runRepo(context.Background(), args); code == 0 {
			t.Fatalf("expected parse error: %v", args)
		}
	}
}

func TestRepoLoadConfigErrors(t *testing.T) {
	app, _ := newTestApp(t)
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("path: %v", err)
	}
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	cases := [][]string{
		{"list"},
		{"status"},
		{"recent"},
		{"open", "alpha"},
		{"path", "alpha"},
		{"info", "alpha"},
		{"graph", "alpha"},
		{"clone", "dummy"},
		{"exec", "--cmd", "echo ok"},
	}
	for _, args := range cases {
		if code := app.runRepo(context.Background(), args); code == 0 {
			t.Fatalf("expected load error: %v", args)
		}
	}
}

func TestRepoScanErrors(t *testing.T) {
	app, cfg := newTestApp(t)
	cfg.ReposRoot = filepath.Join(t.TempDir(), "missing")
	writeConfig(t, cfg)
	cases := [][]string{
		{"list"},
		{"status"},
		{"recent"},
		{"open", "alpha"},
		{"path", "alpha"},
		{"info", "alpha"},
		{"graph", "alpha"},
		{"exec", "--cmd", "echo ok"},
	}
	for _, args := range cases {
		if code := app.runRepo(context.Background(), args); code == 0 {
			t.Fatalf("expected scan error: %v", args)
		}
	}
}
