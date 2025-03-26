package transformer

import (
	"bytes"
	"os"
	"strings"

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

	// gqlparser.formatter は GraphQL Federation の仕様に未対応のためここで改変してる。
	// `extend schema { @link` ではなく `extend schema @link` となるように。
	schemaDocument := buf.String()
	schemaDocument = strings.Replace(
		schemaDocument,
		"extend schema {\n\t@link(",
		"extend schema @link(",
		1,
	)
	schemaDocument = strings.Replace(
		schemaDocument,
		") }",
		")",
		1,
	)
	return schemaDocument
}

// Save は変換後のスキーマを指定されたパスに保存する
func (r *Result) Save(path string) error {
	return os.WriteFile(path, []byte(r.String()), 0644)
}
