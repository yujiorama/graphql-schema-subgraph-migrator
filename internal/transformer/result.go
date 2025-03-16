package transformer

import (
    "bytes"
    "fmt"
    "os"

    "github.com/vektah/gqlparser/v2/ast"
    "github.com/vektah/gqlparser/v2/formatter"
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

    var buf bytes.Buffer
    f := formatter.NewFormatter(
        &buf,
        formatter.WithComments(),
        formatter.WithBuiltin(),
    )
    f.FormatSchemaDocument(r.schema)
    return buf.String()
}

// Save は変換後のスキーマを指定されたパスに保存する
func (r *Result) Save(path string) error {
    if r.schema == nil {
        return fmt.Errorf("schema is nil")
    }
    
    var buf bytes.Buffer
    f := formatter.NewFormatter(
        &buf,
        formatter.WithComments(),
        formatter.WithBuiltin(),
    )
    f.FormatSchemaDocument(r.schema)
    return os.WriteFile(path, buf.Bytes(), 0644)
}
