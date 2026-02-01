# ガイドライン
## Don't Panic
プロダクションで動くコードはパニックを避けなければいけません。
パニックは連鎖的障害の主な原因です。
もしエラーが起きた場合、関数はエラーを返して、呼び出し元がどのようにエラーをハンドリングするか決めさせる必要があります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func foo(bar string) {
  if len(bar) == 0 {
    panic("bar must not be empty")
  }
  // ...
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("USAGE: foo <bar>")
    os.Exit(1)
  }
  foo(os.Args[1])
}
```

</td><td>

```go
func foo(bar string) error {
  if len(bar) == 0 {
    return errors.New("bar must not be empty")
  }
  // ...
  return nil
}

func main() {
  if len(os.Args) != 2 {
    fmt.Println("USAGE: foo <bar>")
    os.Exit(1)
  }
  if err := foo(os.Args[1]); err != nil {
    panic(err)
  }
}
```

</td></tr>
</tbody></table>

`panic`と`recover`はエラーハンドリングではありません。
プログラムはnil参照などの回復不可能な状況が発生したとき以外は出すべきではありません。
ただプログラムの初期化時は例外です。
プログラムが開始するときに異常が起きた場合にはpanicを起こしてもよいでしょう。

```go
var _statusTemplate = template.Must(template.New("name").Parse("_statusHTML"))
```

またテストでは、テストが失敗したことを示すためには`panic`ではなくて `t.Fatal` や `t.FailNow` を使うようにしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// func TestFoo(t *testing.T)

f, err := ioutil.TempFile("", "test")
if err != nil {
  panic("failed to set up test")
}
```

</td><td>

```go
// func TestFoo(t *testing.T)

f, err := ioutil.TempFile("", "test")
if err != nil {
  t.Fatal("failed to set up test")
}
```

</td></tr>
</tbody></table>

## Use go.uber.org/atomic

[sync/atomic](https://golang.org/pkg/sync/atomic)パッケージによるアトミック操作は`int32`や`int64`といった基本的な型を対象としているため、アトミックに操作すべき変数に対する読み出し・変更操作にアトミック操作を用いるということ(つまりsync/atomicパッケージの関数を使うこと自体)を容易に忘却させます。例では`int32`の変数に普通の読み出し操作を行ってしまっていますが、これはコンパイラの型チェック機構を素通ししてしまっているため潜在的に競合条件のあるコードをコンパイルできてしまっています。

[go.uber.org/atomic](https://godoc.org/go.uber.org/atomic)は実際のデータの型を基底型として隠蔽することによりこれらのアトミック操作に対して型安全性を付与できます。これによって読み出し操作を行う方法はアトミックな操作に限定され、普通の読み出し操作はコンパイラの型チェックの機構によってコンパイル時にはじくことが可能となります。
また`sync/atomic`パッケージに加えて便利な`atomic.Bool`型も提供しています。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type foo struct {
  running int32  // アトミックな操作が必要な変数
}

func (f* foo) start() {
  if atomic.SwapInt32(&f.running, 1) == 1 {
     // すでに実行中
     return
  }
  // Fooを開始
}

func (f *foo) isRunning() bool {
  return f.running == 1  // 競合条件! --> 別スレッドから実行されたatomic.SwapInt32による値の更新が見えないことが起こりうる
}
```

</td><td>

```go
type foo struct {
  running atomic.Bool
}

func (f *foo) start() {
  if f.running.Swap(true) {
     // すでに実行中
     return
  }
  // Fooを開始
}

func (f *foo) isRunning() bool {
  return f.running.Load()  // 読み出し操作がアトミックなため決定論的な振る舞いにになる(実はこの振る舞いはGoのメモリバリア指定がデフォルトSeqCstであることに依存するが、本項では深く触れない)
}
```

</td></tr>
</tbody></table>

## Avoid Mutable Globals
グローバル変数を変更するのは避けましょう。
代わりに依存関係の注入を使って構造体に持たせるようにしましょう。
関数ポインタを他の値と同じように構造体にもたせます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// sign.go

var _timeNow = time.Now

func sign(msg string) string {
  now := _timeNow()
  return signWithTime(msg, now)
}
```

</td><td>

```go
// sign.go

type signer struct {
  now func() time.Time
}

func newSigner() *signer {
  return &signer{
    now: time.Now,
  }
}

func (s *signer) Sign(msg string) string {
  now := s.now()
  return signWithTime(msg, now)
}
```
</td></tr>
<tr><td>

```go
// sign_test.go

func TestSign(t *testing.T) {
  oldTimeNow := _timeNow
  _timeNow = func() time.Time {
    return someFixedTime
  }
  defer func() { _timeNow = oldTimeNow }()

  assert.Equal(t, want, sign(give))
}
```

</td><td>

```go
// sign_test.go

func TestSigner(t *testing.T) {
  s := newSigner()
  s.now = func() time.Time {
    return someFixedTime
  }

  assert.Equal(t, want, s.Sign(give))
}
```

</td></tr>
</tbody></table>


## Avoid Embedding Types in Public Structs
型の埋め込みは実装の詳細を漏らし、型のブラッシュアップを阻害し、ドキュメントが曖昧になってしまいます。

共通の `AbstractList` を使って様々なリスト型を実装すると仮定します。
この場合に `AbstractList` を埋め込むことでリスト操作を実装するのはやめましょう。
代わりにメソッドを再度定義してその中で `AbstractList` のメソッドを実装するようにしましょう。

```go
type AbstractList struct {}

// Add adds an entity to the list.
func (l *AbstractList) Add(e Entity) {
  // ...
}

// Remove removes an entity from the list.
func (l *AbstractList) Remove(e Entity) {
  // ...
}
```
<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// ConcreteList is a list of entities.
type ConcreteList struct {
  *AbstractList
}
```

</td><td>

```go
// ConcreteList is a list of entities.
type ConcreteList struct {
  list *AbstractList
}

// Add adds an entity to the list.
func (l *ConcreteList) Add(e Entity) {
  return l.list.Add(e)
}

// Remove removes an entity from the list.
func (l *ConcreteList) Remove(e Entity) {
  return l.list.Remove(e)
}
```

</td></tr>
</tbody></table>

Go では継承が無い代わりに[埋め込み]( https://golang.org/doc/effective_go.html#embedding )を使えます。
外部の型は暗黙的に埋め込まれた型のメソッドを実装しています。
これらのメソッドはデフォルトでは埋め込まれた型のインスタンスのメソッドになります。

構造体は埋め込んだ型と同じ名前のフィールドを作成します。
なので埋め込んだ型が公開されていたら、そのフィールドも公開されます。
後方互換性を保つために、外側の型は埋め込んだ型を保持する必要があります。

埋め込みが必要な場面は殆どありません。
多くは面倒なメソッドの移譲を書かずに済ませるために使われます。

構造体の代わりに AbstractList インタフェースを埋め込むこともできます。
これだと開発者に将来的な自由度をもたせることができます。
しかし、抽象的な実装に依存して実装の詳細が漏れるという問題は解決されません。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// AbstractList is a generalized implementation
// for various kinds of lists of entities.
type AbstractList interface {
  Add(Entity)
  Remove(Entity)
}

// ConcreteList is a list of entities.
type ConcreteList struct {
  AbstractList
}
```

</td><td>

```go
// ConcreteList is a list of entities.
type ConcreteList struct {
  list *AbstractList
}

// Add adds an entity to the list.
func (l *ConcreteList) Add(e Entity) {
  return l.list.Add(e)
}

// Remove removes an entity from the list.
func (l *ConcreteList) Remove(e Entity) {
  return l.list.Remove(e)
}
```

</td></tr>
</tbody></table>

構造体の埋め込みでもインタフェースの埋め込みでも、将来的な型の変更に制限がかかります。

* 埋め込まれたインタフェースにメソッドを追加することは破壊的変更になります
* 埋め込まれた構造体からメソッドを削除すると破壊的変更になります
* 埋め込まれた型を削除することは破壊的変更になります
* 埋め込まれた型を同じインタフェースを実装した別の型に差し替える場合も破壊的変更になります

埋め込みの代わりに同じメソッドを書くのは面倒ですが、その分実装の詳細を外側から隠すことができます。
実装の詳細をそのメソッドが持つことで内部での変更がしやすくなります。
実装がすぐに見えるので、Listの詳細をさらに見に行く必要がなくなります。

## Avoid Using Built-in Names

Go の[言語仕様]( https://golang.org/ref/spec )ではいくつかのビルトイン、または[定義済み識別子]( https://golang.org/ref/spec#Predeclared_identifiers )があります。これらは Go のプログラム内で識別子として使うべきではありません。

状況にもよりますが、これらの識別子を再利用すると、元の識別子がレキシカルスコープ内で隠蔽されるか、元のコードを混乱させます。
コンパイラがエラーを出して気づく場合もありますが、最悪の場合は grep などでは発見困難な潜在的バグを起こす可能性があります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var error string
// `error` はビルトインの型を隠す

// or

func handleErrorMessage(error string) {
    // `error` はビルトインの型を隠す
}
```

</td><td>

```go
var errorMessage string
// `error` はビルトインの型のまま

// or

func handleErrorMessage(msg string) {
    // `error` はビルトインの型のまま
}
```

</td></tr>
<tr><td>

```go
type Foo struct {
    // これらのフィールドは
    // 技術的にはシャドウイングを引き起こしませんが、
    // `error` や `string` という文字列は曖昧です。
    error  error
    string string
}

func (f Foo) Error() error {
    // `error` と `f.error` は見た目が似ている
    return f.error
}

func (f Foo) String() string {
    // `string` と `f.string` は見た目が似ている
    return f.string
}
```

</td><td>

```go
type Foo struct {
    // `error` や `string` という文字列は明確に型名を指します
    err error
    str string
}

func (f Foo) Error() error {
    return f.err
}

func (f Foo) String() string {
    return f.str
}
```

</td></tr>
</tbody></table>

宣言済み識別子名をローカルの識別子に使ってもコンパイラはエラーを出さないことに注意してください。ですが `go vet` などのツールはこれらのシャドウイングを正しく見つけることができます。
