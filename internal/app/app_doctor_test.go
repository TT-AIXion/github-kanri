package app

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/TT-AIXion/github-kanri/internal/config"
	"github.com/TT-AIXion/github-kanri/internal/output"
)

func TestDoctor(t *testing.T) {
	app, cfg := newTestApp(t)
	os.RemoveAll(cfg.ReposRoot)
	os.RemoveAll(cfg.SkillsRoot)
	if code := app.runDoctor(context.Background(), nil); code != 0 {
		t.Fatalf("doctor failed")
	}
	path, _ := config.DefaultConfigPath()
	_ = os.WriteFile(path, []byte("{bad"), 0o644)
	if code := app.runDoctor(context.Background(), nil); code == 0 {
		t.Fatalf("expected doctor error")
	}
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

func TestDoctorParseError(t *testing.T) {
	app, _ := newTestApp(t)
	if code := app.runDoctor(context.Background(), []string{"--bad"}); code == 0 {
		t.Fatalf("expected doctor parse error")
	}
}

func TestDoctorExpandError(t *testing.T) {
	home := t.TempDir()
	config.SetUserHomeDirForTest(func() (string, error) { return home, nil })
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("path error: %v", err)
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
	app := App{Out: output.Writer{Out: os.Stdout, ErrW: os.Stderr}}
	if code := app.runDoctor(context.Background(), nil); code == 0 {
		t.Fatalf("expected expand error")
	}
}

func TestDoctorValidateError(t *testing.T) {
	home := t.TempDir()
	config.SetUserHomeDirForTest(func() (string, error) { return home, nil })
	defer config.ResetUserHomeDirForTest()
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("path error: %v", err)
	}
	cfg := config.Config{
		ProjectsRoot:   "x",
		ReposRoot:      "y",
		SkillsRoot:     "z",
		SkillTargets:   []string{".codex/skills"},
		SyncTargets:    []config.SyncTarget{{}},
		SyncMode:       "bad",
		ConflictPolicy: "fail",
		DenyCommands:   []string{"rm -rf*"},
	}
	if err := config.Save(path, cfg); err != nil {
		t.Fatalf("save: %v", err)
	}
	app := App{Out: output.Writer{Out: os.Stdout, ErrW: os.Stderr}}
	if code := app.runDoctor(context.Background(), nil); code == 0 {
		t.Fatalf("expected validate error")
	}
}
