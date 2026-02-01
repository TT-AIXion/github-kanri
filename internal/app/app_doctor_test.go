package app

import (
	"context"
	"os"
	"runtime"
	"testing"

	"github.com/AIXion-Team/github-kanri/internal/config"
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
