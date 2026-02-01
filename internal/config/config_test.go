package config

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDefaultConfigPath(t *testing.T) {
	tmp := t.TempDir()
	SetUserHomeDirForTest(func() (string, error) { return tmp, nil })
	defer ResetUserHomeDirForTest()
	path, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("DefaultConfigPath error: %v", err)
	}
	if !strings.Contains(path, filepath.Join(tmp, ".config", "github-kanri", "config.json")) {
		t.Fatalf("unexpected path: %s", path)
	}
}

func TestDefaultConfigPathError(t *testing.T) {
	SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	defer ResetUserHomeDirForTest()
	if _, err := DefaultConfigPath(); err == nil {
		t.Fatalf("expected error")
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg, err := DefaultConfig()
	if err != nil {
		t.Fatalf("DefaultConfig error: %v", err)
	}
	if cfg.ProjectsRoot == "" || cfg.ReposRoot == "" || cfg.SkillsRoot == "" {
		t.Fatalf("missing roots")
	}
	if len(cfg.SyncTargets) == 0 {
		t.Fatalf("missing sync targets")
	}
}

func TestSaveLoadExpandValidate(t *testing.T) {
	tmp := t.TempDir()
	SetUserHomeDirForTest(func() (string, error) { return tmp, nil })
	defer ResetUserHomeDirForTest()
	cfg := Config{
		ProjectsRoot:   "~/Projects",
		ReposRoot:      "~/Projects/repos",
		SkillsRoot:     "~/Projects/skills",
		SkillTargets:   []string{".codex/skills"},
		AllowCommands:  []string{"git status*"},
		DenyCommands:   []string{"rm -rf*"},
		SyncMode:       "copy",
		ConflictPolicy: "fail",
		SyncTargets: []SyncTarget{{
			Name:    "skills",
			Src:     "~/Projects/skills",
			Dest:    []string{".codex/skills"},
			Include: []string{"**/*"},
			Exclude: []string{".git/**"},
		}},
	}
	path, err := DefaultConfigPath()
	if err != nil {
		t.Fatalf("path error: %v", err)
	}
	if err := Save(path, cfg); err != nil {
		t.Fatalf("save error: %v", err)
	}
	loaded, err := Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	loaded, err = ExpandConfigPaths(loaded)
	if err != nil {
		t.Fatalf("expand error: %v", err)
	}
	if errs := Validate(loaded); len(errs) != 0 {
		t.Fatalf("validate error: %v", errs)
	}
}

func TestLoadError(t *testing.T) {
	if _, err := Load(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Fatalf("expected load error")
	}
}

func TestSaveError(t *testing.T) {
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "nope")
	if err := os.MkdirAll(dir, 0o000); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path := filepath.Join(dir, "config.json")
	if err := Save(path, Config{}); err == nil {
		t.Fatalf("expected save error")
	}
	_ = os.Chmod(dir, 0o755)

	fileParent := filepath.Join(tmp, "file-parent")
	if err := os.WriteFile(fileParent, []byte("x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := Save(filepath.Join(fileParent, "config.json"), Config{}); err == nil {
		t.Fatalf("expected mkdir error")
	}

	dir2 := filepath.Join(tmp, "readonly")
	if err := os.MkdirAll(dir2, 0o555); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	path2 := filepath.Join(dir2, "config.json")
	if err := Save(path2, Config{}); err == nil {
		t.Fatalf("expected write error")
	}
	_ = os.Chmod(dir2, 0o755)
}

func TestSaveMarshalError(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.json")
	orig := jsonMarshalIndent
	jsonMarshalIndent = func(interface{}, string, string) ([]byte, error) {
		return nil, errors.New("marshal")
	}
	defer func() { jsonMarshalIndent = orig }()
	if err := Save(path, Config{}); err == nil {
		t.Fatalf("expected marshal error")
	}
}

func TestExpandPathError(t *testing.T) {
	SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	defer ResetUserHomeDirForTest()
	if _, err := ExpandPath("~/Projects"); err == nil {
		t.Fatalf("expected error")
	}
}

func TestExpandPathEmpty(t *testing.T) {
	if got, err := ExpandPath(" "); err != nil || got != "" {
		t.Fatalf("expected empty")
	}
}

func TestExpandConfigPathsError(t *testing.T) {
	SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
	defer ResetUserHomeDirForTest()
	cfg := Config{
		ProjectsRoot: "~/Projects",
		ReposRoot:    "~/Projects/repos",
		SkillsRoot:   "~/Projects/skills",
		AllowPaths:   []string{"~/allow"},
		DenyPaths:    []string{"~/deny"},
		SyncTargets: []SyncTarget{{
			Name: "skills",
			Src:  "~/Projects/skills",
			Dest: []string{"~/dest"},
		}},
	}
	if _, err := ExpandConfigPaths(cfg); err == nil {
		t.Fatalf("expected expand error")
	}
}

func TestExpandConfigPathsErrorStages(t *testing.T) {
	t.Run("repos", func(t *testing.T) {
		SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
		defer ResetUserHomeDirForTest()
		cfg := Config{
			ProjectsRoot: "/tmp",
			ReposRoot:    "~/repos",
			SkillsRoot:   "/tmp/skills",
		}
		if _, err := ExpandConfigPaths(cfg); err == nil {
			t.Fatalf("expected repos error")
		}
	})
	t.Run("skills", func(t *testing.T) {
		SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
		defer ResetUserHomeDirForTest()
		cfg := Config{
			ProjectsRoot: "/tmp",
			ReposRoot:    "/tmp/repos",
			SkillsRoot:   "~/skills",
		}
		if _, err := ExpandConfigPaths(cfg); err == nil {
			t.Fatalf("expected skills error")
		}
	})
	t.Run("allow", func(t *testing.T) {
		SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
		defer ResetUserHomeDirForTest()
		cfg := Config{
			ProjectsRoot: "/tmp",
			ReposRoot:    "/tmp/repos",
			SkillsRoot:   "/tmp/skills",
			AllowPaths:   []string{"~/allow"},
		}
		if _, err := ExpandConfigPaths(cfg); err == nil {
			t.Fatalf("expected allow error")
		}
	})
	t.Run("deny", func(t *testing.T) {
		SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
		defer ResetUserHomeDirForTest()
		cfg := Config{
			ProjectsRoot: "/tmp",
			ReposRoot:    "/tmp/repos",
			SkillsRoot:   "/tmp/skills",
			DenyPaths:    []string{"~/deny"},
		}
		if _, err := ExpandConfigPaths(cfg); err == nil {
			t.Fatalf("expected deny error")
		}
	})
	t.Run("sync-src", func(t *testing.T) {
		SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
		defer ResetUserHomeDirForTest()
		cfg := Config{
			ProjectsRoot: "/tmp",
			ReposRoot:    "/tmp/repos",
			SkillsRoot:   "/tmp/skills",
			SyncTargets: []SyncTarget{{
				Name: "skills",
				Src:  "~/src",
				Dest: []string{".codex/skills"},
			}},
		}
		if _, err := ExpandConfigPaths(cfg); err == nil {
			t.Fatalf("expected sync src error")
		}
	})
	t.Run("sync-dest", func(t *testing.T) {
		SetUserHomeDirForTest(func() (string, error) { return "", os.ErrPermission })
		defer ResetUserHomeDirForTest()
		cfg := Config{
			ProjectsRoot: "/tmp",
			ReposRoot:    "/tmp/repos",
			SkillsRoot:   "/tmp/skills",
			SyncTargets: []SyncTarget{{
				Name: "skills",
				Src:  "/tmp/src",
				Dest: []string{"~/dest"},
			}},
		}
		if _, err := ExpandConfigPaths(cfg); err == nil {
			t.Fatalf("expected sync dest error")
		}
	})
}

func TestExpandConfigPathsSuccess(t *testing.T) {
	tmp := t.TempDir()
	SetUserHomeDirForTest(func() (string, error) { return tmp, nil })
	defer ResetUserHomeDirForTest()
	cfg := Config{
		ProjectsRoot: "~/Projects",
		ReposRoot:    "~/Projects/repos",
		SkillsRoot:   "~/Projects/skills",
		AllowPaths:   []string{"~/allow"},
		DenyPaths:    []string{"~/deny"},
		SyncTargets: []SyncTarget{{
			Name: "skills",
			Src:  "~/Projects/skills",
			Dest: []string{"~/dest1", "~/dest2"},
		}},
	}
	if _, err := ExpandConfigPaths(cfg); err != nil {
		t.Fatalf("expand error: %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmp := t.TempDir()
	path := filepath.Join(tmp, "config.json")
	if err := os.WriteFile(path, []byte("{bad"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if _, err := Load(path); err == nil {
		t.Fatalf("expected unmarshal error")
	}
}

func TestApplyDefaults(t *testing.T) {
	cfg := ApplyDefaults(Config{})
	if cfg.SyncMode == "" || cfg.ConflictPolicy == "" {
		t.Fatalf("defaults missing")
	}
	if len(cfg.DenyCommands) == 0 {
		t.Fatalf("deny defaults missing")
	}
}

func TestValidateErrors(t *testing.T) {
	cfg := Config{SyncMode: "bad", ConflictPolicy: "bad"}
	errs := Validate(cfg)
	if len(errs) == 0 {
		t.Fatalf("expected errors")
	}
}

func TestValidateMissingTargets(t *testing.T) {
	cfg := Config{ProjectsRoot: "x", ReposRoot: "y", SkillsRoot: "z", SyncMode: "copy", ConflictPolicy: "fail", SyncTargets: []SyncTarget{{}}}
	if errs := Validate(cfg); len(errs) == 0 {
		t.Fatalf("expected errors")
	}
}
