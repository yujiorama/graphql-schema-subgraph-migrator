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
extend schema @link(url: "https://specs.apollo.dev/federation/v2.9", import: ["@key","@external","@shareable","@provides","@requires"])
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
type Comment @key(fields: "id", resolvable: false) {
	id: ID
	post: Post
	author: User
	body: String
}
type Query {
	getUser(userId: ID): User
	_entities(representations: [_Any!]!): [_Entity]!
	_service: _Service!
}
scalar _Any
scalar federation__FieldSet
scalar link__Import
scalar federation__ContextFieldValue
scalar federation__Scope
scalar federation__Policy
type _Service {
	sdl: String!
}
union _Entity = User | Post
```

### intermediate

```bash
$ graphql-schema-subgraph-migrator -config example/config.json -schema example/intermediate.graphqls
extend schema @link(url: "https://specs.apollo.dev/federation/v2.9", import: ["@key","@external","@shareable","@provides","@requires"])
type User @key(fields: "id", resolvable: true) {
	id: ID
	name: String @external
	email: String
}
type Post @key(fields: "id title", resolvable: true) @key(fields: "id title", resolvable: true) {
	id: ID
	title: String
	body: String
	author: User @external
}
type Comment {
	post: Post
	author: User
	body: String
}
type Query {
	getUser(userId: ID): User
	_entities(representations: [_Any!]!): [_Entity]!
	_service: _Service!
}
union _Entity = Post
scalar _Any
scalar federation__FieldSet
scalar link__Import
scalar federation__ContextFieldValue
scalar federation__Scope
scalar federation__Policy
type _Service {
	sdl: String!
}
union _Entity = User | Post
```

### directive

```bash
$ graphql-schema-subgraph-migrator -config example/config.json -schema example/directive.graphqls
extend schema @link(url: "https://specs.apollo.dev/federation/v2.9", import: ["@key","@external","@shareable","@provides","@requires"])
directive @custom on FIELD_DEFINITION
scalar _Any
scalar federation__FieldSet
scalar link__Import
scalar federation__ContextFieldValue
scalar federation__Scope
scalar federation__Policy
type _Service {
	sdl: String!
}
type Query {
	_entities(representations: [_Any!]!): [_Entity]!
	_service: _Service!
}
```

## ライセンス

MIT

## 作者

<mailto:yujiorama+github@gmail.com>
