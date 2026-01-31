# アーキテクチャ設計書

## 背景

- ローカル運用の定型化
- 安全な一括操作

## 目的と範囲

- 目的: 再現性/安全性のある CLI
- 対象外: サーバ/リモート実行

## 品質属性（非機能要件）

- 安全性: deny 優先、非対話
- 速度: repo 200 で数秒
- 運用性: 単一設定/簡潔ログ

## 全体構成

- CLI バイナリ
- ローカル FS
- Git CLI

## コンポーネント

- Config Loader: JSON 読み込み/検証
- Repo Scanner: repo 検出/フィルタ
- Executor: コマンド実行/timeout
- Sync Engine: copy/link/mirror
- Safety Guard: allow/deny/force
- Output Formatter: tag 付き短文/JSON

## データフロー

1) 引数解析
2) 設定ロード
3) 対象 repo/targets 解決
4) 安全チェック
5) 実行/同期
6) 出力

## ストレージ

- 永続 DB なし
- 設定: JSON 1 ファイル

## セキュリティ

- allow/deny ルール
- 破壊的操作の既定拒否
- ログ永続化なし

## 運用

- ローカル単体実行
- 更新は再ビルド/再配布

## リスク

- 同期上書き事故
- コマンド競合
- パス誤設定
