# graphql-schema-subgraph-migrator

GraphQLスキーマをFederation Subgraphスキーマに変換するツール

## インストール

```bash
go install github.com/yujiorama/graphql-schema-subgraph-migrator/cmd/graphql-schema-subgraph-migrator@latest
```

## 使用方法

```bash
$ graphql-schema-subgraph-migrator
スキーマファイルのパスを指定してください
Usage of graphql-schema-subgraph-migrator
  -config string
        設定ファイルのパス
  -output string
        出力ファイルパス（指定がない場合は標準出力）
  -schema string
        GraphQLスキーマファイルのパス
  -version
        バージョン情報を表示
exit status 1
```

### version

```bash
$ graphql-schema-subgraph-migrator -version
graphql-schema-subgraph-migrator version dev (none) built at unknown
```

### example

```bash
$ graphql-schema-subgraph-migrator -config example/config.json -schema example/example.graphqls
extend schema {
        @link(url: "https://specs.apollo.dev/federation/v2.3", import: ["@key","@external","@shareable","@provides","@requires"]) }
type User @key(fields: "id", resolvable: true) {
        id: ID
        name: String @external
        email: String
}
type Post @key(fields: "id title", resolvable: true) {
        id: ID
        title: String
        body: String
        author: User @external
}

```

## ライセンス

MIT

## 作者

<mailto:yujiorama+github@gmail.com>
