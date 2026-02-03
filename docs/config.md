# Config

Config path:

- `~/.config/github-kanri/config.json`

Related files:

- `config.example.json`
- `docs/config.schema.json`

## Defaults

`gkn config init` writes the default config. Paths accept `~` and are expanded on load.

## Fields

- `projectsRoot` (string, required): base directory for repos/skills.
- `reposRoot` (string, required): repo root (usually `projectsRoot/repos`).
- `skillsRoot` (string, required): skills root (usually `projectsRoot/skills`).
- `skillsRemote` (string, optional): remote URL for skills clone/sync.
- `skillTargets` (string[], optional): relative destinations for skills sync.
- `syncTargets` (object[], required): sync definitions.
  - `name` (string): label.
  - `src` (string): source path.
  - `dest` (string[]): destination paths.
  - `include` (string[]): include globs.
  - `exclude` (string[]): exclude globs.
- `allowCommands` (string[], optional): allowed command globs.
- `denyCommands` (string[], optional): denied command globs (checked first).
- `allowPaths` (string[], optional): allowed path globs.
- `denyPaths` (string[], optional): denied path globs (checked first).
- `syncMode` (string, required): `copy` | `mirror` | `link`.
- `conflictPolicy` (string, required): `fail` | `overwrite`.

## Safety rules

- Deny rules are evaluated before allow rules.
- Ambiguous repo matches stop and print candidates only.
- Destructive actions require explicit `--force`.

## Tips

- Keep `allowCommands` narrow. Prefer explicit entries over wildcards.
- Use `gkn config validate` after edits.
