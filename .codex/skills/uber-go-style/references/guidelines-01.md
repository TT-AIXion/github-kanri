# ガイドライン
## Pointers to Interfaces
インタフェースをポインタとして渡す必要はほぼありません。
インタフェースは値として渡すべきです。
ただインタフェースを実装している要素はポインタでも大丈夫です。

インタフェースには2種類あります。

1. 型付けされた情報へのポインタ。これは type と考えることができます。
2. データポインタ。格納されたデータがポインタならそのまま使えます。格納されたデータが値ならその値のポインタになります。

もしインタフェースのメソッドがそのインタフェースを満たした型のデータをいじりたいなら、インタフェースの裏側の型はポインタである必要があります。

## Verify Interface Compliance
コンパイル時にインタフェースが適切に実装されているかチェックしましょう。
これは以下のことを指します。

* 公開された型がAPIとして適切に要求されたインタフェースを実装しているか
* 公開されてるかどうかに関わらず、ある型の集合が同じインタフェースを実装しているか
* その他にインタフェースを実装しなければ利用できなくなるケース

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Handler struct {
  // ...
}



func (h *Handler) ServeHTTP(
  w http.ResponseWriter,
  r *http.Request,
) {
  ...
}
```

</td><td>

```go
type Handler struct {
  // ...
}

var _ http.Handler = (*Handler)(nil)

func (h *Handler) ServeHTTP(
  w http.ResponseWriter,
  r *http.Request,
) {
  // ...
}
```

</td></tr>
</tbody></table>

`var _ http.Handler = (*Handler)(nil)` という式は `*Handler` 型が `http.Handler` インタフェースを実装していなければ、コンパイルエラーになります。

代入式の右辺はゼロ値にするべきです。
ポインタ型やスライス、マップなどは `nil` ですし、構造体ならその型の空の構造体にします。

```go
type LogHandler struct {
  h   http.Handler
  log *zap.Logger
}

var _ http.Handler = LogHandler{}

func (h LogHandler) ServeHTTP(
  w http.ResponseWriter,
  r *http.Request,
) {
  // ...
}
```

## Receivers and Interfaces
レシーバーが値のメソッドはレシーバーがポインタでも呼び出すことができますが、逆はできません。

```Go
type S struct {
  data string
}

func (s S) Read() string {
  return s.data
}

func (s *S) Write(str string) {
  s.data = str
}

sVals := map[int]S{1: {"A"}}

// You can only call Read using a value
sVals[1].Read()

// This will not compile:
//  sVals[1].Write("test")

sPtrs := map[int]*S{1: {"A"}}

// You can call both Read and Write using a pointer
sPtrs[1].Read()
sPtrs[1].Write("test")
```

同じように、メソッドのレシーバーが値型でも、ポインタがインタフェースを満たしているとみなされます。

```
type F interface {
  f()
}

type S1 struct{}

func (s S1) f() {}

type S2 struct{}

func (s *S2) f() {}

s1Val := S1{}
s1Ptr := &S1{}
s2Val := S2{}
s2Ptr := &S2{}

var i F
i = s1Val
i = s1Ptr
i = s2Ptr

// The following doesn't compile, since s2Val is a value, and there is no value receiver for f.
//   i = s2Val
```

Effective Go の [Pointers vs Values]( https://golang.org/doc/effective_go.html#pointers_vs_values )を見るとよいでしょう。

## Zero-value Mutexes are Valid
`sync.Mutex` や `sync.RWMutex` はゼロ値でも有効です。ポインタで扱う必要はありません。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
mu := new(sync.Mutex)
mu.Lock()
```

</td><td>

```go
var mu sync.Mutex
mu.Lock()
```

</td></tr>
</tbody></table>

もし構造体のポインタを使う場合、mutexはポインタでないフィールドにする必要があります。
外部に公開されてない構造体なら、mutexを埋め込みで使うこともできます。


<table>
<tbody>
<tr><td>

```go
type smap struct {
  sync.Mutex // only for unexported types

  data map[string]string
}

func newSMap() *smap {
  return &smap{
    data: make(map[string]string),
  }
}

func (m *smap) Get(k string) string {
  m.Lock()
  defer m.Unlock()

  return m.data[k]
}
```

</td><td>

```go
type SMap struct {
  mu sync.Mutex

  data map[string]string
}

func NewSMap() *SMap {
  return &SMap{
    data: make(map[string]string),
  }
}

func (m *SMap) Get(k string) string {
  m.mu.Lock()
  defer m.mu.Unlock()

  return m.data[k]
}
```

</td></tr>

</tr>
<tr>
<td>インターナルな型やmutexのインタフェースを実装している必要がある場合には埋め込みを使う</td>
<td>公開されている型にはプライベートなフィールドを使う</td>
</tr>

</tbody></table>

## Copy Slices and Maps at Boundaries
スライスやマップは内部でデータへのポインタが含まれています。なのでコピーする際には注意してください。

### Receiving Slices and Maps
引数として受け取ってフィールドに保存したスライスは、他の箇所でデータが書き換わる可能性があることを覚えておいてください。

<table>
<thead><tr><th>Bad</th> <th>Good</th></tr></thead>
<tbody>
<tr>
<td>

```go
func (d *Driver) SetTrips(trips []Trip) {
  d.trips = trips
}

trips := ...
d1.SetTrips(trips)

// ここで値が変わると d1.trips[0] も変わる
trips[0] = ...
```

</td>
<td>

```go
func (d *Driver) SetTrips(trips []Trip) {
  d.trips = make([]Trip, len(trips))
  copy(d.trips, trips)
}

trips := ...
d1.SetTrips(trips)

// d1.trips に変更が及ばない
trips[0] = ...
```

</td>
</tr>

</tbody>
</table>

### Returning Slices and Maps
同じように、公開せずに内部に保持しているスライスやマップが変更されることもあります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Stats struct {
  mu sync.Mutex
  counters map[string]int
}

// Snapshot returns the current stats.
func (s *Stats) Snapshot() map[string]int {
  s.mu.Lock()
  defer s.mu.Unlock()

  return s.counters
}

// snapshot は mutex で守られない
// レースコンディションが起きる
snapshot := stats.Snapshot()
```

</td><td>

```go
type Stats struct {
  mu sync.Mutex
  counters map[string]int
}

func (s *Stats) Snapshot() map[string]int {
  s.mu.Lock()
  defer s.mu.Unlock()

  result := make(map[string]int, len(s.counters))
  for k, v := range s.counters {
    result[k] = v
  }
  return result
}

// snapshot はただのコピーなので変更しても影響はない
snapshot := stats.Snapshot()
```

</td></tr>
</tbody></table>

## Defer to Clean Up
ファイルや mutex のロックなどをクリーンアップするために defer を使おう

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
p.Lock()
if p.count < 10 {
  p.Unlock()
  return p.count
}

p.count++
newCount := p.count
p.Unlock()

return newCount

// easy to miss unlocks due to multiple returns
```

</td><td>

```go
p.Lock()
defer p.Unlock()

if p.count < 10 {
  return p.count
}

p.count++
return p.count

// more readable
```

</td></tr>
</tbody></table>

defer のオーバーヘッドは非常に小さいです。
関数の実行時間がナノ秒のオーダーである場合には避ける必要があります。
defer を使ってほんの少しの実行コストを払えば可読性がとてもあがります。
これはシンプルなメモリアクセス以上の計算が必要な大きなメソッドに特に当てはまります。

## Channel Size is One or None
channel のサイズは普段は1もしくはバッファなしのものにするべきです。
デフォルトでは channel はバッファなしでサイズが0になっています。
それより大きいサイズにする場合はよく考える必要があります。
どのようにしてサイズを決定するのか、チャネルがいっぱいになり処理がブロックされたときにどのような挙動をするかよく考える必要があります。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// Ought to be enough for anybody!
c := make(chan int, 64)
```

</td><td>

```go
// Size of one
c := make(chan int, 1) // or
// Unbuffered channel, size of zero
c := make(chan int)
```

</td></tr>
</tbody></table>
