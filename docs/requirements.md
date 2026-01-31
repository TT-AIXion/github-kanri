# github-kanri 要件定義（ドラフト）

## 目的

- ローカルの GitHub リポジトリ管理を最短で回す CLI
- エージェント/人間どちらも迷わず実行できる、コマンドのみの決定的操作
- LLM から呼び出しても安全・再現性・失敗理由が明確

## スコープ

- 対象環境: macOS ローカルのみ
- 実行手段: ターミナルでのコマンド実行のみ（自然言語入力なし）
- 連携: 外部サービス連携なし（GitHub API など不要）
- ログ保存: しない（コンソール出力のみ）

## 前提/パス規約

- リポジトリ基準パス: `~/Projects/repos/`
- スキル基準パス: `~/Projects/skills/`
- スキル配布先: `.codex/skills/`, `.claude/skills/`, ほか既知のスキルディレクトリ
- 既存構成に合わせて設定で上書き可

## 用語

- Repo: `~/Projects/repos/` 配下の Git リポジトリ
- Skills: `~/Projects/skills/` 配下に管理される共通スキル群
- Sync: 共有資産（skills 等）をクローン/更新し各 Repo に配布
- Sync Targets: 同期対象のディレクトリ集合（JSON で可変）

## 主要ユースケース

1. Repo 一覧/状態確認を最短で実行
2. Skills をクローンして全 Repo に配布
3. Skills 編集後に各 Repo の `.codex/skills` 等を自動更新

## 機能要件

### 1. Repo 管理

- Repo 走査: `~/Projects/repos/` 配下の `.git` を検出
- Repo 状態: `git status --porcelain` 相当を一覧化
- Repo 一括操作: 指定コマンドの一括実行（例: fetch, pull, status）
- Repo open: repo 名の部分一致で `code <path>` 起動（複数一致は候補一覧出力→`--pick <n>` で選択）

### 2. Skills/Sync 管理

- skills クローン: 設定されたリモートから `~/Projects/skills/` にクローン
- skills 更新: `git fetch/pull` で更新
- 配布先生成: `.codex/skills/`, `.claude/skills/` 等の既知ディレクトリを作成
- 配布方式: 既定は `copy`（安全: 衝突時は fail）
- `mirror`/`link` は明示指定で有効化
- 監視モード: skills 更新検知で自動 sync（任意コマンド）
- 共有資産 sync: templates/scripts/config/docs 等を JSON で指定して配布
- sync 対象はデフォルト値 + JSON で上書き/追加/除外可能
- 衝突時の挙動: 既定 `fail`（`--force` で上書き）
- `include`/`exclude` は glob（`.gitignore` 互換を想定）

### 3. 安全制御（JSON 設定）

- 設定ファイル: JSON（例: `~/.config/github-kanri/config.json`）
- allow/deny ルール: コマンド単位・パス単位の許可/拒否
- 危険操作: deny 時は非対話で即失敗（exit != 0）

### 4. LLM から有利になる仕様

- サブコマンドは決定的で固定
- 出力はタグ付き短文（例: `OK`, `WARN`, `ERR`）で機械読み取り容易
- 非対話: すべて引数で完結、プロンプト表示なし
- 曖昧一致: 候補一覧を出して終了（選択は次コマンドで指定）

## 追加便利機能（案）

- `repo add/clone`: URL から `~/Projects/repos/` へ配置
- `repo exec --parallel --timeout`: 並列/タイムアウト
- `repo dirty/changed`: 変更ありのみ抽出
- `repo open`: repo 名部分一致で開く（複数一致は `--pick`）
- `repo path`: 対象 repo のフルパス出力
- `repo recent`: 最終更新順で一覧
- `repo info`: remote/origin, default branch, dirty を1行表示
- `repo graph`: 簡易ログ要約
- `repo exec --require-clean`: dirty repo を除外
- `skills diff`: 配布先との差分表示
- `skills link`: シンボリックリンク配布
- `skills pin`: commit/tag 固定配布（再現性）
- `skills verify`: ハッシュ一致チェック
- `skills status`: 各配布先の同期状態一覧
- `skills clean`: 設定外ディレクトリの削除（deny 連動）
- `config validate/doctor`: パス/権限/git 有無チェック
- `include/exclude`: 対象 Repo 絞り込み
- `--json` 任意: LLM 向け機械可読出力（標準出力のみ）

## コマンド名（決定）

- `gkn`

## CLI 仕様（案）

- コマンド名: `gkn`

```
gkn repo list
gkn repo status
gkn repo exec --cmd "git fetch --all"
gkn skills clone
gkn skills sync
gkn skills watch
gkn config show
```

### 共通オプション（案）

- `--force`: 衝突時上書き（既定は fail）
- `--dry-run`: 実行せずに差分/対象のみ表示
- `--only <pattern>`: 対象限定（glob）
- `--exclude <pattern>`: 対象除外（glob）

## 設定（案）

```json
{
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
    },
    {
      "name": "templates",
      "src": "~/Projects/shared/templates",
      "dest": ["./.github", "./docs"],
      "include": ["**/*"],
      "exclude": [".git/**"]
    },
    {
      "name": "scripts",
      "src": "~/Projects/shared/scripts",
      "dest": ["./scripts"],
      "include": ["**/*"],
      "exclude": [".git/**"]
    },
    {
      "name": "config",
      "src": "~/Projects/shared/config",
      "dest": ["./"],
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

## 非機能要件

- 依存最小: Go 標準 + 最小の OSS
- 実行速度: Repo 200 でも数秒で完了
- 安全: deny ルール最優先、fail fast
- 可観測性: 標準出力のみ、簡潔ログ

## 受け入れ基準（MVP）

- `gkn repo list/status` が正しく Repo 検出
- `gkn skills clone/sync` で skills 配布が完了
- 設定 JSON で deny されたコマンドが確実に止まる

## 未決事項（要確認）

- skills リモート URL の指定方法（設定のみ/引数でも）
- sync 方式: rsync/コピー/ハードリンク
- watch の実装: fsnotify で十分か
- 既存 skills の上書き方針（安全優先か完全ミラーか）
- 既知の skill ディレクトリの種類（追加候補）
