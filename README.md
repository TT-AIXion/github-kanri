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
gkn repo <list|status|recent|info|graph|open|path|clone|exec>
gkn skills <clone|sync|link|watch|diff|verify|status|pin|clean>
gkn config <init|show|validate>
gkn doctor
gkn version
```

## Safety model

- Deny rules are always checked first
- Ambiguous matches never execute; candidates only
- Destructive actions require explicit `--force`

## Config

Single JSON config file:

- `~/.config/github-kanri/config.json`

## Requirements

- macOS
- Git

## Docs

- `README.ja.md`
- `docs/requirements.md`
- `CONTRIBUTING.md`
- `SECURITY.md`
- `SUPPORT.md`

## License

MIT. See `LICENSE`.
