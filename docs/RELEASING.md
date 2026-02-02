# Releasing

## 自動リリース

- `main` へ push/merge → `auto-tag` が `v0.0.0-main.<UTC>.<sha>` を作成
- 同じ `auto-tag` 内で GoReleaser 実行 → GitHub Release 作成
- `-` を含む tag は pre-release 扱い

## 安定版リリース

- 手動で `vX.Y.Z` を tag して push
- `release` が GitHub Release を作成

## Homebrew 有効化

Homebrew Cask を有効にする場合のみ設定。

- `HOMEBREW_TAP_OWNER` (Actions variables)
- `HOMEBREW_TAP_NAME` (Actions variables)
- `HOMEBREW_TAP_TOKEN` (Actions secrets: tap repo へ write 可能な PAT)

`HOMEBREW_TAP_*` が未設定なら Cask の upload は skip。
