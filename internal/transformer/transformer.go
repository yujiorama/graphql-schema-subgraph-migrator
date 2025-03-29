package transformer

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/parser"
	"github.com/yujiorama/graphql-schema-subgraph-migrator/internal/validator"
)

type SchemaTransformer struct {
	config               Config
	validator            *validator.SubgraphValidator
	compositionValidator *validator.CompositionValidator
}

func New(configPath string) (*SchemaTransformer, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	return &SchemaTransformer{
		config:               config,
		validator:            validator.NewSubgraphValidator(),
		compositionValidator: validator.NewCompositionValidator(),
	}, nil
}

// TransformFile はファイルパスからスキーマを読み込んで変換する
func (t *SchemaTransformer) TransformFile(schemaPath string) (*Result, error) {
	// ファイルからスキーマを読み込む
	source, err := os.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}
	// スキーマをパースする
	schemaDoc, err := parser.ParseSchema(&ast.Source{
		Name:  schemaPath,
		Input: string(source),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}
	// 変換を実行
	transformed, err := t.Transform(schemaDoc)
	if err != nil {
		return nil, err
	}

	return &Result{schema: transformed}, nil
}

// Transform は ast.SchemaDocument を変換する（内部利用）
func (t *SchemaTransformer) Transform(schema *ast.SchemaDocument) (*ast.SchemaDocument, error) {
	transformed := t.transformSchema(schema)

	// Validate transformed schema
	if errors := t.validator.Validate(transformed); len(errors) > 0 {
		return nil, fmt.Errorf("subgraph validation failed:\n%s", formatValidationErrors(errors))
	}

	// Validate composition
	if errors := t.compositionValidator.Validate(transformed); len(errors) > 0 {
		return nil, fmt.Errorf("composition validation failed:\n%s", formatValidationErrors(errors))
	}

	return transformed, nil
}

// ValidationError を Code ごとにグループ化し、Path を昇順に並び替えて文字列化するヘルパー関数
func formatValidationErrors(errors []validator.ValidationError) string {
	// エラーコードごとにグループ化用のマップを定義
	groupedErrors := make(map[string][]validator.ValidationError)

	// エラーをコードごとに分類
	for _, err := range errors {
		groupedErrors[err.Code] = append(groupedErrors[err.Code], err)
	}

	// 結果をフォーマットするためのスライス
	var formattedGroups []string

	// 各エラーコードごとに処理
	for code, errs := range groupedErrors {
		// Path の昇順で並び替え
		sort.Slice(errs, func(i, j int) bool {
			return strings.Join(errs[i].Path, ".") < strings.Join(errs[j].Path, ".")
		})

		// この Code グループのエラーをフォーマット
		var formattedErrors []string
		for _, err := range errs {
			formattedErrors = append(formattedErrors, fmt.Sprintf("- %s (path: %s)", err.Message, strings.Join(err.Path, ".")))
		}

		// Code とそのエラー群を結合
		formattedGroups = append(formattedGroups, fmt.Sprintf("[%s]\n%s", code, strings.Join(formattedErrors, "\n")))
	}

	// 最終結果を結合して返す
	return strings.Join(formattedGroups, "\n\n")
}

func (t *SchemaTransformer) transformSchema(doc *ast.SchemaDocument) *ast.SchemaDocument {
	// Add extend schema with @link directive
	schemaExt := &ast.SchemaDefinition{
		Directives: ast.DirectiveList{
			&ast.Directive{
				Name: "link",
				Arguments: ast.ArgumentList{
					{
						Name: "url",
						Value: &ast.Value{
							Raw:  "https://specs.apollo.dev/federation/v2.9",
							Kind: ast.StringValue,
						},
					},
					{
						Name: "import",
						Value: &ast.Value{
							Kind: ast.ListValue,
							Children: ast.ChildValueList{
								&ast.ChildValue{Value: &ast.Value{Raw: "@key", Kind: ast.StringValue}},
								&ast.ChildValue{Value: &ast.Value{Raw: "@external", Kind: ast.StringValue}},
								&ast.ChildValue{Value: &ast.Value{Raw: "@shareable", Kind: ast.StringValue}},
								&ast.ChildValue{Value: &ast.Value{Raw: "@provides", Kind: ast.StringValue}},
								&ast.ChildValue{Value: &ast.Value{Raw: "@requires", Kind: ast.StringValue}},
							},
						},
					},
				},
			},
		},
	}
	doc.SchemaExtension = ast.SchemaDefinitionList{schemaExt}

	// Transform type definitions
	for _, def := range doc.Definitions {
		if def.Kind == ast.Object {
			t.transformTypeDefinition(def)
		}
	}

	// Add required scalar types
	requiredScalarTypes := ast.DefinitionList{
		&ast.Definition{
			Kind: ast.Scalar,
			Name: "_Any",
		}, &ast.Definition{
			Kind: ast.Scalar,
			Name: "federation__FieldSet",
		}, &ast.Definition{
			Kind: ast.Scalar,
			Name: "link__Import",
		}, &ast.Definition{
			Kind: ast.Scalar,
			Name: "federation__ContextFieldValue",
		}, &ast.Definition{
			Kind: ast.Scalar,
			Name: "federation__Scope",
		}, &ast.Definition{
			Kind: ast.Scalar,
			Name: "federation__Policy",
		},
	}
	for _, requiredScalarTypeDefinition := range requiredScalarTypes {
		scalarTypeExists := false
		for _, definition := range doc.Definitions {
			if definition.Kind == ast.Scalar && definition.Name == requiredScalarTypeDefinition.Name {
				scalarTypeExists = true
				break
			}
		}
		if !scalarTypeExists {
			doc.Definitions = append(doc.Definitions, requiredScalarTypeDefinition)
		}
	}

	// Add _Service type
	serviceTypeExists := false
	for _, definition := range doc.Definitions {
		if definition.Kind == ast.Object && definition.Name == "_Service" {
			serviceTypeExists = true
			break
		}
	}
	if !serviceTypeExists {
		doc.Definitions = append(doc.Definitions, &ast.Definition{
			Kind: ast.Object,
			Name: "_Service",
			Fields: ast.FieldList{
				{Name: "sdl", Type: ast.NonNullNamedType("String", nil)},
			},
		})
	}

	// Add Query extension
	var queryTypeDefinition *ast.Definition
	for _, definition := range doc.Definitions {
		if definition.Kind == ast.Object && definition.Name == "Query" {
			queryTypeDefinition = definition
			break
		}
	}
	if queryTypeDefinition != nil {
		entitiesFieldExists := false
		serviceFieldExists := false
		for _, field := range queryTypeDefinition.Fields {
			if field.Name == "_entities" {
				entitiesFieldExists = true
			}
			if field.Name == "_service" {
				serviceFieldExists = true
			}
		}

		if !entitiesFieldExists {
			queryTypeDefinition.Fields = append(queryTypeDefinition.Fields, &ast.FieldDefinition{
				Name: "_entities",
				Type: ast.NonNullListType(ast.NamedType("_Entity", nil), nil),
				Arguments: ast.ArgumentDefinitionList{
					&ast.ArgumentDefinition{
						Name: "representations",
						Type: ast.NonNullListType(ast.NonNullNamedType("_Any", nil), nil),
					},
				},
			})
		}
		if !serviceFieldExists {
			queryTypeDefinition.Fields = append(queryTypeDefinition.Fields, &ast.FieldDefinition{
				Name: "_service",
				Type: ast.NonNullNamedType("_Service", nil),
			})
		}
	} else {
		doc.Definitions = append(doc.Definitions, &ast.Definition{
			Kind: ast.Object,
			Name: "Query",
			Fields: ast.FieldList{
				&ast.FieldDefinition{
					Name: "_entities",
					Type: ast.NonNullListType(ast.NamedType("_Entity", nil), nil),
					Arguments: ast.ArgumentDefinitionList{
						&ast.ArgumentDefinition{
							Name: "representations",
							Type: ast.NonNullListType(ast.NonNullNamedType("_Any", nil), nil),
						},
					},
				},
				&ast.FieldDefinition{
					Name: "_service",
					Type: ast.NonNullNamedType("_Service", nil),
				},
			},
		},
		)
	}

	// Add _Entity union type
	entityTypes := map[string]bool{}
	for _, typeDefinition := range doc.Definitions {
		if typeDefinition.Kind == ast.Object {
			for _, directive := range typeDefinition.Directives {
				if strings.EqualFold(directive.Name, "key") {
					entityTypes[typeDefinition.Name] = true
					for _, argument := range directive.Arguments {
						if strings.EqualFold(argument.Name, "resolvable") {
							if strings.EqualFold(argument.Value.Raw, "false") {
								entityTypes[typeDefinition.Name] = false
							}
						}
					}
				}
			}
		}
	}

	resolvableEntityTypes := []string{}
	for entityType, isResolvable := range entityTypes {
		if isResolvable {
			resolvableEntityTypes = append(resolvableEntityTypes, entityType)
		}
	}

	if len(resolvableEntityTypes) > 0 {
		doc.Definitions = append(doc.Definitions, &ast.Definition{
			Kind:  ast.Union,
			Name:  "_Entity",
			Types: resolvableEntityTypes,
		})
	}

	return doc
}

func (t *SchemaTransformer) transformTypeDefinition(def *ast.Definition) {
	// 型設定を取得
	typeConfig, ok := t.config.Types[def.Name]
	if !ok {
		// 型固有の設定が見つからない場合、デフォルト設定を使用
		if t.config.Defaults != nil && t.config.Defaults.Key != nil {
			typeConfig = TypeConfig{
				Keys: []KeyConfig{*t.config.Defaults.Key},
			}
		}
	}

	for _, key := range typeConfig.Keys {
		// TODO https://www.apollographql.com/docs/graphos/schema-design/federated-schemas/reference/directives#key
		keyFields := regexp.MustCompile(`\s+`).Split(key.Fields, -1)
		keyFieldsExistsCount := 0
		for _, keyFieldName := range keyFields {
			for _, definedField := range def.Fields {
				if definedField.Name == keyFieldName {
					keyFieldsExistsCount++
					break
				}
			}
		}

		if keyFieldsExistsCount == len(keyFields) {
			// fieldsに指定したフィールドが存在するので `@key` ディレクティブを追加
			keyDirective := &ast.Directive{
				Name: "key",
				Arguments: ast.ArgumentList{
					{
						Name: "fields",
						Value: &ast.Value{
							Raw:  key.Fields,
							Kind: ast.StringValue,
						},
					},
				},
			}
			if key.Resolvable != nil {
				keyDirective.Arguments = append(keyDirective.Arguments, &ast.Argument{
					Name: "resolvable",
					Value: &ast.Value{
						Raw:  fmt.Sprintf("%t", *key.Resolvable),
						Kind: ast.BooleanValue,
					},
				})
			}
			def.Directives = append(def.Directives, keyDirective)
		}
	}

	// @external ディレクティブをフィールドに追加
	for _, field := range def.Fields {
		for _, externalField := range typeConfig.External {
			if field.Name == externalField {
				field.Directives = append(field.Directives, &ast.Directive{
					Name: "external",
				})
			}
		}
	}
}
