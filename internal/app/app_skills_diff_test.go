package app

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/AIXion-Team/github-kanri/internal/config"
)

func TestSkillsDiffStatusVerify(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = os.WriteFile(filepath.Join(cfg.SkillsRoot, "a.txt"), []byte("a"), 0o644)
	repo1 := initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	repo2 := initGitRepo(t, filepath.Join(cfg.ReposRoot, "beta"), true)
	repo3 := initGitRepo(t, filepath.Join(cfg.ReposRoot, "gamma"), true)
	_ = os.MkdirAll(filepath.Join(repo1, ".codex", "skills"), 0o755)
	_ = os.WriteFile(filepath.Join(repo1, ".codex", "skills", "a.txt"), []byte("a"), 0o644)
	_ = os.MkdirAll(filepath.Join(repo2, ".codex", "skills"), 0o755)
	_ = os.WriteFile(filepath.Join(repo2, ".codex", "skills", "b.txt"), []byte("b"), 0o644)
	_ = os.MkdirAll(filepath.Join(repo3, ".codex", "skills"), 0o755)
	_ = os.WriteFile(filepath.Join(repo3, ".codex", "skills", "a.txt"), []byte("changed"), 0o644)
	writeConfig(t, cfg)

	if code := app.runSkillsDiff(context.Background(), []string{"--target", "skills"}); code != 0 {
		t.Fatalf("diff failed")
	}
	if code := app.runSkillsStatus(context.Background(), []string{"--target", "skills"}); code != 0 {
		t.Fatalf("status failed")
	}
	if code := app.runSkillsVerify(context.Background(), []string{"--target", "skills"}); code != 2 {
		t.Fatalf("expected verify mismatch")
	}

	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runSkillsVerify(context.Background(), []string{"--target", "skills"}); code != 2 {
		t.Fatalf("expected verify mismatch json")
	}
}

func TestSkillsDiffError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.SyncTargets = []config.SyncTarget{{
		Name: "skills",
		Src:  filepath.Join(cfg.ProjectsRoot, "missing"),
		Dest: []string{".codex/skills"},
	}}
	writeConfig(t, cfg)
	if code := app.runSkillsDiff(context.Background(), []string{"--target", "skills"}); code == 0 {
		t.Fatalf("expected diff error")
	}
}

func TestSkillsStatusError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.SyncTargets = []config.SyncTarget{{
		Name: "skills",
		Src:  filepath.Join(cfg.ProjectsRoot, "missing"),
		Dest: []string{".codex/skills"},
	}}
	writeConfig(t, cfg)
	if code := app.runSkillsStatus(context.Background(), []string{"--target", "skills"}); code == 0 {
		t.Fatalf("expected status error")
	}
}

func TestSkillsVerifyError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.SyncTargets = []config.SyncTarget{{
		Name: "skills",
		Src:  filepath.Join(cfg.ProjectsRoot, "missing"),
		Dest: []string{".codex/skills"},
	}}
	writeConfig(t, cfg)
	if code := app.runSkillsVerify(context.Background(), []string{"--target", "skills"}); code == 0 {
		t.Fatalf("expected verify error")
	}
}

func TestSkillsDiffChanged(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = os.WriteFile(filepath.Join(cfg.SkillsRoot, "a.txt"), []byte("a"), 0o644)
	repo := initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	_ = os.MkdirAll(filepath.Join(repo, ".codex", "skills"), 0o755)
	_ = os.WriteFile(filepath.Join(repo, ".codex", "skills", "a.txt"), []byte("changed"), 0o644)
	writeConfig(t, cfg)
	if code := app.runSkillsDiff(context.Background(), []string{"--target", "skills"}); code != 0 {
		t.Fatalf("expected diff changed")
	}
}

func TestSkillsCleanErrors(t *testing.T) {
	t.Run("listfiles", func(t *testing.T) {
		app, cfg := newTestApp(t)
		_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
		cfg.SyncTargets = []config.SyncTarget{{
			Name: "skills",
			Src:  filepath.Join(cfg.ProjectsRoot, "missing"),
			Dest: []string{".codex/skills"},
		}}
		writeConfig(t, cfg)
		if code := app.runSkillsClean(context.Background(), []string{"--target", "skills", "--force"}); code == 0 {
			t.Fatalf("expected clean listfiles error")
		}
	})

	t.Run("deny", func(t *testing.T) {
		app, cfg := newTestApp(t)
		_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
		cfg.DenyPaths = []string{"**"}
		writeConfig(t, cfg)
		if code := app.runSkillsClean(context.Background(), []string{"--target", "skills", "--force"}); code == 0 {
			t.Fatalf("expected deny error")
		}
	})

	t.Run("missing-dest", func(t *testing.T) {
		app, cfg := newTestApp(t)
		_ = os.WriteFile(filepath.Join(cfg.SkillsRoot, "a.txt"), []byte("a"), 0o644)
		_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
		cfg.SyncTargets = []config.SyncTarget{{
			Name: "skills",
			Src:  cfg.SkillsRoot,
			Dest: []string{".codex/missing"},
		}}
		writeConfig(t, cfg)
		if code := app.runSkillsClean(context.Background(), []string{"--target", "skills", "--force"}); code == 0 {
			t.Fatalf("expected clean error")
		}
	})
}

func TestSkillsSyncErrors(t *testing.T) {
	t.Run("deny", func(t *testing.T) {
		app, cfg := newTestApp(t)
		_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
		cfg.DenyPaths = []string{"**"}
		writeConfig(t, cfg)
		if code := app.runSkillsSync(context.Background(), []string{"--target", "skills"}); code == 0 {
			t.Fatalf("expected sync deny error")
		}
	})

	t.Run("syncdir", func(t *testing.T) {
		app, cfg := newTestApp(t)
		_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
		cfg.SyncTargets = []config.SyncTarget{{
			Name: "skills",
			Src:  filepath.Join(cfg.ProjectsRoot, "missing"),
			Dest: []string{".codex/skills"},
		}}
		writeConfig(t, cfg)
		if code := app.runSkillsSync(context.Background(), []string{"--target", "skills"}); code == 0 {
			t.Fatalf("expected sync error")
		}
	})
}

func TestSkillsSyncJSONAndDenyDest(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runSkillsSync(context.Background(), []string{"--target", "skills"}); code != 0 {
		t.Fatalf("expected sync json")
	}

	cfg.DenyPaths = []string{"**/.codex/skills"}
	writeConfig(t, cfg)
	if code := app.runSkillsSync(context.Background(), []string{"--target", "skills"}); code == 0 {
		t.Fatalf("expected deny dest error")
	}
}

func TestSkillsSyncForce(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.ConflictPolicy = "fail"
	writeConfig(t, cfg)
	if code := app.runSkillsSync(context.Background(), []string{"--target", "skills", "--force"}); code != 0 {
		t.Fatalf("expected sync force")
	}
}
