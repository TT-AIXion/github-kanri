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
gkn clone https://github.com/OWNER/REPO.git
gkn repo list
gkn repo status
gkn repo recent --limit 10
```

## 主要コマンド

```
gkn <command> [args]
```

- `clone` リポ clone（reposRoot 配下）
- `repo` リポ操作
- `skills` スキル同期
- `config` 設定
- `doctor` 環境チェック
- `version` バージョン
- `--json` JSON 出力（グローバル）

詳細:

```
gkn <command> --help
```

## コマンド一覧（完全）

```
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

補足:

- `gkn clone` は `gkn repo clone` の別名
- `reposRoot` 既定: `~/Projects/repos`
- `skills clone` は `skillsRemote` or `--remote`

## 使い方（目的→コマンド→結果）

### はじめに設定

やりたいこと: 初期設定を作る  
使う:

```
gkn config init
```

結果: `~/.config/github-kanri/config.json` 作成。既定 `reposRoot=~/Projects/repos`。

やりたいこと: 現在の設定を確認  
使う:

```
gkn config show
```

結果: 設定を表示。

やりたいこと: 設定の不整合を検出  
使う:

```
gkn config validate
```

結果: OK or エラー詳細。

### リポを増やす

やりたいこと: HTTP/SSH で clone して管理対象に入れる  
使う:

```
gkn clone <url>
gkn clone <url> --name repo
```

結果: `reposRoot` 配下に clone。既に存在なら失敗。

### 管理対象の把握

やりたいこと: 管理対象リポ一覧  
使う:

```
gkn repo list
gkn repo list --only "foo*"
gkn repo list --exclude "tmp*"
```

結果: `name path` を出力。

やりたいこと: 変更の有無だけ知りたい  
使う:

```
gkn repo status
```

結果: `clean/dirty` を表示。

やりたいこと: 最近更新された順で見たい  
使う:

```
gkn repo recent --limit 20
```

結果: 最新コミット時刻順。

### 特定リポに対する操作

やりたいこと: そのリポを開く  
使う:

```
gkn repo open <pattern> [--pick n]
```

結果: `code <path>` 実行。曖昧一致は候補のみ。

やりたいこと: パスだけ欲しい  
使う:

```
gkn repo path <pattern> [--pick n]
```

結果: パス文字列。

やりたいこと: 主要情報を見たい  
使う:

```
gkn repo info <pattern> [--pick n]
```

結果: origin/current/default/dirty を表示。

やりたいこと: 簡易ログを見たい  
使う:

```
gkn repo graph <pattern> [--pick n] [--limit n]
```

結果: oneline のログ。

### 一括コマンド実行

やりたいこと: すべてのリポで同じコマンド  
使う:

```
gkn repo exec --cmd "git status -sb"
gkn repo exec --cmd "go test ./..." --parallel 4
gkn repo exec --cmd "make lint" --require-clean
```

結果: 各リポで実行。失敗/dirty は警告 or エラー。`--dry-run` で実行せず計画だけ。

### skills を配布・同期

やりたいこと: skills リポを取得/更新  
使う:

```
gkn skills clone --remote <url>
```

結果: `skillsRoot` に clone。既に git repo なら pull。

やりたいこと: 全リポへ skills をコピー/ミラー/リンク  
使う:

```
gkn skills sync
gkn skills sync --mode mirror
gkn skills link
```

結果: `syncTargets` で定義した `src -> dest` を各リポへ反映。

やりたいこと: 差分やズレを確認  
使う:

```
gkn skills diff
gkn skills status
gkn skills verify
```

結果: 追加/削除/変更 or 一致/不一致を表示。`verify` は不一致で exit 2。

やりたいこと: 特定 ref に固定  
使う:

```
gkn skills pin --target skills --ref <commit|tag>
```

結果: `skillsRoot` を指定 ref に checkout。dirty の場合 `--force` 必須。

やりたいこと: 余計なファイルを消す  
使う:

```
gkn skills clean --force
```

結果: `src` にないファイルを `dest` から削除（`--dry-run` で確認）。

### 監視

やりたいこと: skills の変更を検知して自動同期  
使う:

```
gkn skills watch --interval 5
```

結果: `skillsRoot` の状態変化で `sync` 実行。

### 出力形式

やりたいこと: 機械処理したい  
使う:

```
gkn --json repo list
```

結果: JSON 出力。

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
  "denyCommands": ["rm -rf*", "git reset --hard*"],
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
