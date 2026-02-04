package fsutil

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestIsGitRepoAndListGitRepos(t *testing.T) {
	root := t.TempDir()
	repo1 := filepath.Join(root, "repo1")
	repo2 := filepath.Join(root, "repo2")
	if err := os.MkdirAll(filepath.Join(repo1, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(repo2, ".git"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if !IsGitRepo(repo1) {
		t.Fatalf("expected git repo")
	}
	repos, err := ListGitRepos(root)
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(repos) != 2 {
		t.Fatalf("unexpected repos: %v", repos)
	}
}

func TestListGitReposError(t *testing.T) {
	if _, err := ListGitRepos(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatalf("expected error")
	}
}

func TestFilterNames(t *testing.T) {
	names := []string{"alpha", "beta", "gamma"}
	filtered := FilterNames(names, []string{"*a"}, []string{"gamma"})
	if len(filtered) != 2 {
		t.Fatalf("unexpected: %v", filtered)
	}
}

func TestListFilesAndEnsureDir(t *testing.T) {
	root := t.TempDir()
	if err := EnsureDir(filepath.Join(root, "a/b"), false); err != nil {
		t.Fatalf("ensure: %v", err)
	}
	if err := EnsureDir(filepath.Join(root, "dry"), true); err != nil {
		t.Fatalf("dry ensure: %v", err)
	}
	_ = os.WriteFile(filepath.Join(root, "a", "b", "file.txt"), []byte("hi"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "skip.log"), []byte("skip"), 0o644)
	files, err := ListFiles(root, []string{"**/*.txt"}, []string{"skip*"})
	if err != nil {
		t.Fatalf("list files: %v", err)
	}
	if len(files) != 1 || !strings.Contains(files[0], "file.txt") {
		t.Fatalf("unexpected files: %v", files)
	}
}

func TestCopyLinkRemove(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	dst := filepath.Join(root, "dst.txt")
	if err := os.WriteFile(src, []byte("data"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := CopyFile(src, dst, false); err != nil {
		t.Fatalf("copy: %v", err)
	}
	if err := CopyFile(src, filepath.Join(root, "dry.txt"), true); err != nil {
		t.Fatalf("dry copy: %v", err)
	}
	link := filepath.Join(root, "link.txt")
	if err := LinkFile(src, link, false); err != nil {
		t.Fatalf("link: %v", err)
	}
	if err := LinkFile(src, filepath.Join(root, "drylink.txt"), true); err != nil {
		t.Fatalf("dry link: %v", err)
	}
	if err := RemovePath(link, false); err != nil {
		t.Fatalf("remove: %v", err)
	}
	if err := RemovePath(filepath.Join(root, "dryremove"), true); err != nil {
		t.Fatalf("dry remove: %v", err)
	}
}

func TestSyncDirCopyConflict(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")
	_ = os.MkdirAll(src, 0o755)
	_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
	_ = os.MkdirAll(dst, 0o755)
	_ = os.WriteFile(filepath.Join(dst, "a.txt"), []byte("old"), 0o644)
	if err := SyncDir(src, dst, SyncOptions{Mode: ModeCopy, ConflictPolicy: ConflictFail}); err == nil {
		t.Fatalf("expected conflict")
	}
	if err := SyncDir(src, dst, SyncOptions{Mode: ModeCopy, ConflictPolicy: ConflictOverwrite}); err != nil {
		t.Fatalf("overwrite: %v", err)
	}
	if err := SyncDir(src, dst, SyncOptions{Mode: ModeCopy, ConflictPolicy: ConflictOverwrite, DryRun: true}); err != nil {
		t.Fatalf("dry overwrite: %v", err)
	}
}

func TestSyncDirLinkMirror(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")
	_ = os.MkdirAll(src, 0o755)
	_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
	if err := SyncDir(src, dst, SyncOptions{Mode: ModeLink, ConflictPolicy: ConflictOverwrite}); err != nil {
		t.Fatalf("link: %v", err)
	}
	if err := SyncDir(src, dst, SyncOptions{Mode: ModeMirror, ConflictPolicy: ConflictFail}); err == nil {
		t.Fatalf("expected mirror error")
	}
	_ = os.WriteFile(filepath.Join(dst, "extra.txt"), []byte("x"), 0o644)
	if err := SyncDir(src, dst, SyncOptions{Mode: ModeMirror, ConflictPolicy: ConflictOverwrite}); err != nil {
		t.Fatalf("mirror: %v", err)
	}
}

func TestDiffAndClean(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dst := filepath.Join(root, "dst")
	_ = os.MkdirAll(src, 0o755)
	_ = os.MkdirAll(dst, 0o755)
	_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
	_ = os.WriteFile(filepath.Join(dst, "b.txt"), []byte("b"), 0o644)
	var changed []string
	added, removed, _, err := DiffDir(src, dst, []string{"**/*"}, nil)
	if err != nil {
		t.Fatalf("diff: %v", err)
	}
	if len(added) == 0 || len(removed) == 0 {
		t.Fatalf("expected diff")
	}
	_ = os.WriteFile(filepath.Join(dst, "a.txt"), []byte("c"), 0o644)
	_, _, changed, err = DiffDir(src, dst, []string{"**/*"}, nil)
	if err != nil || len(changed) == 0 {
		t.Fatalf("expected changed")
	}
	files, err := ListFiles(src, []string{"**/*"}, nil)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if err := CleanDir(dst, files, false); err != nil {
		t.Fatalf("clean: %v", err)
	}
}

func TestResolvePath(t *testing.T) {
	root := t.TempDir()
	if got := ResolvePath(root, ""); got != root {
		t.Fatalf("expected root")
	}
	abs := filepath.Join(root, "abs")
	if got := ResolvePath(root, abs); got != abs {
		t.Fatalf("expected abs")
	}
	if got := ResolvePath(root, "rel"); got != filepath.Join(root, "rel") {
		t.Fatalf("expected rel")
	}
}
