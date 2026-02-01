package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/AIXion-Team/github-kanri/internal/config"
)

func TestConsumeJSONFlag(t *testing.T) {
	jsonMode, args := consumeJSONFlag([]string{"--json", "repo", "list"})
	if !jsonMode {
		t.Fatalf("expected json mode")
	}
	if len(args) != 2 || args[0] != "repo" {
		t.Fatalf("unexpected args")
	}
}

func TestRunMain(t *testing.T) {
	tmp := t.TempDir()
	config.SetUserHomeDirForTest(func() (string, error) { return tmp, nil })
	defer config.ResetUserHomeDirForTest()
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("path error: %v", err)
	}
	cfg := config.Config{ProjectsRoot: "~/Projects", ReposRoot: tmp, SkillsRoot: filepath.Join(tmp, "skills"), AllowCommands: []string{"*"}, DenyCommands: []string{"rm -rf*"}, SyncMode: "copy", ConflictPolicy: "fail"}
	if err := config.Save(path, cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}
	if code := runMain([]string{"version"}); code != 0 {
		t.Fatalf("expected 0")
	}
	if code := runMain([]string{"help"}); code != 0 {
		t.Fatalf("expected 0")
	}
	if code := runMain([]string{"unknown"}); code == 0 {
		t.Fatalf("expected non-zero")
	}
	_ = os.RemoveAll(tmp)
}

func TestMainFunc(t *testing.T) {
	orig := exitFunc
	defer func() { exitFunc = orig }()
	exitFunc = func(code int) {}
	main()
}
