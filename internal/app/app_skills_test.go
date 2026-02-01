package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSkillsCommands(t *testing.T) {
	app, cfg := newTestApp(t)
	if code := app.runSkills(context.Background(), []string{}); code != 0 {
		t.Fatalf("expected help")
	}
	if code := app.runSkills(context.Background(), []string{"--help"}); code != 0 {
		t.Fatalf("expected help")
	}
	if code := app.runSkills(context.Background(), []string{"clone"}); code == 0 {
		t.Fatalf("expected clone error")
	}

	bare := initBareRepo(t, filepath.Join(cfg.ProjectsRoot, "skills-remote.git"))
	seedBareRepo(t, bare)
	cfg.SkillsRemote = bare
	writeConfig(t, cfg)
	_ = os.RemoveAll(cfg.SkillsRoot)
	if code := app.runSkills(context.Background(), []string{"clone"}); code != 0 {
		t.Fatalf("clone failed")
	}
	if code := app.runSkills(context.Background(), []string{"clone"}); code != 0 {
		t.Fatalf("clone pull failed")
	}
	_ = os.WriteFile(filepath.Join(cfg.SkillsRoot, "tool.txt"), []byte("x"), 0o644)

	_ = os.RemoveAll(cfg.SkillsRoot)
	_ = os.MkdirAll(cfg.SkillsRoot, 0o755)
	writeConfig(t, cfg)
	if code := app.runSkills(context.Background(), []string{"clone"}); code == 0 {
		t.Fatalf("expected clone error")
	}
	if code := app.runSkills(context.Background(), []string{"clone", "--force"}); code != 0 {
		t.Fatalf("clone force failed")
	}

	if code := app.runSkills(context.Background(), []string{"sync", "--target", "skills"}); code != 0 {
		t.Fatalf("sync failed")
	}
	if code := app.runSkills(context.Background(), []string{"link", "--target", "skills"}); code != 0 {
		t.Fatalf("link failed")
	}
	if code := app.runSkills(context.Background(), []string{"sync", "--mode", "bad"}); code == 0 {
		t.Fatalf("expected mode error")
	}
	if code := app.runSkills(context.Background(), []string{"sync", "--target", "missing"}); code == 0 {
		t.Fatalf("expected target error")
	}

	if code := app.runSkills(context.Background(), []string{"diff", "--target", "skills"}); code != 0 {
		t.Fatalf("diff failed")
	}
	if code := app.runSkills(context.Background(), []string{"verify", "--target", "skills"}); code != 0 {
		t.Fatalf("verify failed")
	}
	if code := app.runSkills(context.Background(), []string{"status", "--target", "skills"}); code != 0 {
		t.Fatalf("status failed")
	}

	if code := app.runSkills(context.Background(), []string{"pin"}); code == 0 {
		t.Fatalf("expected pin error")
	}
	if code := app.runSkills(context.Background(), []string{"pin", "--target", "missing", "--ref", "HEAD"}); code == 0 {
		t.Fatalf("expected target error")
	}
	if code := app.runSkills(context.Background(), []string{"pin", "--target", "skills", "--ref", "HEAD"}); code != 0 {
		t.Fatalf("pin failed")
	}

	cfg.ConflictPolicy = "fail"
	writeConfig(t, cfg)
	if code := app.runSkills(context.Background(), []string{"clean"}); code == 0 {
		t.Fatalf("expected clean error")
	}
	if code := app.runSkills(context.Background(), []string{"clean", "--force"}); code != 0 {
		t.Fatalf("clean failed")
	}
}

func TestSkillsErrorsAndJSON(t *testing.T) {
	app, cfg := newTestApp(t)
	if code := app.runSkills(context.Background(), []string{"unknown"}); code == 0 {
		t.Fatalf("expected skills error")
	}
	if code := app.runSkills(context.Background(), []string{"sync", "--bad"}); code == 0 {
		t.Fatalf("expected sync parse error")
	}
	appJSON := app
	appJSON.Out.JSON = true
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	writeConfig(t, cfg)
	if code := appJSON.runSkills(context.Background(), []string{"status"}); code != 0 {
		t.Fatalf("json status failed")
	}
	if code := appJSON.runSkills(context.Background(), []string{"diff"}); code != 0 {
		t.Fatalf("json diff failed")
	}
	if code := appJSON.runSkills(context.Background(), []string{"verify"}); code != 0 {
		t.Fatalf("json verify failed")
	}
}

func TestSkillsWatch(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	writeConfig(t, cfg)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	go func() {
		_ = app.runSkillsWatch(ctx, []string{"--interval", "1"})
	}()
	time.Sleep(30 * time.Millisecond)
}

func TestSkillsWatchError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = os.RemoveAll(cfg.SkillsRoot)
	writeConfig(t, cfg)
	if code := app.runSkillsWatch(context.Background(), []string{"--interval", "1"}); code == 0 {
		t.Fatalf("expected watch error")
	}
}
