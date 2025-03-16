package transformer

import (
    "fmt"

    "github.com/vektah/gqlparser/v2/ast"
    "github.com/yujiorama/graphql-schema-subgraph-migrator/internal/validator"
)

type SchemaTransformer struct {
    config              Config
    validator           *validator.SubgraphValidator
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

func (t *SchemaTransformer) Transform(schema *ast.SchemaDocument) (*ast.SchemaDocument, error) {
    transformed := t.transformSchema(schema)

    // Validate transformed schema
    if errors := t.validator.Validate(transformed); len(errors) > 0 {
        return nil, fmt.Errorf("subgraph validation failed")
    }

    // Validate composition
    if errors := t.compositionValidator.Validate(transformed); len(errors) > 0 {
        return nil, fmt.Errorf("composition validation failed")
    }

    return transformed, nil
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
                            Raw:  "https://specs.apollo.dev/federation/v2.3",
                            Kind: ast.StringValue,
                        },
                    },
                    {
                        Name: "import",
                        Value: &ast.Value{
                            Kind: ast.ListValue,
                            // ChildValueList に変更
                            List: []*ast.Value{
                                {Raw: "@key", Kind: ast.StringValue},
                                {Raw: "@external", Kind: ast.StringValue},
                                {Raw: "@shareable", Kind: ast.StringValue},
                                {Raw: "@provides", Kind: ast.StringValue},
                                {Raw: "@requires", Kind: ast.StringValue},
                            },
                        },
                    },
                },
            },
        },
    }
    
    // SchemaDefinitionList に変換
    doc.Schema = ast.SchemaDefinitionList{schemaExt}

    // Transform type definitions - Types を Definitions に変更
    for _, def := range doc.Definitions {
        if def.Kind == ast.Object {
            t.transformTypeDefinition(def)
        }
    }

    return doc
}

func (t *SchemaTransformer) transformTypeDefinition(def *ast.Definition) {
    typeConfig, ok := t.config.Types[def.Name]
    if !ok {
        // Use default configuration if type-specific config is not found
        if t.config.Defaults != nil && t.config.Defaults.Key != nil {
            typeConfig = TypeConfig{
                Keys: []KeyConfig{*t.config.Defaults.Key},
            }
        }
    }

    // Add @key directives
    for _, key := range typeConfig.Keys {
        keyDir := &ast.Directive{
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
            keyDir.Arguments = append(keyDir.Arguments, &ast.Argument{
                Name: "resolvable",
                Value: &ast.Value{
                    Raw:  fmt.Sprintf("%t", *key.Resolvable),
                    Kind: ast.BooleanValue,
                },
            })
        }
        def.Directives = append(def.Directives, keyDir)
    }

    // Add @external directives to fields
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
