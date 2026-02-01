---
name: go-cli-best-practices
description: Go製CLIの設計/実装/配布のベストプラクティス参照ガイド。UX仕様、フレームワーク選定、品質ゲート、配布、参考リポを体系化。
---

# Go CLI Best Practices

目的: Go製CLIの設計/実装/配布判断を高速化。参照リポと優先順を固定。
前提: 2026-01-31時点のスナップショット。数値/リリース日は要更新。

## 使いどころ

- Go CLIの仕様策定/実装方針決定
- フレームワーク選定、lint/CI、配布導線の設計
- 参照実装/テンプレの取捨選択

## 進め方（最小）

1) 前提整理: 内部/外部配布、対象OS、互換性、署名/配布要否
2) UX仕様固定: cli-guidelines で入出力/終了コード/ヘルプ/エラー規約を先に決める
3) フレームワーク選定: Cobra / urfave/cli / Kong
4) 品質ゲート: uber-go/guide + golangci-lint を基準化
5) 配布: GoReleaser で配布経路/署名/attestation 方針決定
6) 参照実装: cli/cli, docker/cli, terraform（ライセンス注意）
7) 補完/テンプレ: carapace, templates（必要時のみ）

## 注意（罠）

- project-layout は標準ではない。部分採用のみ
- terraform は BSL 1.1。流用注意
- stars/リリース日は変動。最新は `gh repo view` かWebで再確認

## 参照

- 詳細一覧/メモ: `references/REPOSITORIES.md`
