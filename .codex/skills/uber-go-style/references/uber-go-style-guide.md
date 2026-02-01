# Uber Go Style Guide ベストプラクティス（要約）

出典:
- https://github.com/uber-go/guide
- https://github.com/uber-go/guide/tree/master/src

## 全体構成

- ガイドは「Guidelines」「Performance」「Style」「Patterns」の大区分で整理される。
- インターフェース、リソース管理、並行性、パフォーマンス、スタイル、命名、変数スコープ、テストなどを横断的に扱う。

## インターフェース設計

- インターフェースは値で渡す。ポインタのインターフェースは使わない。
- 値レシーバーは値/ポインタ両方でインターフェース実装とみなされる。ポインタレシーバーはポインタのみ。
- map の要素はアドレス不可。レシーバーの選択はアドレス取得可否に影響する。
- インターフェース適合はコンパイル時に明示確認する（`var _ Interface = (*T)(nil)` など）。
- 公開構造体に型を埋め込まない。互換性を壊さないよう明示的に委譲する。

## エラー処理

- エラー型の選択は「マッチが必要か」「静的/動的か」で決める。
- `fmt.Errorf` でラップする場合、マッチさせたいなら `%w`、マッチさせたくないなら `%v`。
- エラーメッセージは簡潔に。`failed to` のような冗長語は避ける。
- エラー名はエクスポートは `Err`、非エクスポートは `err` をプレフィックスにする。型名は `Error` で終える。
- エラーは一度だけ処理する（処理する/ラップして返す/ログして終了 のいずれか）。二重ログを避ける。

## 並行性/リソース管理

- “Fire-and-forget” を避ける。起動した goroutine の終了を待つ。
- goroutine を `init` で起動しない。
- goroutine は `sync.WaitGroup` で終了待ちし、必要なら `context` でキャンセルする。
- チャネルバッファは 1 か 0 を基本。大きくする場合は根拠を示す。
- 原子的操作は `sync/atomic` ではなく `go.uber.org/atomic` を使う。
- `Mutex` のゼロ値は有効。ポインタ化しない。
- リソース解放は `defer` を使う。
- 境界で slice/map をコピーし、外部からのミューテーションを遮断する。

## プログラム構造

- `init` 関数は避ける。必要性がない限り使わない。
- `main` でのみ `os.Exit` を呼ぶ。`run()` でエラーを返し、`main` が1回だけ終了処理する。
- ミュータブルなグローバルを避ける。依存は注入する。
- ビルトイン名（`error` など）を変数名で上書きしない。
- 変数のスコープは最小化する（必要なら読みやすさ優先）。

## パフォーマンス

- 変換は `fmt` より `strconv` を優先する。
- `string` と `[]byte` の繰り返し変換を避ける。
- コンテナの容量は事前に指定する。
- ゼロ値構造体は `var` 宣言を使う（意味と意図が明確）。
- 最適化は計測の後に行う。

## 命名/スタイル

- パッケージ名は小文字・短い・単語1つ。`util` のような汎用名は避ける。
- 関数名は MixedCaps（Test は `TestSomething` など）。
- Printf 系は `f` サフィックス、または既存慣例名を使う。`go vet -printfuncs` で検証する。
- import のエイリアスは必要時のみ。パッケージ名と一致するなら避ける。
- 行は長くしない（目安 99 文字）。
- 近い宣言はまとめ、import はグループ順を一貫させる。

## テスト

- テーブル駆動テストを基本にする。
- テストケースは `tests := []struct{...}` で定義し、`tt` でループする。
- subtest は `t.Run(tt.name, func(t *testing.T){...})` を使う。
- `t.Parallel()` を使う場合、`tt := tt` でループ変数を再束縛する。
- テーブル内のフィールドは全ケースで使う。混在させない。
- エラー期待は `wantErr` などを使い、期待値を明示する。

## 元資料（src 主要ファイル）

- `src/SUMMARY.md`
- `src/interface-compliance.md`
- `src/interface-ptr.md`
- `src/interface-receiver.md`
- `src/handle-errors-once.md`
- `src/error-types.md`
- `src/error-wrapping.md`
- `src/error-name.md`
- `src/goroutine-forget.md`
- `src/goroutine-exit.md`
- `src/goroutine-init.md`
- `src/channel-size.md`
- `src/atomic.md`
- `src/avoid-init.md`
- `src/exit-once.md`
- `src/mutable-global.md`
- `src/dont-name-variables-after-builtins.md`
- `src/string-format.md`
- `src/strconv.md`
- `src/string-byte-slice.md`
- `src/container-capacity.md`
- `src/package-name.md`
- `src/function-name.md`
- `src/import-alias.md`
- `src/guideline-name.md`
- `src/table-driven-tests.md`
- `src/test-table.md`
- `src/test-table-parallel.md`
