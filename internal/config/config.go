package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var userHomeDir = os.UserHomeDir
var jsonMarshalIndent = json.MarshalIndent

func SetUserHomeDirForTest(fn func() (string, error)) {
	userHomeDir = fn
}

func ResetUserHomeDirForTest() {
	userHomeDir = os.UserHomeDir
}

type Config struct {
	ProjectsRoot   string       `json:"projectsRoot"`
	ReposRoot      string       `json:"reposRoot"`
	SkillsRoot     string       `json:"skillsRoot"`
	SkillsRemote   string       `json:"skillsRemote,omitempty"`
	SkillTargets   []string     `json:"skillTargets"`
	SyncTargets    []SyncTarget `json:"syncTargets"`
	AllowCommands  []string     `json:"allowCommands"`
	DenyCommands   []string     `json:"denyCommands"`
	AllowPaths     []string     `json:"allowPaths,omitempty"`
	DenyPaths      []string     `json:"denyPaths,omitempty"`
	SyncMode       string       `json:"syncMode"`
	ConflictPolicy string       `json:"conflictPolicy"`
}

type SyncTarget struct {
	Name    string   `json:"name"`
	Src     string   `json:"src"`
	Dest    []string `json:"dest"`
	Include []string `json:"include"`
	Exclude []string `json:"exclude"`
}

func DefaultConfigPath() (string, error) {
	home, err := userHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "github-kanri", "config.json"), nil
}

func DefaultConfig() (Config, error) {
	projects := "~/Projects"
	repos := filepath.Join(projects, "repos")
	skills := filepath.Join(projects, "skills")
	skillTargets := []string{".codex/skills", ".claude/skills"}
	return Config{
		ProjectsRoot: projects,
		ReposRoot:    repos,
		SkillsRoot:   skills,
		SkillsRemote: "",
		SkillTargets: skillTargets,
		SyncTargets: []SyncTarget{
			{
				Name:    "skills",
				Src:     skills,
				Dest:    skillTargets,
				Include: []string{"**/*"},
				Exclude: []string{".git/**"},
			},
		},
		AllowCommands: []string{
			"gh auth status*",
			"gh repo create*",
			"git init*",
			"git add*",
			"git commit*",
			"git status*",
			"git log*",
			"git rev-parse*",
			"git config*",
			"git remote*",
			"git clone*",
			"git fetch*",
			"git pull*",
			"git checkout*",
			"git push*",
			"code *",
		},
		DenyCommands: []string{
			"rm -rf*",
			"git reset --hard*",
		},
		SyncMode:       "copy",
		ConflictPolicy: "fail",
	}, nil
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	cfg = ApplyDefaults(cfg)
	return cfg, nil
}

func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := jsonMarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0o644)
}

func ApplyDefaults(cfg Config) Config {
	if cfg.ProjectsRoot == "" {
		cfg.ProjectsRoot = "~/Projects"
	}
	if cfg.ReposRoot == "" {
		cfg.ReposRoot = filepath.Join(cfg.ProjectsRoot, "repos")
	}
	if cfg.SkillsRoot == "" {
		cfg.SkillsRoot = filepath.Join(cfg.ProjectsRoot, "skills")
	}
	if len(cfg.SkillTargets) == 0 {
		cfg.SkillTargets = []string{".codex/skills", ".claude/skills"}
	}
	if len(cfg.SyncTargets) == 0 {
		cfg.SyncTargets = []SyncTarget{
			{
				Name:    "skills",
				Src:     cfg.SkillsRoot,
				Dest:    cfg.SkillTargets,
				Include: []string{"**/*"},
				Exclude: []string{".git/**"},
			},
		}
	}
	if cfg.SyncMode == "" {
		cfg.SyncMode = "copy"
	}
	if cfg.ConflictPolicy == "" {
		cfg.ConflictPolicy = "fail"
	}
	if len(cfg.DenyCommands) == 0 {
		cfg.DenyCommands = []string{
			"rm -rf*",
			"git reset --hard*",
		}
	}
	return cfg
}

func ExpandPath(p string) (string, error) {
	p = strings.TrimSpace(p)
	if p == "" {
		return "", nil
	}
	if strings.HasPrefix(p, "~") {
		home, err := userHomeDir()
		if err != nil {
			return "", err
		}
		p = filepath.Join(home, strings.TrimPrefix(p, "~"))
	}
	return filepath.Clean(p), nil
}

func ExpandConfigPaths(cfg Config) (Config, error) {
	var err error
	if cfg.ProjectsRoot, err = ExpandPath(cfg.ProjectsRoot); err != nil {
		return Config{}, err
	}
	if cfg.ReposRoot, err = ExpandPath(cfg.ReposRoot); err != nil {
		return Config{}, err
	}
	if cfg.SkillsRoot, err = ExpandPath(cfg.SkillsRoot); err != nil {
		return Config{}, err
	}
	for i, p := range cfg.AllowPaths {
		if cfg.AllowPaths[i], err = ExpandPath(p); err != nil {
			return Config{}, err
		}
	}
	for i, p := range cfg.DenyPaths {
		if cfg.DenyPaths[i], err = ExpandPath(p); err != nil {
			return Config{}, err
		}
	}
	for i, t := range cfg.SyncTargets {
		if t.Src, err = ExpandPath(t.Src); err != nil {
			return Config{}, err
		}
		for j, d := range t.Dest {
			if t.Dest[j], err = ExpandPath(d); err != nil {
				return Config{}, err
			}
		}
		cfg.SyncTargets[i] = t
	}
	return cfg, nil
}

func Validate(cfg Config) []error {
	var errs []error
	if strings.TrimSpace(cfg.ProjectsRoot) == "" {
		errs = append(errs, fmt.Errorf("projectsRoot is required"))
	}
	if strings.TrimSpace(cfg.ReposRoot) == "" {
		errs = append(errs, fmt.Errorf("reposRoot is required"))
	}
	if strings.TrimSpace(cfg.SkillsRoot) == "" {
		errs = append(errs, fmt.Errorf("skillsRoot is required"))
	}
	if cfg.SyncMode != "copy" && cfg.SyncMode != "mirror" && cfg.SyncMode != "link" {
		errs = append(errs, fmt.Errorf("syncMode must be copy|mirror|link"))
	}
	if cfg.ConflictPolicy != "fail" && cfg.ConflictPolicy != "overwrite" {
		errs = append(errs, fmt.Errorf("conflictPolicy must be fail|overwrite"))
	}
	for i, t := range cfg.SyncTargets {
		if strings.TrimSpace(t.Name) == "" {
			errs = append(errs, fmt.Errorf("syncTargets[%d].name is required", i))
		}
		if strings.TrimSpace(t.Src) == "" {
			errs = append(errs, fmt.Errorf("syncTargets[%d].src is required", i))
		}
		if len(t.Dest) == 0 {
			errs = append(errs, fmt.Errorf("syncTargets[%d].dest is required", i))
		}
	}
	return errs
}
