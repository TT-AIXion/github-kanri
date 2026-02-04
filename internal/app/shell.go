package app

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var errShellUnknown = errors.New("unknown shell")

const shellMarkerStart = "# >>> gkn shell >>>"
const shellMarkerEnd = "# <<< gkn shell <<<"

func (a App) runShell(_ context.Context, args []string) int {
	if len(args) == 0 || isHelp(args[0]) {
		a.Out.Raw(`gkn shell <shell>
gkn shell install --shell <shell> [--profile path] [--force] [--dry-run]

Shells:
  zsh
  bash
  fish
  powershell`)
		return 0
	}
	switch args[0] {
	case "install":
		return a.runShellInstall(args[1:])
	default:
		snippet, err := shellSnippet(args[0])
		if err != nil {
			a.Out.Err(err.Error(), nil)
			return 1
		}
		a.Out.Raw(snippet)
		return 0
	}
}

func (a App) runShellInstall(args []string) int {
	fs := flag.NewFlagSet("shell install", flag.ContinueOnError)
	fs.SetOutput(os.Stdout)
	shell := fs.String("shell", "", "shell name")
	profile := fs.String("profile", "", "profile path")
	force := fs.Bool("force", false, "replace existing block")
	dryRun := fs.Bool("dry-run", false, "show changes only")
	if err := fs.Parse(args); err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if *shell == "" {
		a.Out.Err("shell required", nil)
		return 1
	}
	snippet, err := shellSnippet(*shell)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	target := *profile
	if target == "" {
		target, err = shellDefaultProfile(*shell)
		if err != nil {
			a.Out.Err(err.Error(), nil)
			return 1
		}
	}
	updated, err := installShellSnippet(target, snippet, *force, *dryRun)
	if err != nil {
		a.Out.Err(err.Error(), nil)
		return 1
	}
	if a.Out.JSON {
		a.Out.OK("shell install", map[string]interface{}{
			"profile": target,
			"updated": updated,
			"dryRun":  *dryRun,
		})
		return 0
	}
	if *dryRun {
		a.Out.OK(fmt.Sprintf("dry-run: %s", target), nil)
		return 0
	}
	if updated {
		a.Out.OK(fmt.Sprintf("installed: %s", target), nil)
		return 0
	}
	a.Out.OK(fmt.Sprintf("already installed: %s", target), nil)
	return 0
}

func shellSnippet(shell string) (string, error) {
	switch shell {
	case "zsh", "bash":
		return strings.TrimSpace(zshBashSnippet), nil
	case "fish":
		return strings.TrimSpace(fishSnippet), nil
	case "powershell":
		return strings.TrimSpace(powershellSnippet), nil
	default:
		return "", fmt.Errorf("%w: %s", errShellUnknown, shell)
	}
}

func shellDefaultProfile(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch shell {
	case "zsh":
		if zdot := os.Getenv("ZDOTDIR"); zdot != "" {
			return filepath.Join(zdot, ".zshrc"), nil
		}
		return filepath.Join(home, ".zshrc"), nil
	case "bash":
		bashrc := filepath.Join(home, ".bashrc")
		if _, err := os.Stat(bashrc); err == nil {
			return bashrc, nil
		}
		bashProfile := filepath.Join(home, ".bash_profile")
		if _, err := os.Stat(bashProfile); err == nil {
			return bashProfile, nil
		}
		return bashrc, nil
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish"), nil
	case "powershell":
		if runtime.GOOS == "windows" {
			return filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"), nil
		}
		return filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"), nil
	default:
		return "", fmt.Errorf("%w: %s", errShellUnknown, shell)
	}
}

func installShellSnippet(profile string, snippet string, force bool, dryRun bool) (bool, error) {
	content, err := os.ReadFile(profile)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	updated, next, err := updateShellContent(string(content), snippet, force)
	if err != nil {
		return false, err
	}
	if !updated || dryRun {
		return updated, nil
	}
	if err := os.MkdirAll(filepath.Dir(profile), 0o755); err != nil {
		return false, err
	}
	if err := os.WriteFile(profile, []byte(next), 0o644); err != nil {
		return false, err
	}
	return updated, nil
}

func updateShellContent(content string, snippet string, force bool) (bool, string, error) {
	block := shellMarkerStart + "\n" + strings.TrimSpace(snippet) + "\n" + shellMarkerEnd + "\n"
	start := strings.Index(content, shellMarkerStart)
	if start == -1 {
		if content != "" && !strings.HasSuffix(content, "\n") {
			content += "\n"
		}
		return true, content + block, nil
	}
	end := strings.Index(content[start:], shellMarkerEnd)
	if end == -1 {
		return false, "", fmt.Errorf("shell marker end missing")
	}
	end += start + len(shellMarkerEnd)
	if !force {
		return false, content, nil
	}
	next := content[:start] + block + strings.TrimPrefix(content[end:], "\n")
	return true, next, nil
}

const zshBashSnippet = `
gkn() {
  if [ "$1" = "cd" ]; then
    local p
    p="$(command gkn cd "${@:2}")" || return
    [ -n "$p" ] || return 1
    cd "$p" || return
    return
  fi
  command gkn "$@"
}
`

const fishSnippet = `
function gkn
  if test (count $argv) -ge 1; and test $argv[1] = "cd"
    set -l p (command gkn cd $argv[2..-1])
    if test $status -ne 0
      return $status
    end
    if test -n "$p"
      cd "$p"
    end
    return
  end
  command gkn $argv
end
`

const powershellSnippet = `
function gkn {
  param([Parameter(ValueFromRemainingArguments = $true)][string[]]$Args)
  $exe = (Get-Command gkn -CommandType Application).Source
  if ($Args.Length -gt 0 -and $Args[0] -eq 'cd') {
    $rest = @()
    if ($Args.Length -gt 1) { $rest = $Args[1..($Args.Length - 1)] }
    $path = & $exe cd @rest
    if ($LASTEXITCODE -ne 0) { return }
    if ($path) { Set-Location $path }
    return
  }
  & $exe @Args
}
`
