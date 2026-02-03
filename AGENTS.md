# Repository Guidelines

## Project Structure & Module Organization

- `README.md`: Product overview and CLI concept.
- `.codex/skills/`: Persistent project docs for agents and maintainers.
- `LICENSE`: Licensing details.
- Source code and tests are not committed yet. When adding them, document the chosen paths here (e.g., `cmd/gkn/`, `internal/`, `pkg/`, `testdata/`).

## Build, Test, and Development Commands

- No build or test scripts are defined yet.
- Planned Go workflow (once `go.mod` exists):
  - `go build ./...` for a full build.
  - `go test ./...` for all tests.
  - `gofmt -w <path>` for formatting.

## Coding Style & Naming Conventions

- Language: Go (planned).
- Formatting: `gofmt` is required; default Go indentation (tabs).
- Naming: short, clear identifiers; prefer explicit over clever.
- Dependencies: keep minimal and justified.

## Testing Guidelines

- Tests are not implemented yet.
- When added, use Go's standard testing package and name files `*_test.go`.
- Prefer regression tests for bugs and run `go test ./...` before delivery.

## Commit & Pull Request Guidelines

- Commit messages: Conventional Commits (e.g., `feat: add config loader`).
- Default branch: `main` (direct commits acceptable unless a PR is requested).
- If using PRs, include a short summary, test status, and any config changes.

## Release (Homebrew)

- Push only is insufficient; Homebrew upgrades require a new tag + Formula bump.
- Update `Formula/gkn.rb`:
  - `version` to new release id (timestamp + short SHA or semver).
  - `url` tag to `v<version>`.
  - `revision` to the release commit SHA.
- Create tag on the release commit: `v<version>` and push both commit + tag.
- Verify: `brew update` then `brew upgrade gkn` and `gkn version`.

## Security & Configuration Tips

- Configuration lives at `~/.config/github-kanri/config.json`.
- Respect allow/deny command rules; deny should fail fast.
- Avoid introducing secrets; this project expects local, non-networked operation.

## References

- CLI name: `gkn`.
