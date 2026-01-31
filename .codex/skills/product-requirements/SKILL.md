# プロダクト要求定義書

## 背景/課題

- ローカル GitHub リポジトリ管理が手作業で散発
- エージェント実行時の迷い/暴走リスク
- 共通 assets/skills の同期が面倒

## 目的

- コマンドのみ・非対話で安全に運用
- ローカル repo/skills の一括管理を高速化
- LLM から呼んでも再現性/失敗理由が明確

## 対象ユーザー

- 主要: 開発者/AI エージェント
- 二次: 運用/PM
- ユースケース: repo 状態確認/一括実行/skills 同期

## 要求事項

- Must
  - macOS ローカル CLI
  - 単一 JSON 設定
  - allow/deny による安全制御
  - repo 走査/状態/一括実行
  - skills + 共有資産 sync
- Should
  - 競合時 fail 既定、`--force` で上書き
  - 曖昧一致は候補出力のみ
  - glob include/exclude
- Could
  - watch/schedule
  - verify/diff

## 非ゴール

- GUI/Web/リモート実行
- 自然言語入力
- 外部サービス連携（GitHub API 等）
- ログ永続化

## 成功指標

- repo/skills 操作が数秒で完了
- 破壊的操作の誤実行ゼロ
- LLM 実行で再現性 100%

## 制約/前提

- ローカルのみ
- 既定パスは `projectsRoot` 配下
- 1 設定ファイルのみ参照

## リスク

- 名前衝突
- 同期先の上書き事故
- repo 数増加で性能劣化

## 未決事項

- sync 詳細挙動（link/mirror）
- 既定除外パターンの確定
- 配布先の標準セット
