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
        outputPath  = flag.String("output", "", "出力ファイルパス（指定がない場合は標準出力）")
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

    result, err := t.TransformFile(*schemaPath)
    if err != nil {
        fmt.Fprintln(os.Stderr, "エラー:", err)
        os.Exit(1)
    }

    // 出力先の処理
    if *outputPath != "" {
        if err := result.Save(*outputPath); err != nil {
            fmt.Fprintln(os.Stderr, "出力エラー:", err)
            os.Exit(1)
        }
    } else {
        fmt.Println(result.String())
    }
}
