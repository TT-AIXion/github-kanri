package app

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/TT-AIXion/github-kanri/internal/output"
)

func TestSkillsClonePullError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	cfg.SkillsRemote = "dummy"
	cfg.AllowCommands = []string{"git status*"}
	writeConfig(t, cfg)
	if code := app.runSkillsClone(context.Background(), []string{"--remote", cfg.SkillsRemote}); code == 0 {
		t.Fatalf("expected pull error")
	}
}

func TestSkillsCloneRemoveAllError(t *testing.T) {
	app, cfg := newTestApp(t)
	cfg.SkillsRemote = "dummy"
	writeConfig(t, cfg)
	orig := removeAll
	removeAll = func(string) error { return errors.New("remove") }
	defer func() { removeAll = orig }()
	if code := app.runSkillsClone(context.Background(), []string{"--remote", cfg.SkillsRemote, "--force"}); code == 0 {
		t.Fatalf("expected remove error")
	}
}

func TestSkillsCloneCloneError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = os.RemoveAll(cfg.SkillsRoot)
	cfg.SkillsRemote = "dummy"
	cfg.AllowCommands = []string{"git status*"}
	writeConfig(t, cfg)
	if code := app.runSkillsClone(context.Background(), []string{"--remote", cfg.SkillsRemote}); code == 0 {
		t.Fatalf("expected clone error")
	}
}

func TestSkillsPinDirtyNoForce(t *testing.T) {
	app, cfg := newTestApp(t)
	repo := initGitRepo(t, cfg.SkillsRoot, true)
	_ = os.WriteFile(filepath.Join(repo, "dirty.txt"), []byte("dirty"), 0o644)
	writeConfig(t, cfg)
	if code := app.runSkillsPin(context.Background(), []string{"--target", "skills", "--ref", "HEAD"}); code == 0 {
		t.Fatalf("expected dirty error")
	}
}

func TestSkillsPinFetchError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	cfg.AllowCommands = []string{"git status*"}
	writeConfig(t, cfg)
	if code := app.runSkillsPin(context.Background(), []string{"--target", "skills", "--ref", "HEAD"}); code == 0 {
		t.Fatalf("expected fetch error")
	}
}

func TestSkillsPinIsCleanError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	cfg.AllowCommands = []string{"git fetch*"}
	writeConfig(t, cfg)
	if code := app.runSkillsPin(context.Background(), []string{"--target", "skills", "--ref", "HEAD"}); code == 0 {
		t.Fatalf("expected is-clean error")
	}
}

func TestSkillsPinCheckoutError(t *testing.T) {
	app, cfg := newTestApp(t)
	bare := initBareRepo(t, filepath.Join(cfg.ProjectsRoot, "skills-remote.git"))
	seedBareRepo(t, bare)
	_ = os.RemoveAll(cfg.SkillsRoot)
	_ = runGit(cfg.ProjectsRoot, "clone", bare, cfg.SkillsRoot)
	cfg.AllowCommands = []string{"git status*", "git fetch*"}
	writeConfig(t, cfg)
	if code := app.runSkillsPin(context.Background(), []string{"--target", "skills", "--ref", "HEAD"}); code == 0 {
		t.Fatalf("expected checkout error")
	}
}

func TestSkillsStateErrors(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = app
	cfg.SkillsRoot = filepath.Join(t.TempDir(), "missing")
	if _, err := skillsState(context.Background(), cfg); err == nil {
		t.Fatalf("expected stat error")
	}

	cfg2 := cfg
	cfg2.SkillsRoot = t.TempDir()
	_ = initGitRepo(t, cfg2.SkillsRoot, true)
	cfg2.AllowCommands = []string{"git status*"}
	if _, err := skillsState(context.Background(), cfg2); err == nil {
		t.Fatalf("expected rev-parse error")
	}

	cfg3 := cfg2
	cfg3.AllowCommands = []string{"git rev-parse*"}
	if _, err := skillsState(context.Background(), cfg3); err == nil {
		t.Fatalf("expected status error")
	}
}

func TestSkillsWatchIntervalReset(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	writeConfig(t, cfg)
	orig := timeAfter
	defer func() { timeAfter = orig }()
	got := make(chan time.Duration, 1)
	timeAfter = func(d time.Duration) <-chan time.Time {
		got <- d
		return make(chan time.Time)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() { done <- app.runSkillsWatch(ctx, []string{"--interval", "0"}) }()
	if d := <-got; d != 5*time.Second {
		t.Fatalf("expected default interval")
	}
	cancel()
	<-done
}

func TestSkillsWatchSyncOnChange(t *testing.T) {
	app, cfg := newTestApp(t)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	app.Out = output.Writer{Out: &out, ErrW: &errBuf}
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	writeConfig(t, cfg)
	orig := timeAfter
	defer func() { timeAfter = orig }()
	tick := make(chan time.Time, 1)
	started := make(chan struct{})
	timeAfter = func(time.Duration) <-chan time.Time {
		select {
		case <-started:
		default:
			close(started)
		}
		return tick
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan int, 1)
	go func() { done <- app.runSkillsWatch(ctx, []string{"--interval", "1"}) }()
	<-started
	_ = os.WriteFile(filepath.Join(cfg.SkillsRoot, "new.txt"), []byte("x"), 0o644)
	tick <- time.Now()
	time.Sleep(100 * time.Millisecond)
	cancel()
	if code := <-done; code != 0 {
		t.Fatalf("expected watch ok: %s", errBuf.String())
	}
}

func TestSkillsWatchStateError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	writeConfig(t, cfg)
	orig := timeAfter
	defer func() { timeAfter = orig }()
	tick := make(chan time.Time, 1)
	started := make(chan struct{})
	timeAfter = func(time.Duration) <-chan time.Time {
		select {
		case <-started:
		default:
			close(started)
		}
		return tick
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan int, 1)
	go func() { done <- app.runSkillsWatch(ctx, []string{"--interval", "1"}) }()
	<-started
	_ = os.RemoveAll(cfg.SkillsRoot)
	tick <- time.Now()
	if code := <-done; code == 0 {
		t.Fatalf("expected watch error")
	}
}

func TestSkillsWatchSyncError(t *testing.T) {
	app, cfg := newTestApp(t)
	_ = initGitRepo(t, cfg.SkillsRoot, true)
	_ = initGitRepo(t, filepath.Join(cfg.ReposRoot, "alpha"), true)
	cfg.DenyPaths = []string{"**"}
	writeConfig(t, cfg)
	orig := timeAfter
	defer func() { timeAfter = orig }()
	tick := make(chan time.Time, 1)
	started := make(chan struct{})
	timeAfter = func(time.Duration) <-chan time.Time {
		select {
		case <-started:
		default:
			close(started)
		}
		return tick
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	done := make(chan int, 1)
	go func() { done <- app.runSkillsWatch(ctx, []string{"--interval", "1"}) }()
	<-started
	_ = os.WriteFile(filepath.Join(cfg.SkillsRoot, "new.txt"), []byte("x"), 0o644)
	tick <- time.Now()
	if code := <-done; code == 0 {
		t.Fatalf("expected watch sync error")
	}
}
