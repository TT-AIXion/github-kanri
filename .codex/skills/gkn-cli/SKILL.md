---
name: gkn-cli
description: Complete command reference and usage guidance for the gkn (github-kanri) CLI. Use when handling gkn repo management, skills sync/verification, config init/show/validate, doctor/version, or when mapping tasks to the correct gkn commands, flags, and examples.
---

# gkn CLI

## Overview

Use this skill to map user goals to exact gkn commands and flags, including global `--json`, repo selection patterns, and skills sync workflows.

## Global Behavior

- Use `gkn help`, `gkn -h`, or `gkn --help` to show top-level usage.
- Use `gkn <group> --help` or `gkn <group> <command> --help` for details.
- Use global `--json` to emit machine-readable output for supported commands.
- Treat `<pattern>` as substring match on repo name unless it contains `*` or `?` (glob match); use `--pick n` when multiple matches exist.
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
