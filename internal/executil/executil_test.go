package executil

import (
	"context"
	"testing"
	"time"

	"github.com/AIXion-Team/github-kanri/internal/safety"
)

func TestRunDry(t *testing.T) {
	r := Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}, DryRun: true}
	res, err := r.Run(context.Background(), "", "echo", "hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Fatalf("expected exit 0")
	}
}

func TestRun(t *testing.T) {
	r := Runner{Guard: safety.Guard{AllowCommands: []string{"echo*", "false*"}}}
	res, err := r.Run(context.Background(), "", "echo", "hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stdout == "" {
		t.Fatalf("expected stdout")
	}
	_, err = r.Run(context.Background(), "", "false")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestRunShell(t *testing.T) {
	r := Runner{Guard: safety.Guard{AllowCommands: []string{"echo*", "exit*"}}}
	res, err := r.RunShell(context.Background(), "", "echo hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Stdout == "" {
		t.Fatalf("expected stdout")
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()
	_, _ = r.RunShell(ctx, "", "exit 1")
}

func TestRunDenied(t *testing.T) {
	r := Runner{Guard: safety.Guard{AllowCommands: []string{"git*"}}}
	if _, err := r.Run(context.Background(), "", "echo", "hi"); err == nil {
		t.Fatalf("expected deny")
	}
}

func TestRunShellDry(t *testing.T) {
	r := Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}, DryRun: true}
	if _, err := r.RunShell(context.Background(), "", "echo hi"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunShellDenied(t *testing.T) {
	r := Runner{Guard: safety.Guard{AllowCommands: []string{"git*"}}}
	if _, err := r.RunShell(context.Background(), "", "echo hi"); err == nil {
		t.Fatalf("expected deny")
	}
}

func TestExitCodeDefault(t *testing.T) {
	if code := exitCode(context.DeadlineExceeded); code != 1 {
		t.Fatalf("expected 1")
	}
}
