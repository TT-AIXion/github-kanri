package app

import (
	"context"
	"os"
	"testing"

	"github.com/TT-AIXion/github-kanri/internal/config"
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
	if code := app.Run(context.Background(), []string{"quickstart"}); code == 0 {
		t.Fatalf("expected quickstart error")
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

func TestLoadConfigErrors(t *testing.T) {
	config.SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	if _, _, err := loadConfig(); err == nil {
		t.Fatalf("expected load error")
	}
	config.ResetUserHomeDirForTest()

	app, _ := newTestApp(t)
	path, _ := config.DefaultConfigPath()
	_ = os.WriteFile(path, []byte("{bad"), 0o644)
	if _, _, err := loadConfig(); err == nil {
		t.Fatalf("expected load error")
	}
	_ = app
}
