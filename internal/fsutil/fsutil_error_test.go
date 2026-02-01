package fsutil

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

type fsutilHooks struct {
	mkdirAll  func(string, os.FileMode) error
	open      func(string) (*os.File, error)
	create    func(string) (*os.File, error)
	chmod     func(string, os.FileMode) error
	symlink   func(string, string) error
	removeAll func(string) error
	lstat     func(string) (os.FileInfo, error)
	walkDir   func(string, fs.WalkDirFunc) error
	relPath   func(string, string) (string, error)
	copy      func(io.Writer, io.Reader) (int64, error)
	stat      func(*os.File) (os.FileInfo, error)
}

func snapshotHooks() fsutilHooks {
	return fsutilHooks{
		mkdirAll:  osMkdirAll,
		open:      osOpen,
		create:    osCreate,
		chmod:     osChmod,
		symlink:   osSymlink,
		removeAll: osRemoveAll,
		lstat:     osLstat,
		walkDir:   walkDir,
		relPath:   relPath,
		copy:      ioCopy,
		stat:      fileStat,
	}
}

func (h fsutilHooks) restore() {
	osMkdirAll = h.mkdirAll
	osOpen = h.open
	osCreate = h.create
	osChmod = h.chmod
	osSymlink = h.symlink
	osRemoveAll = h.removeAll
	osLstat = h.lstat
	walkDir = h.walkDir
	relPath = h.relPath
	ioCopy = h.copy
	fileStat = h.stat
}

func TestCopyFileErrors(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	if err := os.WriteFile(src, []byte("data"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	t.Run("mkdir", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osMkdirAll = func(string, os.FileMode) error { return errors.New("mkdir") }
		if err := CopyFile(src, filepath.Join(root, "a", "dst.txt"), false); err == nil {
			t.Fatalf("expected mkdir error")
		}
	})

	t.Run("open", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osOpen = func(string) (*os.File, error) { return nil, errors.New("open") }
		if err := CopyFile(src, filepath.Join(root, "dst.txt"), false); err == nil {
			t.Fatalf("expected open error")
		}
	})

	t.Run("stat", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		fileStat = func(*os.File) (os.FileInfo, error) { return nil, errors.New("stat") }
		if err := CopyFile(src, filepath.Join(root, "dst.txt"), false); err == nil {
			t.Fatalf("expected stat error")
		}
	})

	t.Run("create", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osCreate = func(string) (*os.File, error) { return nil, errors.New("create") }
		if err := CopyFile(src, filepath.Join(root, "dst.txt"), false); err == nil {
			t.Fatalf("expected create error")
		}
	})

	t.Run("copy", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		ioCopy = func(io.Writer, io.Reader) (int64, error) { return 0, errors.New("copy") }
		if err := CopyFile(src, filepath.Join(root, "dst.txt"), false); err == nil {
			t.Fatalf("expected copy error")
		}
	})

	t.Run("chmod", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osChmod = func(string, os.FileMode) error { return errors.New("chmod") }
		if err := CopyFile(src, filepath.Join(root, "dst.txt"), false); err == nil {
			t.Fatalf("expected chmod error")
		}
	})
}

func TestLinkFileErrors(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src.txt")
	_ = os.WriteFile(src, []byte("data"), 0o644)

	t.Run("mkdir", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osMkdirAll = func(string, os.FileMode) error { return errors.New("mkdir") }
		if err := LinkFile(src, filepath.Join(root, "a", "link.txt"), false); err == nil {
			t.Fatalf("expected mkdir error")
		}
	})

	t.Run("symlink", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osSymlink = func(string, string) error { return errors.New("symlink") }
		if err := LinkFile(src, filepath.Join(root, "link.txt"), false); err == nil {
			t.Fatalf("expected symlink error")
		}
	})
}

func TestListFilesErrors(t *testing.T) {
	root := t.TempDir()
	_ = os.WriteFile(filepath.Join(root, "a.txt"), []byte("a"), 0o644)

	t.Run("walk", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		walkDir = func(string, fs.WalkDirFunc) error { return errors.New("walk") }
		if _, err := ListFiles(root, []string{"**/*"}, nil); err == nil {
			t.Fatalf("expected walk error")
		}
	})

	t.Run("rel", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		relPath = func(string, string) (string, error) { return "", errors.New("rel") }
		if _, err := ListFiles(root, []string{"**/*"}, nil); err == nil {
			t.Fatalf("expected rel error")
		}
	})
}

func TestMirrorCleanupErrors(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dest")
	_ = os.MkdirAll(dest, 0o755)
	_ = os.WriteFile(filepath.Join(dest, "a.txt"), []byte("a"), 0o644)

	t.Run("callback", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		walkDir = func(root string, fn fs.WalkDirFunc) error {
			return fn(root, fakeDirEntry{name: "x", dir: true}, errors.New("walk"))
		}
		if err := mirrorCleanup(root, dest, nil, SyncOptions{ConflictPolicy: ConflictOverwrite}); err == nil {
			t.Fatalf("expected callback error")
		}
	})

	t.Run("conflict-policy", func(t *testing.T) {
		if err := mirrorCleanup(root, dest, nil, SyncOptions{ConflictPolicy: ConflictFail}); err == nil {
			t.Fatalf("expected conflict policy error")
		}
	})

	t.Run("walk", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		walkDir = func(string, fs.WalkDirFunc) error { return errors.New("walk") }
		if err := mirrorCleanup(root, dest, nil, SyncOptions{ConflictPolicy: ConflictOverwrite}); err == nil {
			t.Fatalf("expected walk error")
		}
	})

	t.Run("rel", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		relPath = func(string, string) (string, error) { return "", errors.New("rel") }
		if err := mirrorCleanup(root, dest, nil, SyncOptions{ConflictPolicy: ConflictOverwrite}); err == nil {
			t.Fatalf("expected rel error")
		}
	})

	t.Run("remove", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osRemoveAll = func(string) error { return errors.New("remove") }
		if err := mirrorCleanup(root, dest, nil, SyncOptions{ConflictPolicy: ConflictOverwrite}); err == nil {
			t.Fatalf("expected remove error")
		}
	})
}

func TestCleanDirErrors(t *testing.T) {
	root := t.TempDir()
	dest := filepath.Join(root, "dest")
	_ = os.MkdirAll(dest, 0o755)
	_ = os.WriteFile(filepath.Join(dest, "a.txt"), []byte("a"), 0o644)

	t.Run("walk", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		walkDir = func(string, fs.WalkDirFunc) error { return errors.New("walk") }
		if err := CleanDir(dest, nil, false); err == nil {
			t.Fatalf("expected walk error")
		}
	})

	t.Run("callback", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		walkDir = func(root string, fn fs.WalkDirFunc) error {
			return fn(root, fakeDirEntry{name: "x", dir: true}, errors.New("walk"))
		}
		if err := CleanDir(dest, nil, false); err == nil {
			t.Fatalf("expected callback error")
		}
	})

	t.Run("rel", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		relPath = func(string, string) (string, error) { return "", errors.New("rel") }
		if err := CleanDir(dest, nil, false); err == nil {
			t.Fatalf("expected rel error")
		}
	})

	t.Run("remove", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		osRemoveAll = func(string) error { return errors.New("remove") }
		if err := CleanDir(dest, nil, false); err == nil {
			t.Fatalf("expected remove error")
		}
	})
}

func TestDiffDirErrors(t *testing.T) {
	root := t.TempDir()
	src := filepath.Join(root, "src")
	dest := filepath.Join(root, "dest")
	_ = os.MkdirAll(src, 0o755)
	_ = os.MkdirAll(dest, 0o755)
	_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
	_ = os.WriteFile(filepath.Join(dest, "a.txt"), []byte("a"), 0o644)

	t.Run("src", func(t *testing.T) {
		if _, _, _, err := DiffDir(filepath.Join(root, "missing"), dest, []string{"**/*"}, nil); err == nil {
			t.Fatalf("expected src error")
		}
	})

	t.Run("list", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		walkDir = func(root string, fn fs.WalkDirFunc) error {
			if root == dest {
				return errors.New("walk")
			}
			return h.walkDir(root, fn)
		}
		if _, _, _, err := DiffDir(src, dest, []string{"**/*"}, nil); err == nil {
			t.Fatalf("expected diff list error")
		}
	})

	t.Run("hash", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		ioCopy = func(io.Writer, io.Reader) (int64, error) { return 0, errors.New("copy") }
		if _, _, _, err := DiffDir(src, dest, []string{"**/*"}, nil); err == nil {
			t.Fatalf("expected diff hash error")
		}
	})

	t.Run("dest-hash", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		destFile := filepath.Join(dest, "a.txt")
		osOpen = func(path string) (*os.File, error) {
			if path == destFile {
				return nil, errors.New("open")
			}
			return h.open(path)
		}
		if _, _, _, err := DiffDir(src, dest, []string{"**/*"}, nil); err == nil {
			t.Fatalf("expected dest hash error")
		}
	})
}

func TestSyncDirErrors(t *testing.T) {
	t.Run("list", func(t *testing.T) {
		if err := SyncDir(filepath.Join(t.TempDir(), "missing"), t.TempDir(), SyncOptions{Mode: ModeCopy}); err == nil {
			t.Fatalf("expected list error")
		}
	})

	t.Run("link-conflict", func(t *testing.T) {
		root := t.TempDir()
		src := filepath.Join(root, "src")
		dst := filepath.Join(root, "dst")
		_ = os.MkdirAll(src, 0o755)
		_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
		_ = os.MkdirAll(dst, 0o755)
		_ = os.WriteFile(filepath.Join(dst, "a.txt"), []byte("old"), 0o644)
		if err := SyncDir(src, dst, SyncOptions{Mode: ModeLink, ConflictPolicy: ConflictFail}); err == nil {
			t.Fatalf("expected link conflict error")
		}
	})

	t.Run("link", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		root := t.TempDir()
		src := filepath.Join(root, "src")
		dst := filepath.Join(root, "dst")
		_ = os.MkdirAll(src, 0o755)
		_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
		osSymlink = func(string, string) error { return errors.New("symlink") }
		if err := SyncDir(src, dst, SyncOptions{Mode: ModeLink, ConflictPolicy: ConflictOverwrite}); err == nil {
			t.Fatalf("expected link error")
		}
	})

	t.Run("copy", func(t *testing.T) {
		h := snapshotHooks()
		defer h.restore()
		root := t.TempDir()
		src := filepath.Join(root, "src")
		dst := filepath.Join(root, "dst")
		_ = os.MkdirAll(src, 0o755)
		_ = os.WriteFile(filepath.Join(src, "a.txt"), []byte("a"), 0o644)
		osCreate = func(string) (*os.File, error) { return nil, errors.New("create") }
		if err := SyncDir(src, dst, SyncOptions{Mode: ModeCopy, ConflictPolicy: ConflictOverwrite}); err == nil {
			t.Fatalf("expected copy error")
		}
	})
}

func TestFileHashCopyError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file.txt")
	_ = os.WriteFile(path, []byte("data"), 0o644)
	h := snapshotHooks()
	defer h.restore()
	ioCopy = func(io.Writer, io.Reader) (int64, error) { return 0, errors.New("copy") }
	if _, err := FileHash(path); err == nil {
		t.Fatalf("expected copy error")
	}
}

func TestFileHashOpenError(t *testing.T) {
	if _, err := FileHash(filepath.Join(t.TempDir(), "missing")); err == nil {
		t.Fatalf("expected open error")
	}
}

func TestApplyConflictPolicyRemoveError(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "file.txt")
	_ = os.WriteFile(path, []byte("data"), 0o644)
	h := snapshotHooks()
	defer h.restore()
	osRemoveAll = func(string) error { return errors.New("remove") }
	if err := applyConflictPolicy(path, SyncOptions{ConflictPolicy: ConflictOverwrite}); err == nil {
		t.Fatalf("expected remove error")
	}
}

func TestListGitReposDuplicate(t *testing.T) {
	h := snapshotHooks()
	defer h.restore()
	walkDir = func(root string, fn fs.WalkDirFunc) error {
		dir := root
		_ = fn(filepath.Join(dir, ".git"), fakeDirEntry{name: ".git", dir: true}, nil)
		_ = fn(filepath.Join(dir, ".git"), fakeDirEntry{name: ".git", dir: true}, nil)
		return nil
	}
	repos, err := ListGitRepos("root")
	if err != nil {
		t.Fatalf("list error: %v", err)
	}
	if len(repos) != 1 {
		t.Fatalf("expected single repo")
	}
}

type fakeDirEntry struct {
	name string
	dir  bool
}

func (d fakeDirEntry) Name() string               { return d.name }
func (d fakeDirEntry) IsDir() bool                { return d.dir }
func (d fakeDirEntry) Type() fs.FileMode          { return 0 }
func (d fakeDirEntry) Info() (fs.FileInfo, error) { return nil, nil }
