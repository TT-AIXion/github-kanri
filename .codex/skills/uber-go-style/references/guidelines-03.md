# ガイドライン
## Errors

### Error Types
エラーを定義する方法にはいくつかの種類があります。
ユースケースに合った最適なものを選ぶために以下のことを考慮しましょう。

* 呼び出し側は自身でエラーをハンドリングするためにエラーを検知する必要がありますか？
  その場合、上位のエラー変数または自前の型を定義することで [`errors.Is`] と [`errors.As`] 関数を利用できるようにサポートしなければなりません。
* エラーメッセージは静的な文字列ですか？またはコンテキストを持つ情報が必要な動的な文字列ですか？
  前者ならば [`errors.New`] が利用できます。後者ならば [`fmt.Errorf`] または自前のエラー型を利用しなければなりません。
* 下流のエラーを更に上流に返していますか？もしそうならば[Error Wrappingのセクション](#error-wrapping)を参照してください。

[`errors.Is`]: https://golang.org/pkg/errors/#Is
[`errors.As`]: https://golang.org/pkg/errors/#As

[`errors.New`]: https://golang.org/pkg/errors/#New
[`fmt.Errorf`]: https://golang.org/pkg/fmt/#Errorf

| エラーを検知? | エラーメッセージ | アドバイス                          |
|-----------------|---------------|-----------------------------------|
| いいえ           | 静的        | [`errors.New`]                      |
| いいえ           | 動的       | [`fmt.Errorf`]                      |
| はい             | 静的        | [`errors.New`]を使ったパッケージ変数(`var`で定義) |
| はい             | 動的       |  自前のエラー型                 |

例えば、静的文字列のエラーならば [`errors.New`] を利用しましょう。
呼び出し側がエラーを検知しハンドリングする必要がある場合は、そのエラーをパッケージ変数とし`errors.Is`で検知できるようにしましょう。

<table>
<thead><tr><th>No error matching</th><th>Error matching</th></tr></thead>
<tbody>
<tr><td>

```go
// package foo

func Open() error {
  return errors.New("could not open")
}

// package bar

if err := foo.Open(); err != nil {
  // Can't handle the error.
  panic("unknown error")
}
```

</td><td>

```go
// package foo

var ErrCouldNotOpen = errors.New("could not open")

func Open() error {
  return ErrCouldNotOpen
}

// package bar

if err := foo.Open(); err != nil {
  if errors.Is(err, foo.ErrCouldNotOpen) {
    // handle the error
  } else {
    panic("unknown error")
  }
}
```

</td></tr>
</tbody></table>

動的文字列のエラーの場合、呼び出し側がエラー検知する必要がないならば [`fmt.Errorf`] を使い、検知する必要があるならば自前の`error`インターフェースを実装する型を使いましょう。

<table>
<thead><tr><th>No error matching</th><th>Error matching</th></tr></thead>
<tbody>
<tr><td>

```go
// package foo

func Open(file string) error {
  return fmt.Errorf("file %q not found", file)
}

// package bar

if err := foo.Open("testfile.txt"); err != nil {
  // Can't handle the error.
  panic("unknown error")
}
```

</td><td>

```go
// package foo

type NotFoundError struct {
  File string
}

func (e *NotFoundError) Error() string {
  return fmt.Sprintf("file %q not found", e.File)
}

func Open(file string) error {
  return &NotFoundError{File: file}
}


// package bar

if err := foo.Open("testfile.txt"); err != nil {
  var notFound *NotFoundError
  if errors.As(err, &notFound) {
    // handle the error
  } else {
    panic("unknown error")
  }
}
```

</td></tr>
</tbody></table>

自前のエラー型を公開する場合、それもパッケージの公開APIの一部になることに留意しましょう。

```go
// package foo

type errNotFound struct {
  file string
}

func (e errNotFound) Error() string {
  return fmt.Sprintf("file %q not found", e.file)
}

func IsNotFoundError(err error) bool {
  _, ok := err.(errNotFound)
  return ok
}

func Open(file string) error {
  return errNotFound{file: file}
}

// package bar

if err := foo.Open("foo"); err != nil {
  if foo.IsNotFoundError(err) {
    // handle
  } else {
    panic("unknown error")
  }
}
```

### Error Wrapping
エラーを伝搬させるためには以下の3つの方法が主流です。

* 受けたエラーをそのまま返す。
* `fmt.Errorf` に `%w` を付けてコンテキストを追加する。
* `fmt.Errorf` に `%v` を付けてコンテキストを追加する。

追加するコンテキストが無いならばエラーをそのまま返しましょう。これによりオリジナルのエラー型とメッセージが保たれます。これは下流のエラーメッセージにエラーがどこから来たか追うための十分な情報がある場合に適しています。

別の方法として、"connection refused"のような曖昧なエラーではなく、"call service foo: connection refused"のようなより有益なエラーを得られるように、可能な限りエラーメッセージにコンテキストを追加することもできます。

エラーにコンテキストを追加するには`fmt.Errorf`を使いましょう。このとき、呼び出し側がエラー元の原因を抽出し検知できるようにするべきかどうかに基づき`%w`または`%v`を選ぶことになります。

* 呼び出し側が原因のエラー元を把握する必要がある場合は`%w`を使いましょう。これはほとんどのラップされたエラーにとって良いデフォルトの振る舞いになりますが、呼び出し側がそれに依存しだすかもしないことを考慮しましょう。そのため、ラップされたエラーが既知の変数(var)か型(type)であるケースでは、関数の責務としてその振る舞いのコードドキュメント記載とテストをしましょう。
* 原因のエラー元をあえて曖昧にする場合`%v`を使いましょう。呼び出し側はエラー検知をすることができなくなりますが、将来必要なときに`%w`を使うように変更できます。

返されたエラーにコンテキストを追加する場合、"failed to"のようなエラーがスタックに蓄積されるあたって明白な表現は避け、コンテキストを簡潔に保つようにしてください。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
s, err := store.New()
if err != nil {
    return fmt.Errorf(
        "failed to create new store: %w", err)
}
```

</td><td>

```go
s, err := store.New()
if err != nil {
    return fmt.Errorf(
        "new store: %w", err)
}
```

</td></tr><tr><td>

```
failed to x: failed to y: failed to create new store: the error
```

</td><td>

```
x: y: new store: the error
```

</td></tr>
</tbody></table>

しかし、エラーメッセージが他のシステムに送られる場合は"err"タグを付けたり"Failed"プレフィックスをつけたりすることでエラーメッセージであることを明確にする必要があります。

[Don't just check errors, handle them gracefully]の記事も参照してください。

  [Don't just check errors, handle them gracefully]: https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully

### Error Naming

グローバル変数として使われるエラー値においてそれがパブリックかプライベートかによって`Err`または`err`のプレフィックスを付けましょう。
この助言は[Prefix Unexported Globals with _](#prefix-unexported-globals-with-_)の助言よりも優先されます。

```go
var (
  // 以下の2つのエラーはパブリックであるため
  // このパッケージの利用者はこれらのエラーを
  // errors.Isで検知することができる。

  ErrBrokenLink = errors.New("link is broken")
  ErrCouldNotOpen = errors.New("could not open")

  // このエラーはパッケージのパブリックAPIに
  // させたくないのでプライベートにしている。
  // errors.Isでパッケージ内にてこのエラーを
  // 使うことができる。

  errNotFound = errors.New("not found")
)
```

カスタムエラー型の場合、`Error`を末尾に付けるようにしましょう。

```go
// 同様にこのエラーはパブリックであるため
// このパッケージの利用者はこれらのエラーを
// errors.Asで検知することができる。

type NotFoundError struct {
  File string
}

func (e *NotFoundError) Error() string {
  return fmt.Sprintf("file %q not found", e.File)
}

// このエラーはパッケージのパブリックAPIに
// させたくないのでプライベートにしている。
// errors.Asでパッケージ内にてこのエラーを
// 使うことができる。

type resolveError struct {
  Path string
}

func (e *resolveError) Error() string {
  return fmt.Sprintf("resolve %q", e.Path)
}
```

### Handle Errors Once
呼び出し側が呼び出し先からエラーを受け取ったとき、エラーについてどれほど知っているかで様々な方法があります。

主にこれらですが、他にも色々あります。

- もし、呼び出し先が特定のエラーを定義しているなら、[`errors.Is`]や[`errors.As`]を使って処理を分岐させます。
- もし復帰可能なエラーなら、エラーをログに出力し、安全に処理を継続しましょう
- もしドメインで定義された条件を満たさないエラーなら、事前に定義したエラーを返しましょう
- エラーを返すときは、[Wrap]( #error-wrapping ) するか、一つ一つ忠実に返しましょう。

呼び出し元がどのようにエラーを扱うかに関わらず、通常はどのエラーも一度だけ処理するべきです。例えば、呼び出し元の更に呼び出し元も同様にエラーを処理するので、呼び出し元はエラーをログに記録してから返すべきではありません。

次の具体例を考えます。

<table>
<thead><tr><th>Description</th><th>Code</th></tr></thead>
<tbody>
<tr><td>

**悪い例**: エラーをログに書き出してから返す。

より上位の呼び出し先も同様にエラーをログに出力する可能性があるので、アプリケーションログに多くのノイズが混ざりこのログの価値が薄まります。

</td><td>

```go
u, err := getUser(id)
if err != nil {
  // BAD: See description
  log.Printf("Could not get user %q: %v", id, err)
  return err
}
```

</td></tr>
<tr><td>

**良い例**: エラーをラップして返す。

より上位の呼び出し側がエラーを処理します。 `%w` を使うと、[`errors.Is`]や[`errors.As`]を使ってエラーマッチさせることができます。

</td><td>

```go
u, err := getUser(id)
if err != nil {
  return fmt.Errorf("get user %q: %w", id, err)
}
```

</td></tr>
<tr><td>

**良い例**: エラーをログに出力し、処理を継続する。

もしその処理が絶対必要でないなら、品質は下がりますが復旧して処理を続けることができます。
</td><td>

```go
if err := emitMetrics(); err != nil {
  // メトリクスの書き出し失敗は処理を止めるほどではない
  log.Printf("Could not emit metrics: %v", err)
}

```

</td></tr>
<tr><td>

**良い例**: エラーマッチさせて処理を継続する。

もし呼び出し元がエラーを定義していて、そのエラーが復旧可能なら、エラーを確認して処理を継続しましょう。
その他のエラーだった場合はエラーをラップして返しましょう。

ラップされたエラーは上位の呼び出し元で処理させましょう。
</td><td>

```go
tz, err := getUserTimeZone(id)
if err != nil {
  if errors.Is(err, ErrUserNotFound) {
    // もしユーザーがみつからないなら UTC を使う
    tz = time.UTC
  } else {
    return fmt.Errorf("get user %q: %w", id, err)
  }
}
```

</td></tr>
</tbody></table>

## Handle Type Assertion Failures
[型アサーション]( https://golang.org/ref/spec#Type_assertions )で1つの戻り値を受け取る場合、その型でなかったらパニックを起こします。
型アサーションではその型に変換できたかを示すbool値も同時に返ってくるので、それで事前にチェックしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
t := i.(string)
```

</td><td>

```go
t, ok := i.(string)
if !ok {
  // 安全にエラーを処理する
}
```

</td></tr>
</tbody></table>
