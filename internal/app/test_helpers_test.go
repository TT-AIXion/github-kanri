package app

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/TT-AIXion/github-kanri/internal/config"
	"github.com/TT-AIXion/github-kanri/internal/output"
)

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
