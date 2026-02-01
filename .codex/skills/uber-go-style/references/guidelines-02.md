# ガイドライン
## Start Enums at One
Go で enum を導入するときの標準的な方法は、型を定義して `const` のグループを作り、初期値を `iota` にすることです。
変数のデフォルト値はゼロ値です。なので通常はゼロ値ではない値から enum を始めるべきでしょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Operation int

const (
  Add Operation = iota
  Subtract
  Multiply
)

// Add=0, Subtract=1, Multiply=2
```

</td><td>

```go
type Operation int

const (
  Add Operation = iota + 1
  Subtract
  Multiply
)

// Add=1, Subtract=2, Multiply=3
```

</td></tr>
</tbody></table>

ただゼロ値を使うことに意味があるケースもあります。
例えばゼロ値をデフォルトの挙動として扱いたい場合です。

```go
type LogOutput int

const (
  LogToStdout LogOutput = iota
  LogToFile
  LogToRemote
)

// LogToStdout=0, LogToFile=1, LogToRemote=2
```

## Use `"time"` to handle time
時間を正しく扱うのは非常に困難です。
時間に対する誤解には次のようなものがあります。

1. 1日は24時間である
2. 1時間は60分である
3. 1週間は7日である
4. 1年は365日である
5. [などなど]( https://infiniteundo.com/post/25326999628/falsehoods-programmers-believe-about-time )

例えば、1番について考えると、単純に24時間を足すだけでは正しくカレンダー上次の日になるとは限りません。

そのため、時間を扱う場合は常に[time]( https://pkg.go.dev/time?tab=doc )パッケージを使いましょう。
なぜならこのパッケージで前述の誤解を安全に処理することができるからです。

### Use `time.Time` for instants of time
時刻を扱うときは[time.Time]( https://pkg.go.dev/time?tab=doc#Time )型を使いましょう。
また、時刻を比較したり、足し引きする際にも[time.Time]( https://pkg.go.dev/time?tab=doc#Time )型のメソッドを使いましょう

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func isActive(now, start, stop int) bool {
  return start <= now && now < stop
}
```

</td><td>

```go
func isActive(now, start, stop time.Time) bool {
  return (start.Before(now) || start.Equal(now)) && now.Before(stop)
}
```

</td></tr>
</tbody></table>

### Use `time.Duration` for periods of time
期間を扱うときには[time.Duration]( https://pkg.go.dev/time?tab=doc#Duration )型を使いましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func poll(delay int) {
  for {
    // ...
    time.Sleep(time.Duration(delay) * time.Millisecond)
  }
}

poll(10) // was it seconds or milliseconds?
```

</td><td>

```go
func poll(delay time.Duration) {
  for {
    // ...
    time.Sleep(delay)
  }
}

poll(10*time.Second)
```

</td></tr>
</tbody></table>

時刻に24時間を足す例に戻ります。
もしカレンダー上で次の日の同じ時刻にしたい場合は [`time.AddDate`]( https://pkg.go.dev/time?tab=doc#Time.AddDate )メソッドを使います。
もしその時刻から正確に24時間後にしたい場合は[`time.Add`]( https://pkg.go.dev/time?tab=doc#Time.Add )メソッドを使います。

```go
newDay := t.AddDate(0 /* years */, 0, /* months */, 1 /* days */)
maybeNewDay := t.Add(24 * time.Hour)
```

### Use `time.Time` and `time.Duration` with external systems

できるなら外部システムとのやり取りにも `time.Time` 型や `time.Duration` 型を使うようにしましょう。

* Command-line flags: [`flag`]( https://golang.org/pkg/flag/ ) パッケージは[`time.ParseDuration`]( https://golang.org/pkg/time/#ParseDuration )を使うことで `time.Duration` 型をサポートできます
* JSON: [`encoding/json`]( https://golang.org/pkg/encoding/json/ )パッケージは[`Unmarshal` メソッド]( https://golang.org/pkg/time/#Time.UnmarshalJSON )によって[RFC 3339]( https://tools.ietf.org/html/rfc3339 )フォーマットの時刻を `time.Time` 型にエンコーディングできます
* SQL: [`database/sql`]( https://golang.org/pkg/database/sql/ )パッケージではもしドライバーがサポートしていれば `DATETIME` や `TIMESTAMP` 型のカラムを `time.Time` 型にすることができます
* YAML: [`gopkg.in/yaml.v2`]( https://godoc.org/gopkg.in/yaml.v2 )パッケージは[`time.ParseDuration`]( https://golang.org/pkg/time/#ParseDuration )によって[RFC 3339]( https://tools.ietf.org/html/rfc3339 )フォーマットの時刻を `time.Time` 型にエンコーディングできます

もし `time.Duration` が使えないなら、`int` 型や `float64` 型を使ってフィールド名に単位をもたせるようにしましょう。
次の表のようにします。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// {"interval": 2}
type Config struct {
  Interval int `json:"interval"`
}
```

</td><td>

```go
// {"intervalMillis": 2000}
type Config struct {
  IntervalMillis int `json:"intervalMillis"`
}
```

</td></tr>
</tbody></table>

もしこれらのインタラクションで `time.Time` を使えない場合には、[RFC 3339]( https://tools.ietf.org/html/rfc3339 )フォーマットの `string` 型を使うようにしましょう。
このフォーマットは[time.UnmarshalText]( https://golang.org/pkg/time/#Time.UnmarshalText )メソッドの中でも使われますし、`time.Parse` や `time.Format` 関数でも [`time.RFC3339`]( https://pkg.go.dev/time?tab=doc#RFC3339 )と組み合わせて使えます。
