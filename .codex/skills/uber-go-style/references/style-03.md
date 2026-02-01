# Style
## Avoid Naked Parameters
値をそのまま関数の引数に入れることは可読性を損ないます。
もし分かりづらいならC言語スタイルのコメントで読みやすくしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// func printInfo(name string, isLocal, done bool)

printInfo("foo", true, true)
```

</td><td>

```go
// func printInfo(name string, isLocal, done bool)

printInfo("foo", true /* isLocal */, true /* done */)
```

</td></tr>
</tbody></table>

よりよいのはただの`bool`を自作の型で置き換えることです。こうすると型安全ですし可読性も上がります。
更に将来的にtrue/false以外の状態も利用可能に修正することもできます。

```go
type Region int

const (
  UnknownRegion Region = iota
  Local
)

type Status int

const (
  StatusReady = iota + 1
  StatusDone
  // Maybe we will have a StatusInProgress in the future.
)

func printInfo(name string, region Region, status Status)
```

## Use Raw String Literals to Avoid Escaping
Goは複数行や引用符のために[`Raw string literal`]( https://golang.org/ref/spec#raw_string_lit )をサポートしています。
これらをうまく使って手動でエスケープした読みづらい文字列を避けてください。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
wantError := "unknown name:\"test\""
```

</td><td>

```go
wantError := `unknown error:"test"`
```

</td></tr>
</tbody></table>

## Initializing Structs

### Use Field Names to Initialize Structs

構造体を初期化する際にはフィールド名を書くようにしましょう。
[`go vet`]( https://golang.org/cmd/vet/ )でこのルールは指摘されます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
k := User{"John", "Doe", true}
```

</td><td>

```go
k := User{
    FirstName: "John",
    LastName: "Doe",
    Admin: true,
}
```

</td></tr>
</tbody></table>

例外としてフィールド数が3以下のテストケースなら省略してもよいです。

```go
tests := []struct{
  op Operation
  want string
}{
  {Add, "add"},
  {Subtract, "subtract"},
}
```

### Omit Zero Value Fields in Structs
フィールド名を使って構造体を初期化するときは、意味のあるコンテキストを提供しない場合はフィールド名を省略しましょう。Goが自動的に型に応じたゼロ値を設定してくれます

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
user := User{
  FirstName: "John",
  LastName: "Doe",
  MiddleName: "",
  Admin: false,
}
```

</td><td>

```go
user := User{
  FirstName: "John",
  LastName: "Doe",
}
```

</td></tr>
</tbody></table>

省略されたフィールドはデフォルトのゼロ値持つことは読み手の負荷を下げることができます。デフォルト値ではない値だけが指定されているからです。

ゼロ値をあえてセットする意味がある場合もあります。例えば、[Test Tables]( #test-tables ) で指定するテストケースではゼロ値でも設定することは役に立ちます。

### Use `var` for Zero Value Structs

全てのフィールドを省略して宣言するときは、 `var` を使って構造体の宣言をしましょう。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
user := User{}
```

</td><td>

```go
var user User
```

</td></tr>
</tbody></table>

[map initialization]( #initializing-maps )でも似たようなことをしていますが、この方法を使うと、ゼロ値だけの構造体の宣言と、フィールドに値を指定する宣言を区別することができます。そしてこの方法は[declare empty slices]( #declaring-empty-slices )と同じ理由でおすすめの方法です。

### Initializing Struct References
構造体の初期化と同じように構造体のポインタを初期化するときは`new(T)`ではなく、`&T{}`を使いましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
sval := T{Name: "foo"}

// inconsistent
sptr := new(T)
sptr.Name = "bar"
```

</td><td>

```go
sval := T{Name: "foo"}

sptr := &T{Name: "bar"}
```

</td></tr>
</tbody></table>

## Initializing Maps
空のマップを作る場合は `make(...)` を使い、コード内で実際にデータを入れます。
こうすることで、変数宣言と視覚的に区別でき、後でサイズヒントをつけやすくなります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var (
  // m1 is safe to read and write;
  // m2 will panic on writes.
  m1 = map[T1]T2{}
  m2 map[T1]T2
)
```

</td><td>

```go
var (
  // m1 is safe to read and write;
  // m2 will panic on writes.
  m1 = make(map[T1]T2)
  m2 map[T1]T2
)
```

</td></tr>
<tr><td>

変数宣言と初期化が視覚的に似ている

</td><td>

変数宣言と初期化が視覚的に区別しやすい

</td></tr>
</tbody></table>

可能なら `make()` でマップを初期化する際にキャパシティのヒントを渡しましょう。
詳細は[Prefer Specifying Map Capacity Hints]( #prefer-specifying-map-capacity-hints )を参照してください。

一方で、マップが予め決まった要素だけを保つ場合にはリテラルを使って初期化するほうがよいでしょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
m := make(map[T1]T2, 3)
m[k1] = v1
m[k2] = v2
m[k3] = v3
```

</td><td>

```go
m := map[T1]T2{
  k1: v1,
  k2: v2,
  k3: v3,
}
```

</td></tr>
</tbody></table>

大まかな原則は初期化時に決まった要素を追加するならマップリテラルを使い、それ以外なら `make()` (とあるならキャパシティのヒント)を使いましょう。

## Format Strings outside Printf
フォーマット用の文字列を`Printf`スタイルの外で定義する場合は`const`を使いましょう。
こうすることで`go vet`などの静的解析ツールでチェックしやすくなります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
msg := "unexpected values %v, %v\n"
fmt.Printf(msg, 1, 2)
```

</td><td>

```go
const msg = "unexpected values %v, %v\n"
fmt.Printf(msg, 1, 2)
```

</td></tr>
</tbody></table>

## Naming Printf-style Functions
`Printf`スタイルの関数を使う場合、`go vet`がフォーマットをチェックできるか確認しましょう。

これは可能であれば`Printf`スタイルの関数名を使う必要があることを示しています。
`go vet`はデフォルトでこれらの関数をチェックします。

事前に定義された関数名を使わない場合、関数名の最後を`f`にしましょう。
例えば`Wrap`ではなく`Wrapf`にします。
`go vet`は特定の`Printf`スタイルのチェックができるようになっていますが、末尾が`f`である必要があります。

```shell
$ go vet -printfuncs=wrapf,statusf
```

[go vet: Printf family check]( https://kuzminva.wordpress.com/2017/11/07/go-vet-printf-family-check/ )を更に参照してください。
