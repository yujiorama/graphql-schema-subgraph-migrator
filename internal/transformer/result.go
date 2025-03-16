package transformer

import (
    "fmt"
    "os"

    "github.com/vektah/gqlparser/v2/formatter"
    "github.com/vektah/gqlparser/v2/ast"
)

// Result は変換結果を表す構造体
type Result struct {
    schema *ast.SchemaDocument
}

// String はスキーマを文字列として返す
func (r *Result) String() string {
    if r.schema == nil {
        return ""
    }
    return formatter.Schema(r.schema)
}

// Save は変換後のスキーマを指定されたパスに保存する
func (r *Result) Save(path string) error {
    if r.schema == nil {
        return fmt.Errorf("schema is nil")
    }
    return os.WriteFile(path, []byte(r.String()), 0644)
}
