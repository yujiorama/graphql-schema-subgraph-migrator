# Migration

## example

```bash
$ graphql-schema-subgraph-migrator -config ./config.json -schema ./example.graphqls 
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
type _Service {
	sdl: String!
}
union _Entity = User | Post
```

## intermediate

```bash
$ graphql-schema-subgraph-migrator -config ./config.json -schema ./intermediate.graphqls
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
union _Entity = User | Post
scalar _Any
type _Service {
	sdl: String!
}
```

## directive

```bash
$ graphql-schema-subgraph-migrator -config ./config.json -schema ./directive.graphqls
extend schema @link(url: "https://specs.apollo.dev/federation/v2.9", import: ["@key","@external","@shareable","@provides","@requires"])
directive @custom on FIELD_DEFINITION
scalar _Any
type _Service {
	sdl: String!
}
type Query {
	_entities(representations: [_Any!]!): [_Entity]!
	_service: _Service!
}
```
