# github-kanri / gkn

Safety-first local GitHub repo management CLI.

## Why gkn

- Deny-first safety guard (fail fast on risky commands)
- Non-interactive and reproducible
- Local-only (no GitHub API)
- Fast bulk operations across many repos
- JSON output with `--json`
- Skills sync and verification

## Quickstart

```sh
gkn config init
gkn config show
gkn repo list
gkn repo status
gkn repo recent --limit 10
gkn shell install --shell zsh
```

## Install

Homebrew:

```sh
brew install TT-AIXion/github-kanri/gkn
```

Go install:

```sh
go install github.com/TT-AIXion/github-kanri/cmd/gkn@<tag>
```

Build from source:

```sh
go build -o gkn ./cmd/gkn
```

## Core commands

```text
gkn cd <pattern> [--pick n]
gkn repo <list|status|recent|info|graph|open|path|cd|clone|exec>
gkn shell <shell>
gkn shell install --shell <shell> [--profile path] [--force] [--dry-run]
gkn skills <clone|sync|link|watch|diff|verify|status|pin|clean>
gkn config <init|show|validate>
gkn doctor
gkn version
```

## Safety model

- Deny rules are always checked first
- Ambiguous matches never execute; candidates only
- Destructive actions require explicit `--force`

## Shell integration

`gkn shell install --shell zsh` installs a wrapper so `gkn cd <pattern>` changes directories.

## Config

Single JSON config file:

- `~/.config/github-kanri/config.json`
- Example: `config.example.json`
- Schema: `docs/config.schema.json`

## Requirements

- macOS (primary)
- Linux binaries are provided on a best-effort basis
- Git

## Docs

- `README.ja.md`
- `docs/requirements.md`
- `docs/usage.md`
- `docs/config.md`
- `docs/config.schema.json`
- `docs/RELEASING.md`
- `docs/gkn.1`
- `CONTRIBUTING.md`
- `SECURITY.md`
- `SUPPORT.md`
- `config.example.json`
- `completions/`

## License

MIT. See `LICENSE`.
