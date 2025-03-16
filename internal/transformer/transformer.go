package transformer

import (
	"fmt"

	"github.com/graphql-go/graphql"
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

func (t *SchemaTransformer) Transform(schema *graphql.Schema) (*graphql.Schema, error) {
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

func (t *SchemaTransformer) transformSchema(schema *graphql.Schema) *graphql.Schema {
	// Transform implementation here
	return schema
}
