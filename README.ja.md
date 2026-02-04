# github-kanri / gkn（日本語）

英語版はこちら: `README.md`

GitHub リポ管理 CLI。ローカル運用・非対話・安全設計。

## 強み

- deny 優先で危険操作を即失敗
- 非対話で再現性重視
- ローカルのみ（GitHub API 不使用）
- まとめて観測/実行できる
- `--json` で機械可読
- skills 同期/検証

## クイックスタート

```sh
gkn config init
gkn config show
gkn repo list
gkn repo status
gkn repo recent --limit 10
gkn shell install --shell zsh
```

## 主要コマンド

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

## シェル連携

`gkn shell install --shell zsh` を実行すると `gkn cd <pattern>` で移動できる。

## 設定

- `~/.config/github-kanri/config.json`
- 例: `config.example.json`
- スキーマ: `docs/config.schema.json`

## 動作環境

- macOS（主）
- Linux バイナリはベストエフォート
- Git

## ドキュメント

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
