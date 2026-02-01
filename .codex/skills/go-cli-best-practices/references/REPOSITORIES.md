# Go製CLIベストプラクティス参照リポ（2026-01-31時点スナップショット）

注意: 数値/リリース日は変動。最新は `gh repo view` かWebで確認。

## 0) コピペ用: リポ一覧

```text
cli-guidelines/cli-guidelines
uber-go/guide
spf13/cobra
urfave/cli
alecthomas/kong
cli/cli
docker/cli
hashicorp/terraform
goreleaser/goreleaser
golangci/golangci-lint
golang-standards/project-layout
carapace-sh/carapace
FalcoSuessgott/golang-cli-template
thazelart/golang-cli-template
```

## 1) CLIのUX/設計（言語非依存）

- cli-guidelines/cli-guidelines
  - 得られる: CLI設計原則（引数/フラグ、出力、エラー、設定、対話/非対話、スクリプト適性）
  - 取り入れる: 終了コード、stdout/stderr分離、機械可読出力、`--help`/`--version`一貫性

## 2) Goコードの書き方（流儀固定）

- uber-go/guide
  - 得られる: Goスタイルガイド。レビュー基準の共通化
  - 取り入れる: エラー処理、命名、interface、ゼロ値/初期化、並行処理

## 3) CLIフレームワーク（実装基盤）

- spf13/cobra
  - 得られる: サブコマンド前提のCLI構築、補完/ヘルプ生成
  - 向く: サブコマンド拡張前提、周辺エコシステム重視
- urfave/cli
  - 得られる: シンプルなCLI構築
  - 向く: 依存の思想が合う/シンプルさ重視
- alecthomas/kong
  - 得られる: 構造体タグで宣言的にCLI定義
  - 向く: 宣言的に仕様定義したい場合

## 4) 実運用の大規模CLI（現実解）

- cli/cli
  - 得られる: 大規模Go CLIの設計/テスト/配布の実例
  - 注目: Build Provenance Attestation、`--help`/`--json`設計、認証/失敗時出力
- docker/cli
  - 得られる: ビルド/テスト/lint/クロスビルド導線の実例
- hashicorp/terraform
  - 得られる: 超大規模CLIの実装/互換性/プラグインの現実
  - 注意: Business Source License 1.1

## 5) リリース/配布・品質

- goreleaser/goreleaser
  - 得られる: クロスビルド/配布/パッケージ自動化
- golangci/golangci-lint
  - 得られる: 複数linter統合。CIに載せやすい

## 6) テンプレート（全部入り）

- FalcoSuessgott/golang-cli-template
  - 得られる: Cobra + GoReleaser + golangci-lint + CI + manpages + completions
  - 注意: 更新が古い。構成要素は取捨選択
- thazelart/golang-cli-template
  - 得られる: Cobra + GoReleaser + テスト（testify）
  - 注意: 小規模。最新追随用途は弱い

## 7) プロジェクト構成（罠）

- golang-standards/project-layout
  - 注意: "標準"ではない。IssueでGoコア開発者が否定的
  - 取り入れ方: `cmd/` や `internal/` の考え方だけ部分採用

## 8) 補完

- carapace-sh/carapace
  - 得られる: Cobra向け補完（複数シェル対応）
  - 向く: サブコマンド引数まで補完したい

## 9) 取り入れる順序（優先度）

1. cli-guidelines でCLI仕様固定
2. フレームワーク選定（Cobra/urfave/cli/Kong）
3. 品質ゲート（golangci-lint + uber-go/guide）
4. 配布（GoReleaser）
5. 参照実装（cli/cli, docker/cli）
6. 補完/テンプレ（carapace, templates）

## 10) 評価メモ

- 前提: 内部ツール vs 外部配布で要件が分岐
- リスク: project-layout 標準化の罠 / 供給網(配布物)の真正性 / ライセンス
- 代替: 小規模なら標準 `flag` + 自前ディスパッチで開始→必要時に移行
- 優先: UX仕様 → フレームワーク → lint/CI → 配布

## 11) 実行コマンド（最小）

```bash
# ガイド/スタイル/実戦CLI
# (GitHub CLI使用)
gh repo clone cli-guidelines/cli-guidelines -- --depth=1
gh repo clone uber-go/guide -- --depth=1
gh repo clone cli/cli -- --depth=1
gh repo clone docker/cli -- --depth=1

# GoReleaser / golangci-lint
gh repo clone goreleaser/goreleaser -- --depth=1
gh repo clone golangci/golangci-lint -- --depth=1
```
