package app

import (
	"bytes"
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/AIXion-Team/github-kanri/internal/config"
	"github.com/AIXion-Team/github-kanri/internal/output"
)

func TestAppRunBasic(t *testing.T) {
	app, _ := newTestApp(t)
	if code := app.Run(context.Background(), []string{}); code != 0 {
		t.Fatalf("expected 0")
	}
	if code := app.Run(context.Background(), []string{"help"}); code != 0 {
		t.Fatalf("expected 0")
	}
	if code := app.Run(context.Background(), []string{"version"}); code != 0 {
		t.Fatalf("expected 0")
	}
	if code := app.Run(context.Background(), []string{"repo", "list"}); code != 0 {
		t.Fatalf("expected repo list")
	}
	if code := app.Run(context.Background(), []string{"skills", "status"}); code != 0 {
		t.Fatalf("expected skills status")
	}
	if code := app.Run(context.Background(), []string{"config", "show"}); code != 0 {
		t.Fatalf("expected config show")
	}
	if code := app.Run(context.Background(), []string{"doctor"}); code != 0 {
		t.Fatalf("expected doctor")
	}
	if code := app.Run(context.Background(), []string{"unknown"}); code == 0 {
		t.Fatalf("expected error")
	}
}

func TestExitError(t *testing.T) {
	err := ExitError{Code: 1, Err: os.ErrPermission}
	if err.Error() == "" {
		t.Fatalf("expected error string")
	}
}

func TestMultiFlag(t *testing.T) {
	var m multiFlag
	_ = m.Set("a,b")
	_ = m.Set(" ")
	if m.String() == "" {
		t.Fatalf("expected values")
	}
}

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
	// init should fail without force when config exists
	if code := app.runConfig(context.Background(), []string{"init"}); code == 0 {
		t.Fatalf("expected init error")
	}
	// invalid config for validate
	cfg.SyncMode = "bad"
	writeConfig(t, cfg)
	if code := app.runConfig(context.Background(), []string{"validate"}); code == 0 {
		t.Fatalf("expected validate error")
	}
	// init with force
	if code := app.runConfig(context.Background(), []string{"init", "--force"}); code != 0 {
		t.Fatalf("init force failed")
	}
	// parse error
	if code := app.runConfig(context.Background(), []string{"show", "--bad"}); code == 0 {
		t.Fatalf("expected parse error")
	}
}

func TestLoadConfigErrors(t *testing.T) {
	// invalid home
	config.SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	if _, _, err := loadConfig(); err == nil {
		t.Fatalf("expected load error")
	}
	config.ResetUserHomeDirForTest()

	// invalid json
	app, _ := newTestApp(t)
	path, _ := config.DefaultConfigPath()
	_ = os.WriteFile(path, []byte("{bad"), 0o644)
	if _, _, err := loadConfig(); err == nil {
		t.Fatalf("expected load error")
	}
	_ = app
}

func TestDoctor(t *testing.T) {
	app, cfg := newTestApp(t)
	os.RemoveAll(cfg.ReposRoot)
	os.RemoveAll(cfg.SkillsRoot)
	if code := app.runDoctor(context.Background(), nil); code != 0 {
		t.Fatalf("doctor failed")
	}
	// config load error
	path, _ := config.DefaultConfigPath()
	_ = os.WriteFile(path, []byte("{bad"), 0o644)
	if code := app.runDoctor(context.Background(), nil); code == 0 {
		t.Fatalf("expected doctor error")
	}
	// home error
	config.SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	defer config.ResetUserHomeDirForTest()
	if code := app.runDoctor(context.Background(), nil); code == 0 {
		t.Fatalf("expected doctor error")
	}
}

func TestDoctorNoGit(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("skip on windows")
	}
	app, _ := newTestApp(t)
	orig := os.Getenv("PATH")
	t.Setenv("PATH", "")
	if code := app.runDoctor(context.Background(), nil); code == 0 {
		t.Fatalf("expected doctor error")
	}
	t.Setenv("PATH", orig)
}

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
	// disallow code
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
	bare := initBareRepo(t, filepath.Join(cfg.ReposRoot, "remote.git"))
	if code := app.runRepo(context.Background(), []string{"clone", bare, "--name", "cloned"}); code != 0 {
		t.Fatalf("clone failed")
	}
	if code := app.runRepo(context.Background(), []string{"clone", bare, "--name", "cloned"}); code == 0 {
		t.Fatalf("expected clone error")
	}
	if code := app.runRepo(context.Background(), []string{"exec"}); code == 0 {
		t.Fatalf("expected exec error")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok"}); code != 0 {
		t.Fatalf("exec failed")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "exit 1", "--parallel", "2"}); code != 0 {
		t.Fatalf("exec error failed")
	}
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok", "--require-clean"}); code != 0 {
		t.Fatalf("exec require-clean failed")
	}
	_ = os.WriteFile(filepath.Join(cfg.ReposRoot, "alpha", "a.txt"), []byte("dirty"), 0o644)
	if code := app.runRepo(context.Background(), []string{"exec", "--cmd", "echo ok", "--require-clean"}); code != 0 {
		t.Fatalf("exec dirty failed")
	}
}

func TestSkillsCommands(t *testing.T) {
	app, cfg := newTestApp(t)
	// clone requires remote
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
	// existing git repo -> pull
	if code := app.runSkills(context.Background(), []string{"clone"}); code != 0 {
		t.Fatalf("clone pull failed")
	}
	_ = os.WriteFile(filepath.Join(cfg.SkillsRoot, "tool.txt"), []byte("x"), 0o644)
	// non-git existing without force
	_ = os.RemoveAll(cfg.SkillsRoot)
	_ = os.MkdirAll(cfg.SkillsRoot, 0o755)
	writeConfig(t, cfg)
	if code := app.runSkills(context.Background(), []string{"clone"}); code == 0 {
		t.Fatalf("expected clone error")
	}
	if code := app.runSkills(context.Background(), []string{"clone", "--force"}); code != 0 {
		t.Fatalf("clone force failed")
	}

	// sync/link
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

	// diff/verify/status
	if code := app.runSkills(context.Background(), []string{"diff", "--target", "skills"}); code != 0 {
		t.Fatalf("diff failed")
	}
	if code := app.runSkills(context.Background(), []string{"verify", "--target", "skills"}); code != 0 {
		t.Fatalf("verify failed")
	}
	if code := app.runSkills(context.Background(), []string{"status", "--target", "skills"}); code != 0 {
		t.Fatalf("status failed")
	}

	// pin
	if code := app.runSkills(context.Background(), []string{"pin"}); code == 0 {
		t.Fatalf("expected pin error")
	}
	if code := app.runSkills(context.Background(), []string{"pin", "--target", "missing", "--ref", "HEAD"}); code == 0 {
		t.Fatalf("expected target error")
	}
	if code := app.runSkills(context.Background(), []string{"pin", "--target", "skills", "--ref", "HEAD"}); code != 0 {
		t.Fatalf("pin failed")
	}

	// clean
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

func newTestApp(t *testing.T) (App, config.Config) {
	home := t.TempDir()
	config.SetUserHomeDirForTest(func() (string, error) { return home, nil })
	t.Cleanup(config.ResetUserHomeDirForTest)
	reposRoot := filepath.Join(home, "repos")
	skillsRoot := filepath.Join(home, "skills")
	_ = os.MkdirAll(reposRoot, 0o755)
	_ = os.MkdirAll(skillsRoot, 0o755)
	cfg := config.Config{
		ProjectsRoot:   home,
		ReposRoot:      reposRoot,
		SkillsRoot:     skillsRoot,
		SkillTargets:   []string{".codex/skills"},
		SyncTargets:    []config.SyncTarget{{Name: "skills", Src: skillsRoot, Dest: []string{".codex/skills"}, Include: []string{"**/*"}, Exclude: []string{".git/**"}}},
		AllowCommands:  []string{"*"},
		DenyCommands:   []string{"rm -rf*"},
		SyncMode:       "copy",
		ConflictPolicy: "overwrite",
	}
	writeConfig(t, cfg)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	app := App{Version: "test", Out: output.Writer{Out: &out, ErrW: &errBuf}}
	return app, cfg
}

func writeConfig(t *testing.T, cfg config.Config) {
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("config path: %v", err)
	}
	if err := config.Save(path, cfg); err != nil {
		t.Fatalf("config save: %v", err)
	}
}

func initGitRepo(t *testing.T, path string, commit bool) string {
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := runGit(path, "init"); err != nil {
		t.Fatalf("git init: %v", err)
	}
	if err := runGit(path, "config", "user.email", "test@example.com"); err != nil {
		t.Fatalf("git config: %v", err)
	}
	if err := runGit(path, "config", "user.name", "Tester"); err != nil {
		t.Fatalf("git config: %v", err)
	}
	if commit {
		_ = os.WriteFile(filepath.Join(path, "a.txt"), []byte("a"), 0o644)
		if err := runGit(path, "add", "."); err != nil {
			t.Fatalf("git add: %v", err)
		}
		if err := runGit(path, "commit", "-m", "init"); err != nil {
			t.Fatalf("git commit: %v", err)
		}
	}
	return path
}

func initBareRepo(t *testing.T, path string) string {
	if err := runGit(filepath.Dir(path), "init", "--bare", path); err != nil {
		t.Fatalf("git init bare: %v", err)
	}
	return path
}

func seedBareRepo(t *testing.T, bare string) {
	repo := initGitRepo(t, filepath.Join(t.TempDir(), "seed"), true)
	if err := runGit(repo, "remote", "add", "origin", bare); err != nil {
		t.Fatalf("git remote add: %v", err)
	}
	if err := runGit(repo, "push", "-u", "origin", "HEAD"); err != nil {
		t.Fatalf("git push: %v", err)
	}
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	return cmd.Run()
}

func createStubCommand(t *testing.T, name string) string {
	bin := filepath.Join(t.TempDir(), name)
	content := "#!/bin/sh\n" +
		"echo stub > /dev/null\n"
	if err := os.WriteFile(bin, []byte(content), 0o755); err != nil {
		t.Fatalf("stub write: %v", err)
	}
	return filepath.Dir(bin)
}

func prependPath(t *testing.T, bin string) {
	orig := os.Getenv("PATH")
	if err := os.Setenv("PATH", bin+string(os.PathListSeparator)+orig); err != nil {
		t.Fatalf("setenv: %v", err)
	}
}
