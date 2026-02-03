package app

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TT-AIXion/github-kanri/internal/executil"
)

type quickstartResult struct {
	Name       string `json:"name"`
	Path       string `json:"path"`
	Visibility string `json:"visibility"`
}

func (a App) runQuickstart(ctx context.Context, args []string) int {
	fs := flag.NewFlagSet("quickstart", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	public := fs.Bool("public", false, "make repo public")
	private := fs.Bool("private", false, "make repo private")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if *public && *private {
		a.Out.Err("use only one of --public or --private", nil)
		return 1
	}
	if fs.NArg() == 0 {
		a.Out.Err("repo name required", nil)
		return 1
	}
	if fs.NArg() > 1 {
		a.Out.Err(fmt.Sprintf("unexpected argument: %s", fs.Arg(1)), nil)
		return 1
	}
	ghName, repoName, err := parseQuickstartName(fs.Arg(0))
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	cfg, _, err := loadConfig()
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := ensureSkillsRoot(cfg.SkillsRoot); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	dest := filepath.Join(cfg.ReposRoot, repoName)
	if _, err := os.Stat(dest); err == nil {
		a.Out.Err("destination exists", nil)
		return 1
	}
	guard := guardFromConfig(cfg)
	if err := guard.CheckPath(dest); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	runner := buildRunner(cfg, false)
	if err := ghAuthStatus(ctx, runner); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := os.MkdirAll(cfg.ReposRoot, 0o755); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := initRepoMain(ctx, runner, dest); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := writeQuickstartReadme(dest, repoName); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	targets, err := selectTargets(cfg, "skills")
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if code := a.syncTargets(ctx, cfg, targets, "", false, false, []string{repoName}, nil); code != 0 {
		return code
	}
	if err := gitCommitInit(ctx, runner, dest); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	visibility := "private"
	visFlag := "--private"
	if *public {
		visibility = "public"
		visFlag = "--public"
	}
	if *private {
		visibility = "private"
		visFlag = "--private"
	}
	if err := ghRepoCreate(ctx, runner, ghName, dest, visFlag); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if err := gitPullPush(ctx, runner, dest); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	result := quickstartResult{Name: repoName, Path: dest, Visibility: visibility}
	a.Out.OK(fmt.Sprintf("quickstart %s %s", repoName, dest), result)
	return 0
}

func parseQuickstartName(arg string) (string, string, error) {
	name := strings.TrimSpace(arg)
	if name == "" {
		return "", "", fmt.Errorf("repo name required")
	}
	name = strings.TrimSuffix(name, ".git")
	repoName := filepath.Base(name)
	if repoName == "." || repoName == string(filepath.Separator) || repoName == "" {
		return "", "", fmt.Errorf("invalid repo name: %s", arg)
	}
	return name, repoName, nil
}

func ensureSkillsRoot(path string) error {
	if path == "" {
		return fmt.Errorf("skillsRoot is empty")
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("skillsRoot not found: %s", path)
		}
		return err
	}
	return nil
}

func ghAuthStatus(ctx context.Context, runner executil.Runner) error {
	res, err := runner.Run(ctx, "", "gh", "auth", "status")
	if err != nil {
		return fmt.Errorf("gh auth status failed: %s", commandError(err, res))
	}
	return nil
}

func initRepoMain(ctx context.Context, runner executil.Runner, dest string) error {
	if _, err := runner.Run(ctx, dest, "git", "init", "-b", "main"); err == nil {
		return nil
	}
	resFallback, errFallback := runner.Run(ctx, dest, "git", "init")
	if errFallback != nil {
		return fmt.Errorf("git init failed: %s", commandError(errFallback, resFallback))
	}
	resCheckout, errCheckout := runner.Run(ctx, dest, "git", "checkout", "-b", "main")
	if errCheckout != nil {
		return fmt.Errorf("git checkout -b main failed: %s", commandError(errCheckout, resCheckout))
	}
	return nil
}

func writeQuickstartReadme(dest string, repoName string) error {
	path := filepath.Join(dest, "README.md")
	if _, err := os.Stat(path); err == nil {
		return nil
	}
	body := fmt.Sprintf("# %s\n", repoName)
	return os.WriteFile(path, []byte(body), 0o644)
}

func gitCommitInit(ctx context.Context, runner executil.Runner, dest string) error {
	resAdd, err := runner.Run(ctx, dest, "git", "add", ".")
	if err != nil {
		return fmt.Errorf("git add failed: %s", commandError(err, resAdd))
	}
	resCommit, err := runner.Run(ctx, dest, "git", "commit", "-m", "chore: init", "--allow-empty")
	if err != nil {
		return fmt.Errorf("git commit failed: %s", commandError(err, resCommit))
	}
	return nil
}

func ghRepoCreate(ctx context.Context, runner executil.Runner, name string, source string, visibility string) error {
	args := []string{"repo", "create", name, visibility, "--confirm", "--source", source, "--remote", "origin", "--push"}
	res, err := runner.Run(ctx, "", "gh", args...)
	if err != nil {
		return fmt.Errorf("gh repo create failed: %s", commandError(err, res))
	}
	return nil
}

func gitPullPush(ctx context.Context, runner executil.Runner, dest string) error {
	resPull, err := runner.Run(ctx, dest, "git", "pull", "--rebase", "origin", "main")
	if err != nil {
		return fmt.Errorf("git pull failed: %s", commandError(err, resPull))
	}
	resPush, err := runner.Run(ctx, dest, "git", "push", "-u", "origin", "main")
	if err != nil {
		return fmt.Errorf("git push failed: %s", commandError(err, resPush))
	}
	return nil
}

func commandError(err error, res executil.Result) string {
	if err == nil {
		return ""
	}
	msg := strings.TrimSpace(res.Stderr)
	if msg == "" {
		msg = strings.TrimSpace(res.Stdout)
	}
	if msg == "" {
		return err.Error()
	}
	return fmt.Sprintf("%s: %s", err.Error(), msg)
}
