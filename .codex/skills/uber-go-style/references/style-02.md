# Style
## Prefix Unexported Globals with _
公開されてないtop-levelの`var`や`const`の名前には最初にアンダースコアをつけることでより内部向けてあることが明確になります。
ただ`err`で始まる変数名は例外です。
理由としてtop-levelの変数のスコープはそのパッケージ全体です。一般的な名前を使うと別のファイルで間違った値を使ってしまうことになります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// foo.go

const (
  defaultPort = 8080
  defaultUser = "user"
)

// bar.go

func Bar() {
  defaultPort := 9090
  ...
  fmt.Println("Default port", defaultPort)

  // We will not see a compile error if the first line of
  // Bar() is deleted.
}
```

</td><td>

```go
// foo.go

const (
  _defaultPort = 8080
  _defaultUser = "user"
)
```

</td></tr>
</tbody></table>

## Embedding in Structs

埋め込まれた型は構造体の定義の最初に置くべきです。
また、通常のフィールドと区別するために1行開ける必要があります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Client struct {
  version int
  http.Client
}
```

</td><td>

```go
type Client struct {
  http.Client

  version int
}
```

</td></tr>
</tbody></table>

埋め込みは適切な方法で機能を追加、拡張できる具体的なメリットがある場合に使いましょう。
ユーザーに悪影響を与えずに行う必要があります。[Avoid Embedding Types in Public Structs]( #avoid-embedding-types-in-public-structs ) も参照しましょう。

例外: Mutex は埋め込むべきではありません。公開されないフィールドとして使いましょう。[Zero-value Mutexes are Valid]( #zero-value-mutexes-are-valid ) も参照しましょう。

埋め込みは次のことをすべきではありません。:

- 見た目や、手軽さを重視して使うこと
- 埋め込まれた型が使いにくくなること
- 埋め込まれた型のゼロ値に影響すること。もし埋め込まれた型が便利なゼロ値を持っているなら、埋め込まれた後にもそれが維持されるようにしなければいけません。
- 埋め込まれた副作用として関係ないメソッドやフィールドが公開されてしまうこと
- 公開してない型を公開してしまうこと
- 埋め込まれた型のコピーに影響が出ること
- 埋め込まれた型の API や、意味上の型が変わってしまうこと
- 非標準の形式で埋め込むこと
- 埋め込まれる型の実装の詳細を公開すること
- 型の内部を操作できるようにすること
- ユーザーが意図しない方法で内部関数の挙動を変えること

簡単にまとめると、きちんと意識して埋め込みましょうということになります。
使うべきかチェックする簡単な方法は、「埋め込みたい型の公開されているメソッドやフィールドは全て埋め込まれる型に直接追加する必要があるか？」です。
答えが、「いくつかはある」あるいは「ない」の場合は埋め込みは使わず、フィールドを使いましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type A struct {
    // 悪い例:
    // A.Lock() と A.Unlock() が
    // 使えますが、メリットはありません。
    // さらに、A の内部を操作できる
    // ようになってしまいます。
    sync.Mutex
}
```

</td><td>

```go
type countingWriteCloser struct {
    // 良い例:
    // Write() は特定の目的のために
    // 外側のレイヤーに提供されます。
    // そして、内部の型の Write() メソッドに
    // 役割を移譲しています。
    io.WriteCloser

    count int
}

func (w *countingWriteCloser) Write(bs []byte) (int, error) {
    w.count += len(bs)
    return w.WriteCloser.Write(bs)
}
```

</td></tr>
<tr><td>

```go
type Book struct {
    // 悪い例:
    // ポインタなのでゼロ値の便利さを消してしまう
    io.ReadWriter

    // other fields
}

// later

var b Book
b.Read(...)  // panic: nil pointer
b.String()   // panic: nil pointer
b.Write(...) // panic: nil pointer
```

</td><td>

```go
type Book struct {
    // 良い例:
    // 便利なゼロ値を持っている
    bytes.Buffer

    // other fields
}

// later

var b Book
b.Read(...)  // ok
b.String()   // ok
b.Write(...) // ok
```

</td></tr>
<tr><td>

```go
type Client struct {
    sync.Mutex
    sync.WaitGroup
    bytes.Buffer
    url.URL
}
```

</td><td>

```go
type Client struct {
    mtx sync.Mutex
    wg  sync.WaitGroup
    buf bytes.Buffer
    url url.URL
}
```

</td></tr>
</tbody></table>

## Local Variable Declarations
変数が明示的に設定される場合、`:=` 演算子を利用しましょう。
<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var s = "foo"
```

</td><td>

```go
s := "foo"
```

</td></tr>
</tbody></table>

しかし空のスライスを宣言する場合は`var`キーワードを利用したほうがよいでしょう。[参考資料: Declearing Empty Slices]( https://github.com/golang/go/wiki/CodeReviewComments#declaring-empty-slices )

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func f(list []int) {
  filtered := []int{}
  for _, v := range list {
    if v > 10 {
      filtered = append(filtered, v)
    }
  }
}
```

</td><td>

```go
func f(list []int) {
  var filtered []int
  for _, v := range list {
    if v > 10 {
      filtered = append(filtered, v)
    }
  }
}
```

</td></tr>
</tbody></table>

## nil is a valid slice
`nil` は長さ0のスライスとして有効です。
つまり以下のものが有効です。

* 長さ0のスライスを返す代わりに`nil`を返す
    <table>
    <thead><tr><th>Bad</th><th>Good</th></tr></thead>
    <tbody>
    <tr><td>

    ```go
    if x == "" {
    return []int{}
    }
    ```

    </td><td>

    ```go
    if x == "" {
    return nil
    }
    ```

    </td></tr>
    </tbody></table>

* スライスが空かチェックするためには`nil`かチェックするのではなく`len(s) == 0`でチェックする 
  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  func isEmpty(s []string) bool {
    return s == nil
  }
  ```

  </td><td>

  ```go
  func isEmpty(s []string) bool {
    return len(s) == 0
  }
  ```

  </td></tr>
  </tbody></table>

* varで宣言しただけのゼロ値が有効
  <table>
  <thead><tr><th>Bad</th><th>Good</th></tr></thead>
  <tbody>
  <tr><td>

  ```go
  nums := []int{}
  // or, nums := make([]int)

  if add1 {
    nums = append(nums, 1)
  }

  if add2 {
    nums = append(nums, 2)
  }
  ```

  </td><td>

  ```go
  var nums []int

  if add1 {
    nums = append(nums, 1)
  }

  if add2 {
    nums = append(nums, 2)
  }
  ```

  </td></tr>
  </tbody></table>

## Reduce Scope of Variables
できる限り変数のスコープを減らしましょう。ただネストを浅くすることとバランスを考えてください。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
err := ioutil.WriteFile(name, data, 0644)
if err != nil {
 return err
}
```

</td><td>

```go
if err := ioutil.WriteFile(name, data, 0644); err != nil {
 return err
}
```

</td></tr>
</tbody></table>

もし関数の戻り値をifの外で利用する場合、あまりスコープを縮めようとしなくてもよいでしょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
if data, err := ioutil.ReadFile(name); err == nil {
  err = cfg.Decode(data)
  if err != nil {
    return err
  }

  fmt.Println(cfg)
  return nil
} else {
  return err
}
```

</td><td>

```go
data, err := ioutil.ReadFile(name)
if err != nil {
   return err
}

if err := cfg.Decode(data); err != nil {
  return err
}

fmt.Println(cfg)
return nil
```

</td></tr>
</tbody></table>
