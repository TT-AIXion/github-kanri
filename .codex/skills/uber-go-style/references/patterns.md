# Patterns
## Test Tables
[サブテスト]( https://blog.golang.org/subtests )を利用したテーブルドリブンテストでコアのテストロジックを繰り返すときにコードの重複を避けるようにしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// func TestSplitHostPort(t *testing.T)

host, port, err := net.SplitHostPort("192.0.2.0:8000")
require.NoError(t, err)
assert.Equal(t, "192.0.2.0", host)
assert.Equal(t, "8000", port)

host, port, err = net.SplitHostPort("192.0.2.0:http")
require.NoError(t, err)
assert.Equal(t, "192.0.2.0", host)
assert.Equal(t, "http", port)

host, port, err = net.SplitHostPort(":8000")
require.NoError(t, err)
assert.Equal(t, "", host)
assert.Equal(t, "8000", port)

host, port, err = net.SplitHostPort("1:8")
require.NoError(t, err)
assert.Equal(t, "1", host)
assert.Equal(t, "8", port)
```

</td><td>

```go
// func TestSplitHostPort(t *testing.T)

tests := []struct{
  give     string
  wantHost string
  wantPort string
}{
  {
    give:     "192.0.2.0:8000",
    wantHost: "192.0.2.0",
    wantPort: "8000",
  },
  {
    give:     "192.0.2.0:http",
    wantHost: "192.0.2.0",
    wantPort: "http",
  },
  {
    give:     ":8000",
    wantHost: "",
    wantPort: "8000",
  },
  {
    give:     "1:8",
    wantHost: "1",
    wantPort: "8",
  },
}

for _, tt := range tests {
  t.Run(tt.give, func(t *testing.T) {
    host, port, err := net.SplitHostPort(tt.give)
    require.NoError(t, err)
    assert.Equal(t, tt.wantHost, host)
    assert.Equal(t, tt.wantPort, port)
  })
}
```

</td></tr>
</tbody></table>

テストテーブルを使うと、エラーメッセージへの情報の追加やテストケースの追加も簡単ですし、コードも少なくなります。

テストケースのルールはテストケースの構造体のスライス名が `tests` 、ループ内のそれぞれのテストケースの変数名が `tt` とします。
更に入力値と出力値をわかりやすくするために`give`や`want`などのプレフィックスをつけることを推奨しています。

```go
tests := []struct{
  give     string
  wantHost string
  wantPort string
}{
  // ...
}

for _, tt := range tests {
  // ...
}
```

## Avoid Unnecessary Complexity in Table Tests
テーブルテストで実施するテストの中で条件付きのアサーションや分岐ロジックがあると、可読性が低くなり保守がとても難しくなります。
テーブルテストでは `for` の中で、複雑なコードや条件分岐を入れるべきではありません。

テストが失敗してデバッグする必要があるとき、大きくて複雑なテーブルテストは可読性と保守性を大きく損ないます。

そのようなテーブルテストはいくつかのテーブルテストに分割するか、そもそも別のテスト関数に分けるのも良いでしょう。

いくつかの方針を紹介します。

- 振る舞いが小さくなるように意識する
- 条件付きのアサーションを避けて、テストの深さを最小にする
- テーブルのフィールドが全てのテストで使われているか確認する
- 全てのロジックが全てのテストケースで実行されるか確認する

ここでいう「テストの深さ」とは、「そのテストで、前のアサーションを保持する必要がある連続したアサーションの数」といえます。
循環複雑度に近いです。より薄いテストはよりアサーション間の関係が薄く、より重要な点はそれらのアサーションは条件付きになる可能性が低くなることです。

具体的に言うと、次のような状況だとテストを読むのが難しくなります。

- テーブルのフィールドによって条件分岐が複数ある。( `shouldError` や `expectCall` などのフィールドがあると注意です )
- 特定のモックの期待値のためにたくさんの `if` がある。( `shouldCallFoo` などのフィールドがあると怪しいです )
- テーブルの中に関数がある。(フィールドの中に `setupMocks func(*FooMock)` がある )

しかし、変更された入力に基づいてのみ変化する動作をテストする場合、比較可能なユニットを別々のテストに分割して比較しにくくするのではなく、すべての入力に対してどのように動作が変化するかをよりよく説明するために、同様のケースをまとめてテーブルテストとすることが望ましい場合があります。

テスト本体が短くてわかりやすければ、成功ケースと失敗ケースの分岐経路を1つにして、`shouldErr` のようなフィールドでエラーを期待することもできます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func TestComplicatedTable(t *testing.T) {
  tests := []struct {
    give          string
    want          string
    wantErr       error
    shouldCallX   bool
    shouldCallY   bool
    giveXResponse string
    giveXErr      error
    giveYResponse string
    giveYErr      error
  }{
    // ...
  }

  for _, tt := range tests {
    t.Run(tt.give, func(t *testing.T) {
      // setup mocks
      ctrl := gomock.NewController(t)
      xMock := xmock.NewMockX(ctrl)
      if tt.shouldCallX {
        xMock.EXPECT().Call().Return(
          tt.giveXResponse, tt.giveXErr,
        )
      }
      yMock := ymock.NewMockY(ctrl)
      if tt.shouldCallY {
        yMock.EXPECT().Call().Return(
          tt.giveYResponse, tt.giveYErr,
        )
      }

      got, err := DoComplexThing(tt.give, xMock, yMock)

      // verify results
      if tt.wantErr != nil {
        require.EqualError(t, err, tt.wantErr)
        return
      }
      require.NoError(t, err)
      assert.Equal(t, want, got)
    })
  }
}
```

</td><td>

```go
func TestShouldCallX(t *testing.T) {
  // setup mocks
  ctrl := gomock.NewController(t)
  xMock := xmock.NewMockX(ctrl)
  xMock.EXPECT().Call().Return("XResponse", nil)

  yMock := ymock.NewMockY(ctrl)

  got, err := DoComplexThing("inputX", xMock, yMock)

  require.NoError(t, err)
  assert.Equal(t, "want", got)
}

func TestShouldCallYAndFail(t *testing.T) {
  // setup mocks
  ctrl := gomock.NewController(t)
  xMock := xmock.NewMockX(ctrl)

  yMock := ymock.NewMockY(ctrl)
  yMock.EXPECT().Call().Return("YResponse", nil)

  _, err := DoComplexThing("inputY", xMock, yMock)
  assert.EqualError(t, err, "Y failed")
}
```

</td></tr>
</tbody></table>

このコードの複雑さは変更したり、理解したり、このテストが正しいのかを証明するのが困難になります。

厳密なガイドラインはありませんが、システムへの複数の入力がある場合は、可読性と保守性を常に考慮して、個別のテストかテーブルテストか決定しましょう。

### Parallel Tests

並列テストなどの特殊なループ(例えば、ループの一部としてゴルーチンを起動し参照を保持するもの)では、ループのスコープ内の変数が正しい値を保持しているか注意しましょう。

```go
tests := []struct{
  give string
  // ...
}{
  // ...
}

for _, tt := range tests {
  tt := tt // for t.Parallel
  t.Run(tt.give, func(t *testing.T) {
    t.Parallel()
    // ...
  })
}
```

この例では、 `t.Parallel()` を下で読んでいるので、 `tt` という変数をループの中で再度宣言する必要があります。
もしやらないと、ほとんどのテストで `tt` 変数が期待しない値になったり、テスト中に値が変更されてしまいます。

## Functional Options
Functional Option パターンは不透明なOption型を使って内部の構造体に情報を渡すパターンです。
可変長引数を受け取り、それらを順に内部のオプションに渡します。

コンストラクタや、公開されたAPIで3つ以上の多くの引数が必要な場合、このパターンを使うと良いでしょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
// package db

func Open(
  addr string,
  cache bool,
  logger *zap.Logger
) (*Connection, error) {
  // ...
}
```

</td><td>

```go
// package db

type Option interface {
  // ...
}

func WithCache(c bool) Option {
  // ...
}

func WithLogger(log *zap.Logger) Option {
  // ...
}

// Open creates a connection.
func Open(
  addr string,
  opts ...Option,
) (*Connection, error) {
  // ...
}
```

</td></tr>
<tr><td>

キャッシュやロガーはデフォルトを使う場合でも常に指定する必要があります

```go
db.Open(addr, db.DefaultCache, zap.NewNop())
db.Open(addr, db.DefaultCache, log)
db.Open(addr, false /* cache */, zap.NewNop())
db.Open(addr, false /* cache */, log)
```

</td><td>

オプションは必要なら提供されます

```go
db.Open(addr)
db.Open(addr, db.WithLogger(log))
db.Open(addr, db.WithCache(false))
db.Open(
  addr,
  db.WithCache(false),
  db.WithLogger(log),
)
```

</td></tr>
</tbody></table>

私達が紹介する方法は非公開のメソッドを持った `Option` インタフェースを使う方法です。
指定されたオプションは非公開の `options` 構造体に保持されます。

```go
type options struct {
  cache  bool
  logger *zap.Logger
}

type Option interface {
  apply(*options)
}

type cacheOption bool

func (c cacheOption) apply(opts *options) {
  opts.cache = bool(c)
}

func WithCache(c bool) Option {
  return cacheOption(c)
}

type loggerOption struct {
  Log *zap.Logger
}

func (l loggerOption) apply(opts *options) {
  opts.logger = l.Log
}

func WithLogger(log *zap.Logger) Option {
  return loggerOption{Log: log}
}

// Open creates a connection.
func Open(
  addr string,
  opts ...Option,
) (*Connection, error) {
  options := options{
    cache:  defaultCache,
    logger: zap.NewNop(),
  }

  for _, o := range opts {
    o.apply(&options)
  }

  // ...
}
```

このパターンを実装するためにクロージャを使っていますが、この方法は開発者により柔軟性をもたせ、デバッグやテストをしやすくなると考えています。
特に、テストやモックなどで比較する際にクロージャを使うことで比較しやすくなります。
更に、`options` に他のインタフェースを実装させる事もできます。
`fmt.Stringer` インタフェースを実装すると、設定を人間にわかりやすく表示させることも可能です。

更に以下の資料が参考になります。

* [Self-referential functions and the design of options]( https://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html )
* [Functional options for friendly APIs]( https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis )
