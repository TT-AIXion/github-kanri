# github-kanri

GitHub リポジトリ管理 CLI（エージェント/人間両対応）。
コマンドのみ・非対話・安全設計。ローカル運用前提。

## 名前

- コマンド名: `gkn`

## 前提

- macOS ローカル
- 1 つの設定ファイルのみ

## インストール

### Homebrew（準備後）

- `brew install --cask <tap>/gkn`
- tap は GitHub Actions 変数 `HOMEBREW_TAP_OWNER` / `HOMEBREW_TAP_NAME` に合わせる

### GitHub Release

- Release のアセットから取得

### ビルド（暫定）

- `go build -o gkn ./cmd/gkn`

## 品質チェック

- `scripts/quality.sh` (go vet + go test coverage)

## 設定

- 例: `~/.config/github-kanri/config.json`
- 1 箇所のみを参照
- `skillsRemote` は `gkn skills clone --remote` で上書き可

```json
{
  "projectsRoot": "~/Projects",
  "reposRoot": "~/Projects/repos",
  "skillsRoot": "~/Projects/skills",
  "skillsRemote": "git@github.com:org/skills.git",
  "skillTargets": [".codex/skills", ".claude/skills"],
  "syncTargets": [
    {
      "name": "skills",
      "src": "~/Projects/skills",
      "dest": [".codex/skills", ".claude/skills"],
      "include": ["**/*"],
      "exclude": [".git/**"]
    },
    {
      "name": "templates",
      "src": "~/Projects/shared/templates",
      "dest": ["./.github", "./docs"],
      "include": ["**/*"],
      "exclude": [".git/**"]
    }
  ],
  "allowCommands": ["git status", "git fetch", "git pull"],
  "denyCommands": ["rm -rf", "git reset --hard"],
  "syncMode": "copy",
  "conflictPolicy": "fail"
}
```

## 主要コマンド（案）

### repo

- `gkn repo list`
- `gkn repo status`
- `gkn repo open <pattern> --pick <n>`
- `gkn repo path <pattern> --pick <n>`
- `gkn repo recent --limit <n>`
- `gkn repo info <pattern> --pick <n>`
- `gkn repo graph <pattern> --pick <n> --limit <n>`
- `gkn repo clone <url> [--name <repo>]`
- `gkn repo exec --cmd "<command>" [--parallel <n>] [--timeout <sec>] [--require-clean]`

### skills / sync

- `gkn skills clone`
- `gkn skills sync [--target <name>]`
- `gkn skills watch [--target <name>]`
- `gkn skills diff [--target <name>]`
- `gkn skills verify [--target <name>]`
- `gkn skills status [--target <name>]`
- `gkn skills link [--target <name>]`
- `gkn skills pin --target <name> --ref <commit|tag>`
- `gkn skills clean [--target <name>]`

### config / system

- `gkn config show`
- `gkn config init`
- `gkn config validate`
- `gkn doctor`
- `gkn version`

## 共通オプション（案）

- `--force`: 衝突時上書き（既定は fail）
- `--dry-run`: 実行せずに差分/対象のみ表示
- `--only <pattern>`: 対象限定（glob）
- `--exclude <pattern>`: 対象除外（glob）
- `--json`: 機械可読出力

## 安全設計

- deny ルール最優先、違反時は即失敗
- 既定は `copy` + `fail`（上書きは `--force`）
- 曖昧一致は候補のみ出力（非対話）

## 目的

- ローカル GitHub プロジェクト管理を最短で回す
- LLM から呼び出しても安全/再現性/明確な失敗理由

## 開発メモ

- 詳細要件は `docs/requirements.md`
