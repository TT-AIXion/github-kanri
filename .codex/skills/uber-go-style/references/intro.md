# 導入
コーディングスタイルは私達のコードを統治する規則です。
これらのスタイルは、gofmt がやってくれることから少しだけ発展したものです。

このガイドのゴールはUber社内でのGoのコードでやるべき、もしくはやるべからずを説明し、コードの複雑さを管理することです。
これらのルールはコードを管理しやすくし、かつエンジニアがGoの言語機能をより生産的に利用できるようにします。

このガイドは元々同僚がGoを使ってより開発しやすくするために[Prashant Varanasi]( https://github.com/prashantv )と[Simon Newton]( https://github.com/nomis52 )によって作成されました。
長年にわたって多くのフィードバックを受けて修正されています。

このドキュメントはUber社内で使われる規則を文書化したものです。
多くは以下のリソースでも見ることができるような一般的なものです。

1. Effective Go
2. The Go common mistakes guide

全てのコードは `golint` や `go vet` を通してエラーが出ない状態にするべきです。
エディタに以下の設定を導入することを推奨しています。

1. 保存するごとに `goimports` を実行する
2. `golint` と `go vet` を実行してエラーがないかチェックする

Goのエディタのサポートについては以下の資料を参考にしてください。
https://github.com/golang/go/wiki/IDEsAndTextEditorPlugins
