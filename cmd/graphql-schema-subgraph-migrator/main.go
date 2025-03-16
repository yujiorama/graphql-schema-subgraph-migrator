package main

import (
    "flag"
    "fmt"
    "os"

    "github.com/yujiorama/graphql-schema-subgraph-migrator/internal/transformer"
)

var (
    version = "dev"
    commit  = "none"
    date    = "unknown"
)

func main() {
    var (
        configPath  = flag.String("config", "", "設定ファイルのパス")
        schemaPath  = flag.String("schema", "", "GraphQLスキーマファイルのパス")
        showVersion = flag.Bool("version", false, "バージョン情報を表示")
    )
    flag.Parse()

    if *showVersion {
        fmt.Printf("graphql-schema-subgraph-migrator version %s (%s) built at %s\n", version, commit, date)
        return
    }

    if *schemaPath == "" {
        fmt.Fprintln(os.Stderr, "スキーマファイルのパスを指定してください")
        flag.Usage()
        os.Exit(1)
    }

    t, err := transformer.New(*configPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "エラー:", err)
        os.Exit(1)
    }

    if err := t.Transform(*schemaPath); err != nil {
        fmt.Fprintln(os.Stderr, "エラー:", err)
        os.Exit(1)
    }
}
