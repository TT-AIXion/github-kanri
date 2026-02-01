# Performance
パフォーマンスガイドラインは特によく実行される箇所にのみ適用されます。

## Prefer strconv over fmt
数字と文字列を単に変換する場合、`fmt` パッケージよりも `strconv` パッケージのほうが高速に実行されます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
for i := 0; i < b.N; i++ {
  s := fmt.Sprint(rand.Int())
}
```

</td><td>

```go
for i := 0; i < b.N; i++ {
  s := strconv.Itoa(rand.Int())
}
```

</td></tr>
<tr><td>

```
BenchmarkFmtSprint-4    143 ns/op    2 allocs/op
```

</td><td>

```
BenchmarkStrconv-4    64.2 ns/op    1 allocs/op
```

</td></tr>
</tbody></table>

## Avoid string-to-byte conversion
固定の文字列からバイト列を何度も生成するのは避けましょう。
代わりに変数に格納してそれを使うようにしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
for i := 0; i < b.N; i++ {
  w.Write([]byte("Hello world"))
}
```

</td><td>

```go
data := []byte("Hello world")
for i := 0; i < b.N; i++ {
  w.Write(data)
}
```

</tr>
<tr><td>

```
BenchmarkBad-4   50000000   22.2 ns/op
```

</td><td>

```
BenchmarkGood-4  500000000   3.25 ns/op
```

</td></tr>
</tbody></table>

## Prefer Specifying Map Capacity Hints
スライスやmapの容量のヒントが事前にある場合は初期化時に次のように設定しましょう。

```go
make(map[T1]T2, hint)
```

`make()` の引数にキャパシティを渡すと、初期化時に適切なサイズにしようとします。
なので要素をマップに追加する際にアロケーションの回数を減らすことができます。
ただし、キャパシティのヒントは必ずしも保証されるものではありません。
もし事前にキャパシティを渡していても、要素の追加時にアロケーションが発生する場合もあります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
m := make(map[string]os.FileInfo)
files, _ := ioutil.ReadDir("./files")
for _, f := range files {
    m[f.Name()] = f
}
```

</td><td>

```go
files, _ := ioutil.ReadDir("./files")
m := make(map[string]os.FileInfo, len(files))
for _, f := range files {
    m[f.Name()] = f
}
```

</td></tr>
<tr><td>

`m` は初期化時にサイズのヒントが与えられませんでした。
そのため余計にアロケーションが発生する可能性があります。

</td><td>

`m` にはサイズのヒントが与えられています。
要素の追加時にアロケーションの回数を押さえられます。

</td></tr>
</tbody></table>
