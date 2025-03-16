# graphql-schema-subgraph-migrator

GraphQLスキーマをFederation Subgraphスキーマに変換するツール

## インストール

```bash
go install github.com/yujiorama/graphql-schema-subgraph-migrator@latest
```

## 使用方法

```bash
graphql-schema-subgraph-migrator -config config.json -schema schema.graphql
```

### オプション

```
    -config: 設定ファイルのパス（オプション）
    -schema: GraphQLスキーマファイルのパス（必須）
    -version: バージョン情報を表示
```

### 設定ファイル例

```JSON

{
  "types": {
    "User": {
      "keys": [
        {
          "fields": "id",
          "resolvable": true
        },
        {
          "fields": "email",
          "resolvable": false
        }
      ],
      "external": ["email", "name"]
    },
    "Post": {
      "keys": [
        {
          "fields": "id,title",
          "resolvable": true
        }
      ],
      "external": ["author"]
    }
  },
  "defaults": {
    "key": {
      "fields": "id",
      "resolvable": true
    }
  }
}
```

## ライセンス

MIT

## 作者

yujiorama+github@gmail.com
