package gitutil

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/AIXion-Team/github-kanri/internal/executil"
	"github.com/AIXion-Team/github-kanri/internal/safety"
)

func TestGitUtils(t *testing.T) {
	root := t.TempDir()
	repoPath := filepath.Join(root, "repo")
	if err := runGit(root, "init", repoPath); err != nil {
		t.Fatalf("git init: %v", err)
	}
	if err := runGit(repoPath, "config", "user.email", "test@example.com"); err != nil {
		t.Fatalf("git config: %v", err)
	}
	if err := runGit(repoPath, "config", "user.name", "Tester"); err != nil {
		t.Fatalf("git config: %v", err)
	}
	file := filepath.Join(repoPath, "a.txt")
	_ = os.WriteFile(file, []byte("a"), 0o644)
	if err := runGit(repoPath, "add", "."); err != nil {
		t.Fatalf("git add: %v", err)
	}
	if err := runGit(repoPath, "commit", "-m", "init"); err != nil {
		t.Fatalf("git commit: %v", err)
	}

	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	clean, err := IsClean(context.Background(), runner, repoPath)
	if err != nil || !clean {
		t.Fatalf("expected clean")
	}
	_ = os.WriteFile(file, []byte("b"), 0o644)
	status, err := StatusPorcelain(context.Background(), runner, repoPath)
	if err != nil || status == "" {
		t.Fatalf("expected status")
	}
	branch, err := CurrentBranch(context.Background(), runner, repoPath)
	if err != nil || branch == "" {
		t.Fatalf("expected branch")
	}
	unix, err := LastCommitUnix(context.Background(), runner, repoPath)
	if err != nil || unix == 0 {
		t.Fatalf("expected unix")
	}
	log, err := LogOneline(context.Background(), runner, repoPath, 1)
	if err != nil || log == "" {
		t.Fatalf("expected log")
	}

	bare := filepath.Join(root, "bare.git")
	if err := runGit(root, "init", "--bare", bare); err != nil {
		t.Fatalf("git init bare: %v", err)
	}
	if err := runGit(repoPath, "remote", "add", "origin", bare); err != nil {
		t.Fatalf("git remote: %v", err)
	}
	if err := runGit(repoPath, "push", "-u", "origin", branch); err != nil {
		t.Fatalf("git push: %v", err)
	}
	originRefDir := filepath.Join(repoPath, ".git", "refs", "remotes", "origin")
	_ = os.MkdirAll(originRefDir, 0o755)
	_ = os.WriteFile(filepath.Join(originRefDir, "HEAD"), []byte("ref: refs/remotes/origin/"+branch+"\n"), 0o644)
	origin, err := OriginURL(context.Background(), runner, repoPath)
	if err != nil || !strings.Contains(origin, "bare.git") {
		t.Fatalf("expected origin")
	}
	def, err := DefaultBranch(context.Background(), runner, repoPath)
	if err != nil || def == "" {
		t.Fatalf("expected default branch")
	}

	clonePath := filepath.Join(root, "clone")
	if err := Clone(context.Background(), runner, bare, clonePath); err != nil {
		t.Fatalf("clone: %v", err)
	}
	if err := Fetch(context.Background(), runner, repoPath); err != nil {
		t.Fatalf("fetch: %v", err)
	}
	if err := Pull(context.Background(), runner, repoPath); err != nil {
		t.Fatalf("pull: %v", err)
	}
	if err := Checkout(context.Background(), runner, repoPath, branch); err != nil {
		t.Fatalf("checkout: %v", err)
	}
}

func TestDefaultBranchEmpty(t *testing.T) {
	root := t.TempDir()
	repoPath := filepath.Join(root, "repo")
	if err := runGit(root, "init", repoPath); err != nil {
		t.Fatalf("git init: %v", err)
	}
	if err := runGit(repoPath, "remote", "add", "origin", "https://example.com/foo.git"); err != nil {
		t.Fatalf("git remote: %v", err)
	}
	// create empty origin/HEAD ref to simulate bad output
	refPath := filepath.Join(repoPath, ".git", "refs", "remotes", "origin")
	_ = os.MkdirAll(refPath, 0o755)
	_ = os.WriteFile(filepath.Join(refPath, "HEAD"), []byte("ref: refs/remotes/origin/"), 0o644)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	_, _ = DefaultBranch(context.Background(), runner, repoPath)
}

func TestIsCleanError(t *testing.T) {
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"git log*"}}}
	if _, err := IsClean(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected IsClean error")
	}
}

func TestLastCommitUnixError(t *testing.T) {
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"git status*"}}}
	if _, err := LastCommitUnix(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected LastCommitUnix error")
	}
}

func TestDefaultBranchCommandError(t *testing.T) {
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"git status*"}}}
	if _, err := DefaultBranch(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected DefaultBranch error")
	}
}

func TestDefaultBranchUnexpected(t *testing.T) {
	orig := symbolicRef
	symbolicRef = func(context.Context, executil.Runner, string) (executil.Result, error) {
		return executil.Result{Stdout: "refs/remotes/origin/"}, nil
	}
	defer func() { symbolicRef = orig }()
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if _, err := DefaultBranch(context.Background(), runner, ""); err == nil {
		t.Fatalf("expected unexpected origin HEAD error")
	}
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	return cmd.Run()
}
