---
name: functional-design
description: github-kanri 機能設計書。目的/機能要件/フロー/例外/受け入れ基準。
---

# 機能設計書

## 目的と範囲

- 目的: repo/skills 管理の定型操作を CLI 化
- 対象: ローカル repo 走査/実行/同期
- 非対象: GUI/自然言語/外部連携

## ユーザーフロー

- 基本
  - `gkn config init` → `gkn repo list/status` → `gkn skills sync`
- 例外
  - deny 命中: 即失敗
  - 曖昧一致: 候補出力のみ

## 画面/インタラクション

- CLI のみ
- 対話入力なし（必要なら再実行で `--pick`）

## 機能要件

- repo
  - list/status/open/path/recent/info/graph/clone/exec
- skills/sync
  - clone/sync/watch/diff/verify/status/link/pin/clean
- config/system
  - show/init/validate/doctor/version
- 共通
  - `--force`/`--dry-run`/`--only`/`--exclude`/`--json`

## 例外・エッジケース

- 同名 repo 複数: 候補のみ出力
- conflict: 既定 fail、`--force` で上書き
- 設定不備: validate で失敗
- 未インストール git: doctor で検知

## 受け入れ基準

- 主要コマンドが要件通り動作
- deny ルールで危険操作が止まる
- sync が全 repo に適用

## 依存関係

- Git CLI
- macOS 標準ツール
- Go ランタイム（ビルド用）

## 変更履歴

- 2026-01-31 初版
