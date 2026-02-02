package app

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/TT-AIXion/github-kanri/internal/executil"
)

func TestRepoListScanError(t *testing.T) {
	app, cfg := newTestApp(t)
	cfg.ReposRoot = filepath.Join(t.TempDir(), "missing")
	writeConfig(t, cfg)
	if code := app.runRepoList(context.Background(), nil); code == 0 {
		t.Fatalf("expected list error")
	}
}

func TestRepoStatusCommandError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.AllowCommands = []string{"git log*"}
	writeConfig(t, cfg)
	if code := app.runRepoStatus(context.Background(), nil); code == 0 {
		t.Fatalf("expected status error")
	}
}

func TestRepoRecentJSON(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "beta"), false)
	writeConfig(t, cfg)
	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runRepoRecent(context.Background(), []string{"--limit", "10"}); code != 0 {
		t.Fatalf("expected recent json")
	}
}

func TestRepoPathJSON(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runRepoPath(context.Background(), []string{"alpha"}); code != 0 {
		t.Fatalf("expected path json")
	}
}

func TestRepoPathNoMatch(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	if code := app.runRepoPath(context.Background(), []string{"missing"}); code == 0 {
		t.Fatalf("expected path no match")
	}
}

func TestRepoGraphJSONAndEmpty(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)

	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runRepoGraph(context.Background(), []string{"alpha"}); code != 0 {
		t.Fatalf("expected graph json")
	}

	orig := logOneline
	logOneline = func(context.Context, executil.Runner, string, int) (string, error) {
		return "", nil
	}
	defer func() { logOneline = orig }()
	if code := app.runRepoGraph(context.Background(), []string{"alpha"}); code != 0 {
		t.Fatalf("expected graph empty ok")
	}

	logOneline = func(context.Context, executil.Runner, string, int) (string, error) {
		return "a\n\nb", nil
	}
	if code := app.runRepoGraph(context.Background(), []string{"alpha"}); code != 0 {
		t.Fatalf("expected graph skip blank")
	}
}

func TestRepoInfoJSON(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runRepoInfo(context.Background(), []string{"alpha"}); code != 0 {
		t.Fatalf("expected info json")
	}
}

func TestRepoInfoParseAndMultiMatch(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alps"), true)
	writeConfig(t, cfg)
	if code := app.runRepoInfo(context.Background(), []string{"--bad"}); code == 0 {
		t.Fatalf("expected info parse error")
	}
	if code := app.runRepoInfo(context.Background(), []string{"a"}); code == 0 {
		t.Fatalf("expected info multi match")
	}
	if code := app.runRepoInfo(context.Background(), []string{"missing"}); code == 0 {
		t.Fatalf("expected info no match")
	}
}

func TestRepoCloneDefaultName(t *testing.T) {
	app, cfg := newTestApp(t)
	bare := initBareRepo(t, filepath.Join(t.TempDir(), "remote.git"))
	seedBareRepo(t, bare)
	writeConfig(t, cfg)
	if code := app.runRepoClone(context.Background(), []string{bare}); code != 0 {
		t.Fatalf("expected clone")
	}
}

func TestRepoCloneError(t *testing.T) {
	app, cfg := newTestApp(t)
	bare := initBareRepo(t, filepath.Join(t.TempDir(), "remote.git"))
	seedBareRepo(t, bare)
	cfg.AllowCommands = []string{"git status*"}
	writeConfig(t, cfg)
	if code := app.runRepoClone(context.Background(), []string{bare}); code == 0 {
		t.Fatalf("expected clone error")
	}
}

func TestRepoExecJSONErrorAndParallel(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.AllowCommands = []string{"git*"}
	writeConfig(t, cfg)
	appJSON := app
	appJSON.Out.JSON = true
	if code := appJSON.runRepoExec(context.Background(), []string{"--cmd", "echo ok"}); code == 0 {
		t.Fatalf("expected exec json error")
	}

	cfg.AllowCommands = []string{"*"}
	writeConfig(t, cfg)
	if code := app.runRepoExec(context.Background(), []string{"--cmd", "echo ok", "--parallel", "0"}); code != 0 {
		t.Fatalf("expected exec parallel reset")
	}
}

func TestRepoHelpFlag(t *testing.T) {
	app, _ := newTestApp(t)
	if code := app.runRepo(context.Background(), []string{"--help"}); code != 0 {
		t.Fatalf("expected repo help")
	}
}

func TestRepoHelpKeyword(t *testing.T) {
	app, _ := newTestApp(t)
	if code := app.runRepo(context.Background(), []string{"help"}); code != 0 {
		t.Fatalf("expected repo help keyword")
	}
}

func TestRepoExecStderr(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	if code := app.runRepoExec(context.Background(), []string{"--cmd", "echo err 1>&2"}); code != 0 {
		t.Fatalf("expected exec stderr")
	}
}

func TestRepoGraphMultiMatch(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alps"), true)
	writeConfig(t, cfg)
	if code := app.runRepoGraph(context.Background(), []string{"a"}); code == 0 {
		t.Fatalf("expected graph multi match")
	}
}

func TestRepoGraphNoMatch(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	if code := app.runRepoGraph(context.Background(), []string{"missing"}); code == 0 {
		t.Fatalf("expected graph no match")
	}
}

func TestRepoRecentLimit(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "beta"), true)
	writeConfig(t, cfg)
	if code := app.runRepoRecent(context.Background(), []string{"--limit", "1"}); code != 0 {
		t.Fatalf("expected recent limit")
	}
}

func TestRepoRecentNoCommitOutput(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "beta"), false)
	cfg.AllowCommands = []string{"git status*"}
	writeConfig(t, cfg)
	if code := app.runRepoRecent(context.Background(), []string{"--limit", "10"}); code != 0 {
		t.Fatalf("expected recent output")
	}
}

func TestRepoExecRequireCleanError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.AllowCommands = []string{"git log*"}
	writeConfig(t, cfg)
	if code := app.runRepoExec(context.Background(), []string{"--cmd", "echo ok", "--require-clean"}); code == 0 {
		t.Fatalf("expected require-clean error")
	}
}
