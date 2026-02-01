# Linting
どんなおすすめの linter のセットよりも重要なのは、コードベース全体で一貫した linter を使うことです。

私達は最小限の linter として以下のものをおすすめしています。これだけあれば、一般的な問題を発見でき、不必要に厳しすぎることもなく、良い品質を確立することができるからです。

- [errcheck]( https://github.com/kisielk/errcheck ) はエラーが正しく処理されているかを担保します
- [goimports]( https://godoc.org/golang.org/x/tools/cmd/goimports ) はライブラリのインポートの管理もするフォーマッタです
- [golint]( https://github.com/golang/lint ) は一般的なスタイルのミスを見つけます
- [govet]( https://golang.org/cmd/vet/ ) は一般的なミスを見つけます
- [staticcheck]( https://staticcheck.io/ ) は多くの静的チェックを行います

## Lint Runners

私達は、大きなコードベースでのパフォーマンスと、多くの linter を一度に設定できる点から、[golangci-lint]( https://github.com/golangci/golangci-lint ) をおすすめしています。
このガイドのリポジトリにはおすすめの [.golangci.yaml]( https://github.com/uber-go/guide/blob/master/.golangci.yml ) 設定ファイルがあります。

golangci-lint は[多くのlinter]( https://golangci-lint.run/usage/linters/ )が使えます。前述した linter は基本のセットですが、チームが必要に応じて linter を追加することを推奨しています。
