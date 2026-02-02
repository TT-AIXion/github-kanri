# github-kanri / gkn

GitHub リポ管理 CLI。
ローカル運用・非対話・安全設計。

## できること

- リポ一覧・状態・最近更新を一括で把握
- 名前あいまい一致でパス/情報/履歴を取得
- リポの open/clone/exec を安全に実行
- スキルの clone/sync/diff/verify をまとめて運用
- 設定の init/show/validate と doctor で整備

設計方針:

- 迷わせない: 曖昧一致は候補のみ
- 壊さない: deny 最優先・違反即失敗
- 速い: まとめて観測・最小出力
- 使える: `--json` で機械可読

## 動作環境

- macOS
- Git

## インストール

Homebrew:

```
brew install TT-AIXion/github-kanri/gkn
```

更新:

```
brew update
brew upgrade gkn
```

GitHub Releases:

- Releases のアセットから取得

Go でインストール:

```
go install github.com/TT-AIXion/github-kanri/cmd/gkn@<tag>
```

ソースからビルド:

```
go build -o gkn ./cmd/gkn
```

## クイックスタート

```
gkn config init
gkn config show
gkn repo list
gkn repo status
gkn repo recent --limit 10
```

## 主要コマンド

```
gkn <command> [args]
```

- `repo` リポ操作
- `skills` スキル同期
- `config` 設定
- `doctor` 環境チェック
- `version` バージョン

詳細:

```
gkn <command> --help
```

## 設定

パス:

```
~/.config/github-kanri/config.json
```

例:

```json
{
  "projectsRoot": "~/Projects",
  "reposRoot": "~/Projects/repos",
  "skillsRoot": "~/Projects/skills",
  "skillTargets": [".codex/skills", ".claude/skills"],
  "syncTargets": [
    {
      "name": "skills",
      "src": "~/Projects/skills",
      "dest": [".codex/skills", ".claude/skills"],
      "include": ["**/*"],
      "exclude": [".git/**"]
    }
  ],
  "allowCommands": [
    "git status*",
    "git log*",
    "git rev-parse*",
    "git config*",
    "git remote*",
    "git clone*",
    "git fetch*",
    "git pull*",
    "git checkout*",
    "code *"
  ],
  "denyCommands": [
    "rm -rf*",
    "git reset --hard*"
  ],
  "syncMode": "copy",
  "conflictPolicy": "fail"
}
```

ポイント:

- `denyCommands` 最優先
- `allowCommands` ワイルドカード対応
- `~` は展開

## 出力

- デフォルト: 人間向け
- `--json`: JSON

## 品質チェック

```
scripts/quality.sh
```

## ライセンス

MIT
