# ガイドライン
## Don't fire-and-forget goroutines

ゴルーチンは軽量ですが、コストはかかります。少なくとも、スタックのメモリとスケジュールされたCPUを使います。
典型的な使い方をする限りコストは小さいですが、ライフタイムを考えずに大量に作り出すと大きなパフォーマンス問題を引き起こします。
ライフタイムが管理されてないゴルーチンは、ガベージコレクションの邪魔になったり、使用されなくなったリソースを保持し続けるなどの問題も引き起こす可能性があります。

なので、絶対に本番コードでゴルーチンをリークさせないようにしましょう。 [go.uber.org/goleak]( https://pkg.go.dev/go.uber.org/goleak ) でゴルーチンを使うところでリークさせてないかテストしましょう。

一般的に、全てのゴルーチンはどちらかの方法を持っている必要があります。

- 停止する時間が予測できる
- 停止すべきゴルーチンに通知する方法がある

どちらのケースでも、処理をブロックしてゴルーチンの終了を待つコードも無ければいけません。


例:
<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
go func() {
  for {
    flush()
    time.Sleep(delay)
  }
}()
```

</td><td>

```go
var (
  stop = make(chan struct{}) // ゴルーチンに停止を伝える
  done = make(chan struct{}) // ゴルーチンが停止したことを伝える
)
go func() {
  defer close(done)

  ticker := time.NewTicker(delay)
  defer ticker.Stop()
  for {
    select {
    case <-ticker.C:
      flush()
    case <-stop:
      return
    }
  }
}()

// Elsewhere...
close(stop)  // ゴルーチンに停止シグナルを送る
<-done       // そしてゴルーチンが止まるのを待つ
```

</td></tr>
<tr><td>

ゴルーチンを止める方法は無い。プログラムが終了するまで残り続ける

</td><td>

ゴルーチンは `close(stop)` で停止できる。そして `<-done` で待つこともできる。

</td></tr>
</tbody></table>

### Wait for goroutines to exit

システムによって起動されたゴルーチンが与えられた場合、ゴルーチンの終了を待つ方法を用意する必要があります。次の2つがよく使われる方法です。

* `sync.WaitGroup` を使います。終了を待つべきゴルーチンが複数ある場合、こちらを使う

```go
var wg sync.WaitGroup
for i := 0; i < N; i++ {
  wg.Add(1)
  go func() {
    defer wg.Done()
    // ...
  }()
}

// 全ての終了を待つ
wg.Wait()
```

* 別の `chan struct{}` を作り、ゴルーチンが終了したとき `close` します。待つべきゴルーチンが1つだけのときはこちらを使いましょう

```go
done := make(chan struct{})
go func() {
  defer close(done)
  // ...
}()

// ここでゴルーチンの終了を待つ
<-done
```

### No goroutines in `init()`

`init()` 関数ではゴルーチンを起動するのを避けましょう。 [Avoid init](#avoid-init) も参照してください。

もしパッケージがバックグラウンドで動くゴルーチンを作る必要があるなら、ゴルーチンのライフタイムを管理するオブジェクトを作りそれを公開しましょう。
そのオブジェクトは `Close`、`Stop`、`Shutdown` など、バックグラウンドのゴルーチンを停止し、その終了を待つメソッドを提供しましょう。

<table>
<thead><tr><th>Bad</th><th>Good</th></tr></thead>
<tbody>
<tr><td>

```go
func init() {
  go doWork()
}

func doWork() {
  for {
    // ...
  }
}
```

</td><td>

```go
type Worker struct{ /* ... */ }

func NewWorker(...) *Worker {
  w := &Worker{
    stop: make(chan struct{}),
    done: make(chan struct{}),
    // ...
  }
  go w.doWork()
}

func (w *Worker) doWork() {
  defer close(w.done)
  for {
    // ...
    case <-w.stop:
      return
  }
}

// Shutdown はワーカーに終了を伝え、
// ワーカーが終了するのを待ちます。
func (w *Worker) Shutdown() {
  close(w.stop)
  <-w.done
}
```

</td></tr>
<tr><td>

ユーザーがこのパッケージを公開すると、ゴルーチンを無条件に作成し、止める方法はありません。

</td><td>

ユーザーがリクエストすると、ワーカーを起動します。`Shutdown` メソッドで、ワーカーを停止し使っているリソースを開放できる手段も提供しています。

[Wait for goroutines to exit]( #wait-for-gorutines-to-exit ) でも話したように、ワーカーが複数ゴルーチンを使う場合は `WaitGroup` を使いましょう。

</td></tr>
</tbody></table>
