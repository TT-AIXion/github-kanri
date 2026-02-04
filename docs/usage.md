# Usage

## Quickstart

```sh
gkn config init
gkn config show
gkn repo list
gkn repo status
gkn repo recent --limit 10
gkn shell install --shell zsh
```

## Commands

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

## Common flags

- `--json` output JSON
- `--only <glob>` repeatable
- `--exclude <glob>` repeatable
- `--dry-run`
- `--force`

## Examples

```sh
gkn repo list
gkn repo status --only "**/github-kanri"
gkn repo exec --cmd "git status" --parallel 4

gkn shell install --shell zsh
gkn cd github-kanri

gkn skills sync
gkn skills diff --only "**/skills"

gkn config validate
```

## Shell integration

`gkn shell install --shell zsh` adds a wrapper so `gkn cd <pattern>` changes directories.

## Output

- Human-readable by default
- `--json` for machine-readable output

## Exit codes

- `0` success
- `1` error

## Shell completions

- Bash: `source completions/gkn.bash`
- Zsh: add `completions/` to `$fpath` and run `compinit`

## Manpage

```sh
man ./docs/gkn.1
```
