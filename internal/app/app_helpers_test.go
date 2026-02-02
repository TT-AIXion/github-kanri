package app

import (
	"os"
	"testing"

	"github.com/TT-AIXion/github-kanri/internal/config"
)

func TestLoadConfigPathError(t *testing.T) {
	config.SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	defer config.ResetUserHomeDirForTest()
	if _, _, err := loadConfig(); err == nil {
		t.Fatalf("expected path error")
	}
}

func TestLoadConfigExpandError(t *testing.T) {
	home := t.TempDir()
	config.SetUserHomeDirForTest(func() (string, error) { return home, nil })
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("path: %v", err)
	}
	cfg := config.Config{
		ProjectsRoot: "~/Projects",
		ReposRoot:    "~/Projects/repos",
		SkillsRoot:   "~/Projects/skills",
		SkillTargets: []string{".codex/skills"},
		SyncTargets: []config.SyncTarget{{
			Name: "skills",
			Src:  "~/Projects/skills",
			Dest: []string{".codex/skills"},
		}},
		SyncMode:       "copy",
		ConflictPolicy: "fail",
		DenyCommands:   []string{"rm -rf*"},
	}
	if err := config.Save(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}
	calls := 0
	config.SetUserHomeDirForTest(func() (string, error) {
		calls++
		if calls == 1 {
			return home, nil
		}
		return "", os.ErrPermission
	})
	defer config.ResetUserHomeDirForTest()
	if _, _, err := loadConfig(); err == nil {
		t.Fatalf("expected expand error")
	}
}
