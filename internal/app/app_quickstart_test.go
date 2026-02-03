package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/TT-AIXion/github-kanri/internal/config"
	"github.com/TT-AIXion/github-kanri/internal/executil"
	"github.com/TT-AIXion/github-kanri/internal/fsutil"
	"github.com/TT-AIXion/github-kanri/internal/safety"
)

func TestQuickstartParseErrors(t *testing.T) {
	app, _ := newTestApp(t)
	cases := [][]string{
		{},
		{"--bad"},
		{"--public", "--private", "alpha"},
		{"alpha", "beta"},
		{"/"},
	}
	for _, args := range cases {
		if code := app.runQuickstart(context.Background(), args); code == 0 {
			t.Fatalf("expected error: %v", args)
		}
	}
}

func TestQuickstartPathErrors(t *testing.T) {
	app, cfg := newTestApp(t)
	cfg.SkillsRoot = filepath.Join(t.TempDir(), "missing")
	cfg.SyncTargets = []config.SyncTarget{{
		Name:    "skills",
		Src:     cfg.SkillsRoot,
		Dest:    []string{".codex/skills"},
		Include: []string{"**/*"},
		Exclude: []string{".git/**"},
	}}
	writeConfig(t, cfg)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected skillsRoot error")
	}

	app, cfg = newTestApp(t)
	dest := filepath.Join(cfg.ReposRoot, "alpha")
	if err := os.MkdirAll(dest, 0o755); err != nil {
		t.Fatalf("mkdir dest: %v", err)
	}
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected destination exists error")
	}

	app, cfg = newTestApp(t)
	cfg.DenyPaths = []string{filepath.Join(cfg.ReposRoot, "alpha")}
	writeConfig(t, cfg)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected deny path error")
	}
}

func TestQuickstartLoadConfigError(t *testing.T) {
	app, _ := newTestApp(t)
	path, err := config.DefaultConfigPath()
	if err != nil {
		t.Fatalf("config path: %v", err)
	}
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected load error")
	}
}

func TestQuickstartReposRootMkdirError(t *testing.T) {
	app, cfg := newTestApp(t)
	reposFile := filepath.Join(t.TempDir(), "repos")
	if err := os.WriteFile(reposFile, []byte("x"), 0o644); err != nil {
		t.Fatalf("write repos file: %v", err)
	}
	cfg.ReposRoot = reposFile
	writeConfig(t, cfg)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected reposRoot mkdir error")
	}
}

func TestQuickstartDestMkdirError(t *testing.T) {
	app, cfg := newTestApp(t)
	reposRoot := filepath.Join(t.TempDir(), "repos")
	if err := os.MkdirAll(reposRoot, 0o555); err != nil {
		t.Fatalf("mkdir reposRoot: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(reposRoot, 0o755) })
	cfg.ReposRoot = reposRoot
	writeConfig(t, cfg)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected dest mkdir error")
	}
}

func TestQuickstartAuthFailure(t *testing.T) {
	app, _ := newTestApp(t)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "echo no-auth >&2\nexit 1\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected gh auth error")
	}
}

func TestQuickstartSelectTargetsError(t *testing.T) {
	app, cfg := newTestApp(t)
	cfg.SyncTargets = []config.SyncTarget{{
		Name:    "other",
		Src:     cfg.SkillsRoot,
		Dest:    []string{".codex/skills"},
		Include: []string{"**/*"},
		Exclude: []string{".git/**"},
	}}
	writeConfig(t, cfg)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected target not found error")
	}
}

func TestQuickstartSyncTargetsError(t *testing.T) {
	app, cfg := newTestApp(t)
	cfg.DenyPaths = []string{cfg.SkillsRoot}
	writeConfig(t, cfg)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected sync target error")
	}
}

func TestQuickstartInitRepoMainError(t *testing.T) {
	app, _ := newTestApp(t)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ]; then echo init-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected init error")
	}
}

func TestQuickstartReadmeError(t *testing.T) {
	app, cfg := newTestApp(t)
	dest := filepath.Join(cfg.ReposRoot, "alpha")
	t.Cleanup(func() { _ = os.Chmod(dest, 0o755) })
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ]; then mkdir -p .git; chmod 0555 .; exit 0; fi\nif [ \"$1\" = \"checkout\" ]; then exit 0; fi\nexit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected readme error")
	}
}

func TestQuickstartCommitError(t *testing.T) {
	app, _ := newTestApp(t)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ]; then mkdir -p .git; exit 0; fi\nif [ \"$1\" = \"checkout\" ]; then exit 0; fi\nif [ \"$1\" = \"add\" ]; then exit 0; fi\nif [ \"$1\" = \"commit\" ]; then echo commit-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected commit error")
	}
}

func TestQuickstartRepoCreateError(t *testing.T) {
	app, _ := newTestApp(t)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "if [ \"$1\" = \"repo\" ] && [ \"$2\" = \"create\" ]; then echo gh-fail >&2; exit 1; fi\nexit 0\n")
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ]; then mkdir -p .git; exit 0; fi\nif [ \"$1\" = \"checkout\" ]; then exit 0; fi\nif [ \"$1\" = \"add\" ]; then exit 0; fi\nif [ \"$1\" = \"commit\" ]; then exit 0; fi\nexit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected gh repo create error")
	}
}

func TestQuickstartPullPushError(t *testing.T) {
	app, _ := newTestApp(t)
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "exit 0\n")
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ]; then mkdir -p .git; exit 0; fi\nif [ \"$1\" = \"checkout\" ]; then exit 0; fi\nif [ \"$1\" = \"add\" ]; then exit 0; fi\nif [ \"$1\" = \"commit\" ]; then exit 0; fi\nif [ \"$1\" = \"pull\" ]; then echo pull-fail >&2; exit 1; fi\nif [ \"$1\" = \"push\" ]; then exit 0; fi\nexit 0\n")
	setPath(t, binDir)
	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code == 0 {
		t.Fatalf("expected pull/push error")
	}
}

func TestQuickstartHappyPath(t *testing.T) {
	app, cfg := newTestApp(t)
	if err := os.WriteFile(filepath.Join(cfg.SkillsRoot, "sample.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatalf("skills write: %v", err)
	}
	logPath := filepath.Join(t.TempDir(), "gh.log")
	remoteRoot := filepath.Join(t.TempDir(), "remotes")
	binDir := createGHStub(t, remoteRoot, logPath)
	setPath(t, binDir)
	t.Setenv("GIT_AUTHOR_NAME", "Test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")

	if code := app.runQuickstart(context.Background(), []string{"alpha"}); code != 0 {
		t.Fatalf("expected quickstart")
	}
	dest := filepath.Join(cfg.ReposRoot, "alpha")
	if !fsutil.IsGitRepo(dest) {
		t.Fatalf("expected git repo: %s", dest)
	}
	if _, err := os.Stat(filepath.Join(dest, ".codex", "skills", "sample.txt")); err != nil {
		t.Fatalf("expected skills copy: %v", err)
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(data), "repo create") {
		t.Fatalf("expected gh repo create call")
	}
}

func TestQuickstartPublicFlag(t *testing.T) {
	app, cfg := newTestApp(t)
	if err := os.WriteFile(filepath.Join(cfg.SkillsRoot, "sample.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatalf("skills write: %v", err)
	}
	logPath := filepath.Join(t.TempDir(), "gh.log")
	remoteRoot := filepath.Join(t.TempDir(), "remotes")
	binDir := createGHStub(t, remoteRoot, logPath)
	setPath(t, binDir)
	t.Setenv("GIT_AUTHOR_NAME", "Test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")

	if code := app.runQuickstart(context.Background(), []string{"--public", "public-repo"}); code != 0 {
		t.Fatalf("expected quickstart public")
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(data), "--public") {
		t.Fatalf("expected public flag")
	}
}

func TestQuickstartPrivateFlag(t *testing.T) {
	app, cfg := newTestApp(t)
	if err := os.WriteFile(filepath.Join(cfg.SkillsRoot, "sample.txt"), []byte("ok"), 0o644); err != nil {
		t.Fatalf("skills write: %v", err)
	}
	logPath := filepath.Join(t.TempDir(), "gh.log")
	remoteRoot := filepath.Join(t.TempDir(), "remotes")
	binDir := createGHStub(t, remoteRoot, logPath)
	setPath(t, binDir)
	t.Setenv("GIT_AUTHOR_NAME", "Test")
	t.Setenv("GIT_AUTHOR_EMAIL", "test@example.com")
	t.Setenv("GIT_COMMITTER_NAME", "Test")
	t.Setenv("GIT_COMMITTER_EMAIL", "test@example.com")

	if code := app.runQuickstart(context.Background(), []string{"--private", "private-repo"}); code != 0 {
		t.Fatalf("expected quickstart private")
	}
	data, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("read log: %v", err)
	}
	if !strings.Contains(string(data), "--private") {
		t.Fatalf("expected private flag")
	}
}

func TestParseQuickstartName(t *testing.T) {
	gotName, gotRepo, err := parseQuickstartName("owner/repo.git")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "owner/repo" || gotRepo != "repo" {
		t.Fatalf("unexpected parse: %s %s", gotName, gotRepo)
	}
	gotName, gotRepo, err = parseQuickstartName("repo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if gotName != "repo" || gotRepo != "repo" {
		t.Fatalf("unexpected parse: %s %s", gotName, gotRepo)
	}
	for _, arg := range []string{"", " ", ".", "/"} {
		if _, _, err := parseQuickstartName(arg); err == nil {
			t.Fatalf("expected error: %q", arg)
		}
	}
}

func TestEnsureSkillsRoot(t *testing.T) {
	if err := ensureSkillsRoot(""); err == nil {
		t.Fatalf("expected error")
	}
	missing := filepath.Join(t.TempDir(), "missing")
	if err := ensureSkillsRoot(missing); err == nil {
		t.Fatalf("expected missing error")
	}
	parent := filepath.Join(t.TempDir(), "parent")
	child := filepath.Join(parent, "child")
	if err := os.MkdirAll(child, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.Chmod(parent, 0o000); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(parent, 0o755) })
	if err := ensureSkillsRoot(child); err == nil {
		t.Fatalf("expected permission error")
	}
	if err := ensureSkillsRoot(t.TempDir()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommandError(t *testing.T) {
	if got := commandError(nil, executil.Result{}); got != "" {
		t.Fatalf("expected empty")
	}
	err := errors.New("boom")
	if got := commandError(err, executil.Result{Stderr: "fail\n"}); got != "boom: fail" {
		t.Fatalf("unexpected stderr: %s", got)
	}
	if got := commandError(err, executil.Result{Stdout: "out\n"}); got != "boom: out" {
		t.Fatalf("unexpected stdout: %s", got)
	}
	if got := commandError(err, executil.Result{}); got != "boom" {
		t.Fatalf("unexpected fallback: %s", got)
	}
}

func TestGhAuthStatusError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "echo no-auth >&2\nexit 1\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := ghAuthStatus(context.Background(), runner); err == nil {
		t.Fatalf("expected auth error")
	}
}

func TestInitRepoMainFallback(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ] && [ \"$2\" = \"-b\" ]; then exit 1; fi\nif [ \"$1\" = \"init\" ]; then exit 0; fi\nif [ \"$1\" = \"checkout\" ]; then exit 0; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := initRepoMain(context.Background(), runner, t.TempDir()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestInitRepoMainError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ]; then echo init-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := initRepoMain(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected error")
	}
}

func TestInitRepoMainCheckoutError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "git", "if [ \"$1\" = \"init\" ] && [ \"$2\" = \"-b\" ]; then exit 1; fi\nif [ \"$1\" = \"init\" ]; then exit 0; fi\nif [ \"$1\" = \"checkout\" ]; then echo checkout-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := initRepoMain(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected checkout error")
	}
}

func TestWriteQuickstartReadme(t *testing.T) {
	dest := t.TempDir()
	if err := writeQuickstartReadme(dest, "alpha"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := writeQuickstartReadme(dest, "alpha"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGitCommitInitError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "git", "if [ \"$1\" = \"add\" ]; then exit 0; fi\nif [ \"$1\" = \"commit\" ]; then echo commit-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := gitCommitInit(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected commit error")
	}
}

func TestGitCommitInitAddError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "git", "if [ \"$1\" = \"add\" ]; then echo add-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := gitCommitInit(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected add error")
	}
}

func TestGhRepoCreateError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "gh", "if [ \"$1\" = \"repo\" ] && [ \"$2\" = \"create\" ]; then echo gh-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := ghRepoCreate(context.Background(), runner, "owner/repo", t.TempDir(), "--private"); err == nil {
		t.Fatalf("expected repo create error")
	}
}

func TestGitPullPushError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "git", "if [ \"$1\" = \"pull\" ]; then echo pull-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := gitPullPush(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected pull error")
	}
}

func TestGitPullPushPushError(t *testing.T) {
	binDir := t.TempDir()
	writeStub(t, binDir, "git", "if [ \"$1\" = \"pull\" ]; then exit 0; fi\nif [ \"$1\" = \"push\" ]; then echo push-fail >&2; exit 1; fi\nexit 0\n")
	setPath(t, binDir)
	runner := executil.Runner{Guard: safety.Guard{AllowCommands: []string{"*"}}}
	if err := gitPullPush(context.Background(), runner, t.TempDir()); err == nil {
		t.Fatalf("expected push error")
	}
}

func createGHStub(t *testing.T, remoteRoot string, logPath string) string {
	t.Helper()
	if err := os.MkdirAll(remoteRoot, 0o755); err != nil {
		t.Fatalf("remote root: %v", err)
	}
	bin := filepath.Join(t.TempDir(), "gh")
	script := fmt.Sprintf(`#!/bin/sh
set -e
log=%q
remote_root=%q
mkdir -p "$remote_root"
printf "gh %%s\n" "$*" >> "$log"
if [ "$1" = "auth" ] && [ "$2" = "status" ]; then
  exit 0
fi
if [ "$1" = "repo" ] && [ "$2" = "create" ]; then
  src=""
  remote="origin"
  push=0
  while [ "$#" -gt 0 ]; do
    case "$1" in
      --source)
        src="$2"
        shift 2
        ;;
      --source=*)
        src="${1#--source=}"
        shift
        ;;
      --remote)
        remote="$2"
        shift 2
        ;;
      --remote=*)
        remote="${1#--remote=}"
        shift
        ;;
      --push)
        push=1
        shift
        ;;
      *)
        shift
        ;;
    esac
  done
  if [ -z "$src" ]; then
    echo "missing --source" >&2
    exit 1
  fi
  bare=$(mktemp -d "$remote_root/gkn-remote-XXXXXX")
  git init --bare "$bare" >/dev/null
  git -C "$src" remote add "$remote" "$bare"
  if [ "$push" -eq 1 ]; then
    git -C "$src" push -u "$remote" main >/dev/null
  fi
  exit 0
fi
echo "unexpected gh args" >&2
exit 1
`, logPath, remoteRoot)
	if err := os.WriteFile(bin, []byte(script), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
	return filepath.Dir(bin)
}

func writeStub(t *testing.T, dir string, name string, body string) {
	t.Helper()
	path := filepath.Join(dir, name)
	content := "#!/bin/sh\nset -e\n" + body
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatalf("write stub: %v", err)
	}
}

func setPath(t *testing.T, dir string) {
	t.Helper()
	orig := os.Getenv("PATH")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+orig)
}
