# Style

## Avoid overly long lines
横にスクロールしたり、たくさん首をふるような長過ぎるコードは避けましょう。

横幅は99文字を推奨しています。書く側はこれを超えると改行したほうが良いですが、絶対ではありません。
コードが超えても問題ありません。

## Be Consistent
このガイドラインの一部は客観的に評価することができます。
ですが状況や文脈に依存する主観的なものもあります。

ですが何よりも重要なことは一貫性を保つことです。

一貫性のあるコードは保守しやすく、説明しやすく、読むときのオーバーヘッドも減らせます。
更に新しい規則やバグへの修正が非常に簡単になります。

逆に、1つのコードベース内に複数の異なったりバッティングしているスタイルがあると、メンテナンスのオーバーヘッドや、不確実なコード、認知的不協和が発生します。
これらの全てが開発速度の低下、苦痛なコードレビューを誘発し、更にバグを発生させます。

このガイドラインにある項目を自分たちのコードに適用する場合、パッケージもしくは更に大きな単位で適用することを勧めます。
サブパッケージレベルで適用することは同じコードベースに複数のスタイルを当てはめることになるため、先程述べた悪いパターンに当てはまっています。

## Group Similar Declarations
Go は似たような宣言をグループにまとめることができます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
import "a"
import "b"
```

</td><td>

```go
import (
  "a"
  "b"
)
```

</td></tr>
</tbody></table>

これはパッケージ定数やパッケージ変数、型定義などにも利用できます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go

const a = 1
const b = 2



var a = 1
var b = 2



type Area float64
type Volume float64
```

</td><td>

```go
const (
  a = 1
  b = 2
)

var (
  a = 1
  b = 2
)

type (
  Area float64
  Volume float64
)
```

</td></tr>
</tbody></table>

関係が近いものだけをグループ化しましょう。
関係ないものまでグループ化するのは避けましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
type Operation int

const (
  Add Operation = iota + 1
  Subtract
  Multiply
  ENV_VAR = "MY_ENV"
)
```

</td><td>

```go
type Operation int

const (
  Add Operation = iota + 1
  Subtract
  Multiply
)

const ENV_VAR = "MY_ENV"
```

</td></tr>
</tbody></table>

グループ化を使う場所に制限は無いので、関数内でも使うことができます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func f() string {
  var red = color.New(0xff0000)
  var green = color.New(0x00ff00)
  var blue = color.New(0x0000ff)

  ...
}
```

</td><td>

```go
func f() string {
  var (
    red   = color.New(0xff0000)
    green = color.New(0x00ff00)
    blue  = color.New(0x0000ff)
  )

  ...
}
```

</td></tr>
</tbody></table>

## Import Group Ordering
import のグループは次の2つに分けるべきです。

1. 標準パッケージ
2. それ以外

goimports がデフォルトで適用してくれます。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
import (
  "fmt"
  "os"
  "go.uber.org/atomic"
  "golang.org/x/sync/errgroup"
)
```

</td><td>

```go
import (
  "fmt"
  "os"

  "go.uber.org/atomic"
  "golang.org/x/sync/errgroup"
)
```

</td></tr>
</tbody></table>

## Package Names
パッケージ名をつける場合、以下のルールに従いましょう。

* 全て小文字で大文字やアンダースコアを使わない
* ほとんどの呼び出し側が名前付きインポートをする必要がないようにする
* 短く簡潔にすること。全ての呼び出し側で識別されることを意識してください
* 複数形にしないこと。`net/urls`ではなく`net/url`です
* "common"、"util"、"shared"、"lib"などを使わないこと。これらはなんの情報もない名前です。

[Package Names]( https://blog.golang.org/package-names )や[Style guideline for Go packages]( https://rakyll.org/style-packages/ )も参考にしてください。

## Function Names
関数名にはGoコミュニティの規則である[MixedCaps]( https://golang.org/doc/effective_go.html#mixed-caps )に従います。
例外はテスト関数です。
`TestMyFunction_WhatIsBeingTested`のようにテストの目的を示すためにアンダースコアを使って分割します。

## Import Aliasing
インポートエイリアスはパッケージ名とパッケージパスの末尾が一致していない場合に利用します。

```go
import (
  "net/http"

  client "example.com/client-go"
  trace "example.com/trace/v2"
)
```

他にもインポートするパッケージの名前がバッティングした場合には使います。
それ以外の場合は使わないようにしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
import (
  "fmt"
  "os"


  nettrace "golang.net/x/trace"
)
```

</td><td>

```go
import (
  "fmt"
  "os"
  "runtime/trace"

  nettrace "golang.net/x/trace"
)
```

</td></tr>
</tbody></table>

## Function Grouping and Ordering

* 関数は呼び出される順番におおまかにソートされるべきです
* 関数はレシーバーごとにまとめられているべきです。

なので、`struct`、`const`、`var`の次にパッケージ外に公開されている関数が来るべきです。

`newXYZ()`や`NewXYZ()`は型が定義されたあと、他のメソッドの前に定義されている必要があります。

関数はレシーバーごとにまとめられているので、ユーティリティな関数はファイルの最後の方に出てくるはずです。


<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func (s *something) Cost() {
  return calcCost(s.weights)
}

type something struct{ ... }

func calcCost(n []int) int {...}

func (s *something) Stop() {...}

func newSomething() *something {
    return &something{}
}
```

</td><td>

```go
type something struct{ ... }

func newSomething() *something {
    return &something{}
}

func (s *something) Cost() {
  return calcCost(s.weights)
}

func (s *something) Stop() {...}

func calcCost(n []int) int {...}
```

</td></tr>
</tbody></table>

## Reduce Nesting
エラーや特殊ケースなどは早めにハンドリングして`return`したりループ内では`continue`や`break`してネストが浅いコードを目指しましょう。
ネストが深いコードを減らしていきましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
for _, v := range data {
  if v.F1 == 1 {
    v = process(v)
    if err := v.Call(); err == nil {
      v.Send()
    } else {
      return err
    }
  } else {
    log.Printf("Invalid v: %v", v)
  }
}
```

</td><td>

```go
for _, v := range data {
  if v.F1 != 1 {
    log.Printf("Invalid v: %v", v)
    continue
  }

  v = process(v)
  if err := v.Call(); err != nil {
    return err
  }
  v.Send()
}
```

</td></tr>
</tbody></table>

## Unnecessary Else
if-else のどちらでも変数に代入する場合、条件に一致した場合に上書きするようにしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var a int
if b {
  a = 100
} else {
  a = 10
}
```

</td><td>

```go
a := 10
if b {
  a = 100
}
```

</td></tr>
</tbody></table>

## Top-level Variable Declarations
パッケージ変数で式と同じ型なら型名を指定しないようにしましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
var _s string = F()

func F() string { return "A" }
```

</td><td>

```go
var _s = F()
// Since F already states that it returns a string, we don't need to specify
// the type again.

func F() string { return "A" }
```

</td></tr>
</tbody></table>

式の型と合わない場合は明示するようにしましょう。

```go
type myError struct{}

func (myError) Error() string { return "error" }

func F() myError { return myError{} }

var _e error = F()
// F は myError 型を返すが私達は error 型が欲しい
```
