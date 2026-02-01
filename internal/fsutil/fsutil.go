package fsutil

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/AIXion-Team/github-kanri/internal/match"
)

type SyncMode string

type ConflictPolicy string

const (
	ModeCopy   SyncMode = "copy"
	ModeMirror SyncMode = "mirror"
	ModeLink   SyncMode = "link"
)

const (
	ConflictFail      ConflictPolicy = "fail"
	ConflictOverwrite ConflictPolicy = "overwrite"
)

var (
	osMkdirAll  = os.MkdirAll
	osOpen      = os.Open
	osCreate    = os.Create
	osChmod     = os.Chmod
	osSymlink   = os.Symlink
	osRemoveAll = os.RemoveAll
	osLstat     = os.Lstat
	walkDir     = filepath.WalkDir
	relPath     = filepath.Rel
	ioCopy      = io.Copy
	fileStat    = func(f *os.File) (os.FileInfo, error) { return f.Stat() }
)

type SyncOptions struct {
	Mode           SyncMode
	ConflictPolicy ConflictPolicy
	DryRun         bool
	Include        []string
	Exclude        []string
}

func IsGitRepo(path string) bool {
	info, err := os.Stat(filepath.Join(path, ".git"))
	return err == nil && info.IsDir()
}

func ListGitRepos(root string) ([]string, error) {
	var repos []string
	seen := make(map[string]struct{})
	err := walkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		name := d.Name()
		if name == ".git" {
			repoPath := filepath.Dir(path)
			if _, ok := seen[repoPath]; !ok {
				repos = append(repos, repoPath)
				seen[repoPath] = struct{}{}
			}
			return fs.SkipDir
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(repos)
	return repos, nil
}

func FilterNames(names []string, only []string, exclude []string) []string {
	var out []string
	for _, name := range names {
		if len(only) > 0 && !match.Any(only, name) {
			continue
		}
		if len(exclude) > 0 && match.Any(exclude, name) {
			continue
		}
		out = append(out, name)
	}
	return out
}

func ListFiles(root string, include []string, exclude []string) ([]string, error) {
	var files []string
	if len(include) == 0 {
		include = []string{"**/*"}
	}
	err := walkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == root {
			return nil
		}
		rel, err := relPath(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			if match.Any(exclude, rel) || match.Any(exclude, rel+"/") {
				return fs.SkipDir
			}
			return nil
		}
		if match.Any(exclude, rel) {
			return nil
		}
		if !match.Any(include, rel) {
			return nil
		}
		files = append(files, rel)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	return files, nil
}

func EnsureDir(path string, dryRun bool) error {
	if dryRun {
		return nil
	}
	return osMkdirAll(path, 0o755)
}

func CopyFile(src, dst string, dryRun bool) error {
	if dryRun {
		return nil
	}
	if err := osMkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	srcFile, err := osOpen(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	info, err := fileStat(srcFile)
	if err != nil {
		return err
	}
	dstFile, err := osCreate(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	if _, err := ioCopy(dstFile, srcFile); err != nil {
		return err
	}
	return osChmod(dst, info.Mode())
}

func LinkFile(src, dst string, dryRun bool) error {
	if dryRun {
		return nil
	}
	if err := osMkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return osSymlink(src, dst)
}

func RemovePath(path string, dryRun bool) error {
	if dryRun {
		return nil
	}
	return osRemoveAll(path)
}

func FileHash(path string) (string, error) {
	file, err := osOpen(path)
	if err != nil {
		return "", err
	}
	defer file.Close()
	h := sha256.New()
	if _, err := ioCopy(h, file); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func SyncDir(srcRoot, destRoot string, opts SyncOptions) error {
	files, err := ListFiles(srcRoot, opts.Include, opts.Exclude)
	if err != nil {
		return err
	}
	srcRoot = filepath.Clean(srcRoot)
	destRoot = filepath.Clean(destRoot)

	if opts.Mode == ModeLink {
		for _, rel := range files {
			src := filepath.Join(srcRoot, filepath.FromSlash(rel))
			dst := filepath.Join(destRoot, filepath.FromSlash(rel))
			if err := applyConflictPolicy(dst, opts); err != nil {
				return err
			}
			if err := LinkFile(src, dst, opts.DryRun); err != nil {
				return err
			}
		}
		return nil
	}

	for _, rel := range files {
		src := filepath.Join(srcRoot, filepath.FromSlash(rel))
		dst := filepath.Join(destRoot, filepath.FromSlash(rel))
		if err := applyConflictPolicy(dst, opts); err != nil {
			return err
		}
		if err := CopyFile(src, dst, opts.DryRun); err != nil {
			return err
		}
	}
	if opts.Mode == ModeMirror {
		return mirrorCleanup(srcRoot, destRoot, files, opts)
	}
	return nil
}

func applyConflictPolicy(dst string, opts SyncOptions) error {
	if _, err := osLstat(dst); err != nil {
		return nil
	}
	if opts.ConflictPolicy == ConflictOverwrite {
		return RemovePath(dst, opts.DryRun)
	}
	return fmt.Errorf("conflict detected: %s", dst)
}

func mirrorCleanup(srcRoot, destRoot string, keep []string, opts SyncOptions) error {
	if opts.ConflictPolicy != ConflictOverwrite {
		return errors.New("mirror requires overwrite policy")
	}
	keepSet := make(map[string]struct{}, len(keep))
	for _, rel := range keep {
		keepSet[filepath.ToSlash(rel)] = struct{}{}
	}
	return walkDir(destRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == destRoot {
			return nil
		}
		rel, err := relPath(destRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			return nil
		}
		if _, ok := keepSet[rel]; !ok {
			return RemovePath(path, opts.DryRun)
		}
		return nil
	})
}

func DiffDir(srcRoot, destRoot string, include []string, exclude []string) (added, removed, changed []string, err error) {
	srcFiles, err := ListFiles(srcRoot, include, exclude)
	if err != nil {
		return nil, nil, nil, err
	}
	destFiles, err := ListFiles(destRoot, include, exclude)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return srcFiles, nil, nil, nil
		}
		return nil, nil, nil, err
	}
	srcSet := make(map[string]struct{}, len(srcFiles))
	destSet := make(map[string]struct{}, len(destFiles))
	for _, f := range srcFiles {
		srcSet[f] = struct{}{}
	}
	for _, f := range destFiles {
		destSet[f] = struct{}{}
	}
	for _, f := range srcFiles {
		if _, ok := destSet[f]; !ok {
			added = append(added, f)
			continue
		}
		srcHash, err := FileHash(filepath.Join(srcRoot, filepath.FromSlash(f)))
		if err != nil {
			return nil, nil, nil, err
		}
		destHash, err := FileHash(filepath.Join(destRoot, filepath.FromSlash(f)))
		if err != nil {
			return nil, nil, nil, err
		}
		if srcHash != destHash {
			changed = append(changed, f)
		}
	}
	for _, f := range destFiles {
		if _, ok := srcSet[f]; !ok {
			removed = append(removed, f)
		}
	}
	sort.Strings(added)
	sort.Strings(removed)
	sort.Strings(changed)
	return added, removed, changed, nil
}

func CleanDir(destRoot string, keep []string, dryRun bool) error {
	keepSet := make(map[string]struct{}, len(keep))
	for _, rel := range keep {
		keepSet[filepath.ToSlash(rel)] = struct{}{}
	}
	return walkDir(destRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if path == destRoot {
			return nil
		}
		rel, err := relPath(destRoot, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if d.IsDir() {
			return nil
		}
		if _, ok := keepSet[rel]; !ok {
			return RemovePath(path, dryRun)
		}
		return nil
	})
}

func ResolvePath(repoRoot, dest string) string {
	dest = strings.TrimSpace(dest)
	if dest == "" {
		return repoRoot
	}
	if filepath.IsAbs(dest) {
		return dest
	}
	return filepath.Join(repoRoot, dest)
}
