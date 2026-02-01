# ガイドライン
## Avoid `init()`

できるなら `init()` を使うのは避けましょう。次のケースの場合は避けようがなかったり、推奨されます。

- 実行環境や呼び出し方に関係なく、決定的である場合
- 他の `init()` 関数の実行順や副作用に影響されない場合。`init()` の順序はよく知られていますが、コードが変更され、 `init()` 関数間の関係によってはコードが脆弱になり、エラーが発生しやすくなります
- グローバルな情報や、環境変数、ワーキングディレクトリ、プログラムの引数や入力にアクセスしたり、操作しない場合
- ファイルシステム、ネットワーク、システムコールの操作をしない場合

コードがこれらの要件を満たせない場合、 `main()` 関数か、プログラムのライフサイクルの一部でヘルパー関数として呼び出すか、`main()` 関数で直接呼び出す必要があります。
特にライブラリなど他のプログラムで使われることを想定したコードの場合は特に決定的であることに注意し、 "init magic" を引き起こさないように注意しましょう。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Foo struct {
    // ...
}

var _defaultFoo Foo

func init() {
    _defaultFoo = Foo{
        // ...
    }
}
```

</td><td>

```go
var _defaultFoo = Foo{
    // ...
}

// もしくはテスト可能性に配慮してこのようにします

var _defaultFoo = defaultFoo()

func defaultFoo() Foo {
    return Foo{
        // ...
    }
}
```

</td></tr>
<tr><td>

```go
type Config struct {
    // ...
}

var _config Config

func init() {
    // 悪い例: 現在のディレクトリに依存している
    cwd, _ := os.Getwd()

    // 悪い例: I/O
    raw, _ := os.ReadFile(
        path.Join(cwd, "config", "config.yaml"),
    )

    yaml.Unmarshal(raw, &_config)
}
```

</td><td>

```go
type Config struct {
    // ...
}

func loadConfig() Config {
    cwd, err := os.Getwd()
    // handle err

    raw, err := os.ReadFile(
        path.Join(cwd, "config", "config.yaml"),
    )
    // handle err

    var config Config
    yaml.Unmarshal(raw, &config)

    return config
}
```

</td></tr>
</tbody></table>

これらを考慮すると、次のような状況では `init()` が望ましかったり必要になる可能性があります。

- ただの代入では表現できない複雑な式
- `database/sql` の登録や、エンコーディングタイプの登録など、プラグイン的に使うフック
- [Google Cloud Functions]( https://cloud.google.com/functions/docs/bestpractices/tips#use_global_variables_to_reuse_objects_in_future_invocations ) などの決定的事前処理の最適化

## Exit in Main

Go のプログラムでは即時終了するために [`os.Exit`]( https://golang.org/pkg/os/#Exit ) や [`log.Fatal*`]( https://golang.org/pkg/log/#Fatal ) を使います。`panic()` を使うのは良い方法ではありません [don't panic](#dont-panic)を読んでください。
`os.Exit` や `log.Fatal*` を読んでいいのは `main()` 関数だけです。他の関数ではエラーを返して失敗を通知しましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func main() {
  body := readFile(path)
  fmt.Println(body)
}

func readFile(path string) string {
  f, err := os.Open(path)
  if err != nil {
    log.Fatal(err)
  }

  b, err := io.ReadAll(f)
  if err != nil {
    log.Fatal(err)
  }

  return string(b)
}
```

</td><td>

```go
func main() {
  body, err := readFile(path)
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println(body)
}

func readFile(path string) (string, error) {
  f, err := os.Open(path)
  if err != nil {
    return "", err
  }

  b, err := io.ReadAll(f)
  if err != nil {
    return "", err
  }

  return string(b), nil
}
```

</td></tr>
</tbody></table>

根拠: `exit` する関数が複数あるプログラムにはいくつかの問題があります。

- 明確でない制御フロー: どの関数もプログラムを強制終了できるので、制御フローを推論するのが難しくなります。
- テストが難しくなる: 強制終了するプログラムはテストで呼び出されたときも終了します。これはその関数をテストするのも難しくなりますし、`go test` でテストされるはずだった他のテストがスキップされる危険もあります。
- 後処理のスキップ: プログラムが強制終了されたとき、`defer` で終了時に実行する予定だった関数がスキップされます。これは重要な後処理をスキップする危険があります。

### Exit Once
可能なら、 `os.Exit` か `log.Fatal` を `main()` 関数で **一度だけ** 呼ぶのが好ましいです。もしプログラムを停止する失敗のシナリオがいくつかある場合、そのロジックは別の関数にしてエラーを返しましょう。

こうすると、 `main()` 関数を短くすることができますし、重要なビジネスロジックを分離してテストしやすくすることができます。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
package main

func main() {
  args := os.Args[1:]
  if len(args) != 1 {
    log.Fatal("missing file")
  }
  name := args[0]

  f, err := os.Open(name)
  if err != nil {
    log.Fatal(err)
  }
  defer f.Close()

  // もしここで log.Fatal を呼ぶと、
  // f.Close は実行されません。

  b, err := io.ReadAll(f)
  if err != nil {
    log.Fatal(err)
  }

  // ...
}
```

</td><td>

```go
package main

func main() {
  if err := run(); err != nil {
    log.Fatal(err)
  }
}

func run() error {
  args := os.Args[1:]
  if len(args) != 1 {
    return errors.New("missing file")
  }
  name := args[0]

  f, err := os.Open(name)
  if err != nil {
    return err
  }
  defer f.Close()

  b, err := io.ReadAll(f)
  if err != nil {
    return err
  }

  // ...
}
```

</td></tr>
</tbody></table>

## Use field tags in marshaled structs
JSON や YAML あるいは他のタグを使ったフィールド名をサポートするフォーマットに変換する場合、関連するタグを使ってアノテーションをつけましょう。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Stock struct {
  Price int
  Name  string
}

bytes, err := json.Marshal(Stock{
  Price: 137,
  Name:  "UBER",
})
```

</td><td>

```go
type Stock struct {
  Price int    `json:"price"`
  Name  string `json:"name"`
  // 安全に Name フィールドを name に変換できる
}

bytes, err := json.Marshal(Stock{
  Price: 137,
  Name:  "UBER",
})
```

</td></tr>
</tbody></table>

根拠: 構造体をシリアライズした形式は異なるシステムをつなぐ約束事です。
フィールド名を含む構造体のシリアライズした形式が変わってしまうと、この約束事が破れてしまいます。
タグを使ってフィールド名を指定するとこの約束事がより厳密になり、リファクタリングやフィールドのリネームで不意に壊れてしまうことを防ぐことができます。
