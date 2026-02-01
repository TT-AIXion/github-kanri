package fsutil

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFilterNamesNoFilters(t *testing.T) {
	names := []string{"alpha", "beta"}
	filtered := FilterNames(names, nil, nil)
	if len(filtered) != 2 {
		t.Fatalf("unexpected: %v", filtered)
	}
}

func TestFilterNamesOnlyMismatch(t *testing.T) {
	names := []string{"alpha", "beta"}
	filtered := FilterNames(names, []string{"alpha"}, nil)
	if len(filtered) != 1 || filtered[0] != "alpha" {
		t.Fatalf("unexpected: %v", filtered)
	}
}

func TestListFilesDefaultsAndExclude(t *testing.T) {
	root := t.TempDir()
	_ = os.MkdirAll(filepath.Join(root, "skipdir"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "skipdir", "a.txt"), []byte("a"), 0o644)
	_ = os.WriteFile(filepath.Join(root, "b.txt"), []byte("b"), 0o644)
	files, err := ListFiles(root, nil, []string{"skipdir/**"})
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(files) != 1 || files[0] != "b.txt" {
		t.Fatalf("unexpected files: %v", files)
	}
}

func TestListFilesIncludeSkip(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "a.txt"), []byte("a"), 0o644)
	files, err := ListFiles(root, []string{"**/*.go"}, nil)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(files) != 0 {
		t.Fatalf("expected empty")
	}
}

func TestListGitReposWithFile(t *testing.T) {
	root := t.TempDir()
	repo := filepath.Join(root, "repo")
	_ = os.MkdirAll(filepath.Join(repo, ".git"), 0o755)
	_ = os.WriteFile(filepath.Join(root, "note.txt"), []byte("x"), 0o644)
	if _, err := ListGitRepos(root); err != nil {
		t.Fatalf("list: %v", err)
	}
}

func TestDiffDirDestMissing(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	_ = os.MkdirAll(src, 0o755)
	_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
	added, removed, changed, err := DiffDir(src, filepath.Join(root, "missing"), []string{"**/*"}, nil)
	if err != nil || len(added) == 0 || len(removed) != 0 || len(changed) != 0 {
		t.Fatalf("unexpected diff: %v %v %v %v", added, removed, changed, err)
	}
}

func TestMirrorCleanupKeeps(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dest")
	_ = os.MkdirAll(filepath.Join(dest, "dir"), 0o755)
	_ = os.WriteFile(filepath.Join(dest, "keep.txt"), []byte("a"), 0o644)
	if err := mirrorCleanup(root, dest, []string{"keep.txt"}, SyncOptions{ConflictPolicy: ConflictOverwrite}); err != nil {
		t.Fatalf("mirror: %v", err)
	}
}

func TestMirrorCleanupRemoves(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dest")
	_ = os.MkdirAll(dest, 0o755)
	_ = os.WriteFile(filepath.Join(dest, "remove.txt"), []byte("a"), 0o644)
	if err := mirrorCleanup(root, dest, nil, SyncOptions{ConflictPolicy: ConflictOverwrite}); err != nil {
		t.Fatalf("mirror remove: %v", err)
	}
}

func TestCleanDirKeeps(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dest")
	_ = os.MkdirAll(filepath.Join(dest, "dir"), 0o755)
	_ = os.WriteFile(filepath.Join(dest, "keep.txt"), []byte("a"), 0o644)
	if err := CleanDir(dest, []string{"keep.txt"}, false); err != nil {
		t.Fatalf("clean: %v", err)
	}
}

func TestCleanDirRemoves(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dest")
	_ = os.MkdirAll(dest, 0o755)
	_ = os.WriteFile(filepath.Join(dest, "remove.txt"), []byte("a"), 0o644)
	if err := CleanDir(dest, nil, false); err != nil {
		t.Fatalf("clean remove: %v", err)
	}
}
