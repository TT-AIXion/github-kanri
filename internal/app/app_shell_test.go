package app

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestShellPrintAndInstall(t *testing.T) {
	app, _ := newTestApp(t)
	ctx := context.Background()
	for _, shell := range []string{"zsh", "bash", "fish", "powershell"} {
		if code := app.runShell(ctx, []string{shell}); code != 0 {
			t.Fatalf("shell print failed: %s", shell)
		}
	}
	if code := app.runShell(ctx, []string{"unknown"}); code == 0 {
		t.Fatalf("expected shell error")
	}

	profile := filepath.Join(t.TempDir(), "profile")
	if code := app.runShell(ctx, []string{"install", "--shell", "zsh", "--profile", profile}); code != 0 {
		t.Fatalf("shell install failed")
	}
	data, err := os.ReadFile(profile)
	if err != nil {
		t.Fatalf("read profile: %v", err)
	}
	if !strings.Contains(string(data), shellMarkerStart) {
		t.Fatalf("missing marker")
	}
	if code := app.runShell(ctx, []string{"install", "--shell", "zsh", "--profile", profile}); code != 0 {
		t.Fatalf("shell install idempotent failed")
	}

	dryProfile := filepath.Join(t.TempDir(), "dry", "profile")
	if code := app.runShell(ctx, []string{"install", "--shell", "zsh", "--profile", dryProfile, "--dry-run"}); code != 0 {
		t.Fatalf("shell install dry-run failed")
	}
	if _, err := os.Stat(dryProfile); err == nil {
		t.Fatalf("dry-run should not write")
	}
	if code := app.runShell(ctx, []string{"install", "--shell", "zsh", "--profile", profile, "--force"}); code != 0 {
		t.Fatalf("shell install force failed")
	}
}
