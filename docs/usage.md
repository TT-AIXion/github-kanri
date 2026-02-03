# Usage

## Quickstart

```sh
gkn config init
gkn config show
gkn repo list
gkn repo status
gkn repo recent --limit 10
```

## Commands

```text
gkn repo <list|status|recent|info|graph|open|path|clone|exec>
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

gkn skills sync
gkn skills diff --only "**/skills"

gkn config validate
```

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
