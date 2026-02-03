---
name: gkn-cli
description: Complete command reference and usage guidance for the gkn (github-kanri) CLI. Use when handling gkn repo management, skills sync/verification, config init/show/validate, doctor/version, or when mapping tasks to the correct gkn commands, flags, and examples.
---

# gkn CLI

## Overview

Use this skill to map user goals to exact gkn commands and flags, including global `--json`, repo selection patterns, and skills sync workflows.

## Intent To Command

- Initialize config: `gkn config init` (use `--force` to overwrite).
- Show config: `gkn config show`.
- Validate config: `gkn config validate`.
- Clone repo into `reposRoot`: `gkn clone <url>` or `gkn repo clone <url>`.
- List repos: `gkn repo list`.
- Check dirty/clean: `gkn repo status`.
- Show recent activity: `gkn repo recent --limit n`.
- Open repo in VS Code: `gkn repo open <pattern> [--pick n]`.
- Get repo path: `gkn repo path <pattern> [--pick n]`.
- Show repo info: `gkn repo info <pattern> [--pick n]`.
- Show oneline log: `gkn repo graph <pattern> [--pick n] [--limit n]`.
- Run command across repos: `gkn repo exec --cmd "<command>"`.
- Clone or update skills repo: `gkn skills clone [--remote url]`.
- Sync skills: `gkn skills sync [--mode copy|mirror|link]`.
- Link skills: `gkn skills link`.
- Check drift: `gkn skills diff` or `gkn skills status`.
- Verify match: `gkn skills verify` (exit code `2` on mismatch).
- Pin skills to ref: `gkn skills pin --target name --ref <commit|tag>`.
- Clean extras: `gkn skills clean --force`.
- Watch and auto-sync: `gkn skills watch --interval sec`.
- Environment check: `gkn doctor`.
- Show version: `gkn version`.

## Global Behavior

- Use `gkn help`, `gkn -h`, or `gkn --help` to show top-level usage.
- Use `gkn <group> --help` or `gkn <group> <command> --help` for details.
- Use global `--json` to emit JSON envelopes (`level`, `message`, `data`).
- Treat `<pattern>` as substring match on repo name unless it contains `*` or `?` (glob match).
- Use `--pick n` (1-based) when multiple matches exist.
- Use `--only` and `--exclude` as repeatable or comma-separated globs.
- Note `gkn repo open` runs `code <path>` on the selected repo.

## Commands (Complete)

```text
gkn [--json] <command>

gkn clone <url> [--name repo]

gkn repo list [--only glob] [--exclude glob]
gkn repo status [--only glob] [--exclude glob]
gkn repo open <pattern> [--pick n]
gkn repo path <pattern> [--pick n]
gkn repo recent [--limit n] [--only glob] [--exclude glob]
gkn repo info <pattern> [--pick n]
gkn repo graph <pattern> [--pick n] [--limit n]
gkn repo clone <url> [--name repo]
gkn repo exec --cmd "<command>" [--parallel n] [--timeout sec] [--require-clean] [--dry-run] [--only glob] [--exclude glob]

gkn skills clone [--remote url] [--force]
gkn skills sync [--target name] [--mode copy|mirror|link] [--force] [--dry-run] [--only glob] [--exclude glob]
gkn skills link [--target name] [--force] [--dry-run] [--only glob] [--exclude glob]
gkn skills watch [--target name] [--interval sec]
gkn skills diff [--target name] [--only glob] [--exclude glob]
gkn skills verify [--target name] [--only glob] [--exclude glob]
gkn skills status [--target name] [--only glob] [--exclude glob]
gkn skills pin --target name --ref <commit|tag> [--force]
gkn skills clean [--target name] [--force] [--dry-run] [--only glob] [--exclude glob]

gkn config show
gkn config init [--force]
gkn config validate
gkn doctor
gkn version
```

## Command Notes

- `clone` is an alias for `repo clone` and clones into `reposRoot`.
- `repo list` outputs `name path` for each repo; `repo status` reports `clean/dirty`.
- `repo recent` sorts by last commit time; repos with no commits are flagged.
- `repo info` shows origin/current/default/dirty; `repo graph` shows oneline log.
- `repo exec` runs a shell command per repo with optional parallelism, timeout, `--require-clean`, and `--dry-run`.
- `skills clone` uses `skillsRemote` from config unless `--remote` is provided; pulls if `skillsRoot` already is a git repo.
- `skills sync` applies configured `syncTargets`; `skills link` forces link mode.
- `skills diff`/`status` show drift; `skills verify` exits with code `2` on mismatch.
- `skills pin` requires clean `skillsRoot` unless `--force`.
- `skills clean` requires `--force` or `conflictPolicy=overwrite` and removes files not in `src` (supports `--dry-run`).
- `config show/init/validate` operate on `~/.config/github-kanri/config.json`.

## Exit Codes

- `skills verify` returns `2` on mismatch.
- Most commands return `1` on error and `0` on success.
- `repo exec` returns `1` if any execution error occurs; `--require-clean` skips dirty repos without failing.

## Examples

```text
gkn --json repo list
gkn repo status --only "foo*" --exclude "tmp*"
gkn repo open api --pick 2
gkn repo exec --cmd "git status -sb" --parallel 4
gkn skills sync --mode mirror --dry-run
gkn skills verify --target skills
gkn skills clean --force --dry-run
```
