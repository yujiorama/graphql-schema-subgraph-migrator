package validator

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

type SubgraphValidator struct {
	errors []ValidationError
}

func NewSubgraphValidator() *SubgraphValidator {
	return &SubgraphValidator{}
}

func (v *SubgraphValidator) Validate(schema *ast.SchemaDocument) []ValidationError {
	v.errors = []ValidationError{}
	v.validateSchemaDefinition(schema)
	v.validateTypes(schema)
	v.validateDirectives(schema)
	return v.errors
}

func (v *SubgraphValidator) validateSchemaDefinition(schema *ast.SchemaDocument) {
	if len(schema.Schema) == 0 {
		v.addError(ValidationError{
			Code:     "MISSING_SCHEMA_EXTENSION",
			Message:  "Schema must be defined with 'extend schema'",
			Severity: "error",
			Path:     []string{"schema"},
		})
		return
	}

	hasLinkDirective := false
	for _, schemaDefinition := range schema.Schema {
		for _, dir := range schemaDefinition.Directives {
			if dir.Name == "link" {
				hasLinkDirective = true
				// Validate @link directive arguments
				urlArg := findArgument(dir.Arguments, "url")
				if urlArg == nil || !strings.HasPrefix(urlArg.Value.Raw, "https://specs.apollo.dev/federation/") {
					v.addError(ValidationError{
						Code:     "INVALID_LINK_URL",
						Message:  "Federation @link directive must point to specs.apollo.dev",
						Severity: "error",
						Path:     []string{"schema", "@link", "url"},
					})
				}

				importArg := findArgument(dir.Arguments, "import")
				if importArg == nil {
					v.addError(ValidationError{
						Code:     "MISSING_IMPORTS",
						Message:  "@link directive must specify imports",
						Severity: "error",
						Path:     []string{"schema", "@link", "import"},
					})
				}
			}
		}
	}

	if !hasLinkDirective {
		v.addError(ValidationError{
			Code:     "MISSING_LINK_DIRECTIVE",
			Message:  "Schema must have @link directive",
			Severity: "error",
			Path:     []string{"schema", "@link"},
		})
	}
}

func (v *SubgraphValidator) validateTypes(schema *ast.SchemaDocument) {
	for _, def := range schema.Definitions {
		if def.Kind != ast.Object {
			continue
		}

		hasKey := false
		hasExternal := false

		// Check @key directives
		for _, dir := range def.Directives {
			if dir.Name == "key" {
				hasKey = true
				fieldsArg := findArgument(dir.Arguments, "fields")
				if fieldsArg == nil || fieldsArg.Value.Raw == "" {
					v.addError(ValidationError{
						Code:     "INVALID_KEY_FIELDS",
						Message:  fmt.Sprintf("@key directive on type %s must specify fields", def.Name),
						Severity: "error",
						Path:     []string{def.Name, "@key"},
					})
				}
			}
		}

		// Check fields for @external directive
		for _, field := range def.Fields {
			for _, dir := range field.Directives {
				if dir.Name == "external" {
					hasExternal = true
					break
				}
			}
		}

		if hasExternal && !hasKey {
			v.addError(ValidationError{
				Code:     "EXTERNAL_WITHOUT_KEY",
				Message:  fmt.Sprintf("Type %s has @external fields but no @key directive", def.Name),
				Severity: "error",
				Path:     []string{def.Name},
			})
		}
	}
}

func (v *SubgraphValidator) validateDirectives(schema *ast.SchemaDocument) {
	for _, def := range schema.Definitions {
		if def.Kind != ast.Object {
			continue
		}

		for _, field := range def.Fields {
			// Validate @provides and @requires
			if hasDirective(field.Directives, "provides") && hasDirective(field.Directives, "external") {
				v.addError(ValidationError{
					Code:     "INVALID_PROVIDES_WITH_EXTERNAL",
					Message:  fmt.Sprintf("Field %s.%s cannot have both @provides and @external", def.Name, field.Name),
					Severity: "error",
					Path:     []string{def.Name, field.Name},
				})
			}
		}
	}
}

func (v *SubgraphValidator) addError(err ValidationError) {
	v.errors = append(v.errors, err)
}

// Utility functions（既に定義済みの場合は不要です）
func hasDirective(directives ast.DirectiveList, name string) bool {
	for _, dir := range directives {
		if dir.Name == name {
			return true
		}
	}
	return false
}

func findArgument(args ast.ArgumentList, name string) *ast.Argument {
	for _, arg := range args {
		if arg.Name == name {
			return arg
		}
	}
	return nil
}
