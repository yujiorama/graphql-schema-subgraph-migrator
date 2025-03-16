package validator

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

type CompositionValidator struct {
	errors []ValidationError
}

type EntityReference struct {
	TypeName   string
	FieldNames []string
}

func NewCompositionValidator() *CompositionValidator {
	return &CompositionValidator{}
}

func (v *CompositionValidator) Validate(schema *ast.SchemaDocument) []ValidationError {
	v.errors = []ValidationError{}

	// Validate entity resolvability
	entities := v.collectEntities(schema)
	v.validateEntityResolvability(schema, entities)

	// Validate field dependencies
	v.validateFieldDependencies(schema)

	// Validate key field types
	v.validateKeyFieldTypes(schema)

	return v.errors
}

func (v *CompositionValidator) collectEntities(schema *ast.SchemaDocument) []EntityReference {
	var entities []EntityReference

	for _, def := range schema.Types {
		if def.Kind != ast.Object {
			continue
		}

		for _, dir := range def.Directives {
			if dir.Name == "key" {
				fieldsArg := findArgument(dir.Arguments, "fields")
				if fieldsArg != nil {
					entities = append(entities, EntityReference{
						TypeName:   def.Name,
						FieldNames: parseKeyFields(fieldsArg.Value.Raw),
					})
				}
			}
		}
	}

	return entities
}

func (v *CompositionValidator) validateEntityResolvability(schema *ast.SchemaDocument, entities []EntityReference) {
	for _, entity := range entities {
		typeDef := findType(schema, entity.TypeName)
		if typeDef == nil {
			continue
		}

		for _, fieldName := range entity.FieldNames {
			field := findField(typeDef.Fields, fieldName)
			if field == nil {
				v.addError(ValidationError{
					Code:     "UNRESOLVABLE_KEY_FIELD",
					Message:  fmt.Sprintf("Key field %s.%s is not defined", entity.TypeName, fieldName),
					Severity: "error",
					Path:     []string{entity.TypeName, "@key", fieldName},
				})
				continue
			}

			if hasDirective(field.Directives, "external") {
				v.addError(ValidationError{
					Code:     "EXTERNAL_KEY_FIELD",
					Message:  fmt.Sprintf("Key field %s.%s cannot be @external", entity.TypeName, fieldName),
					Severity: "error",
					Path:     []string{entity.TypeName, fieldName},
				})
			}
		}
	}
}

func (v *CompositionValidator) validateFieldDependencies(schema *ast.SchemaDocument) {
	for _, def := range schema.Types {
		if def.Kind != ast.Object {
			continue
		}

		for _, field := range def.Fields {
			// Validate @provides
			if providesDir := findDirective(field.Directives, "provides"); providesDir != nil {
				fieldsArg := findArgument(providesDir.Arguments, "fields")
				if fieldsArg == nil || fieldsArg.Value.Raw == "" {
					v.addError(ValidationError{
						Code:     "INVALID_PROVIDES",
						Message:  fmt.Sprintf("@provides directive on %s.%s must specify fields", def.Name, field.Name),
						Severity: "error",
						Path:     []string{def.Name, field.Name, "@provides"},
					})
				}
			}

			// Validate @requires
			if requiresDir := findDirective(field.Directives, "requires"); requiresDir != nil {
				fieldsArg := findArgument(requiresDir.Arguments, "fields")
				if fieldsArg == nil || fieldsArg.Value.Raw == "" {
					v.addError(ValidationError{
						Code:     "INVALID_REQUIRES",
						Message:  fmt.Sprintf("@requires directive on %s.%s must specify fields", def.Name, field.Name),
						Severity: "error",
						Path:     []string{def.Name, field.Name, "@requires"},
					})
				}
			}
		}
	}
}

func (v *CompositionValidator) validateKeyFieldTypes(schema *ast.SchemaDocument) {
	validScalarTypes := map[string]bool{
		"ID":     true,
		"String": true,
		"Int":    true,
		"Float":  true,
	}

	for _, def := range schema.Types {
		if def.Kind != ast.Object {
			continue
		}

		for _, dir := range def.Directives {
			if dir.Name != "key" {
				continue
			}

			fieldsArg := findArgument(dir.Arguments, "fields")
			if fieldsArg == nil {
				continue
			}

			fields := parseKeyFields(fieldsArg.Value.Raw)
			for _, fieldName := range fields {
				field := findField(def.Fields, fieldName)
				if field == nil {
					continue
				}

				if !validScalarTypes[field.Type.Name()] {
					v.addError(ValidationError{
						Code:     "INVALID_KEY_FIELD_TYPE",
						Message:  fmt.Sprintf("Key field %s.%s must be a scalar type", def.Name, fieldName),
						Severity: "error",
						Path:     []string{def.Name, fieldName, "type"},
					})
				}
			}
		}
	}
}

func (v *CompositionValidator) addError(err ValidationError) {
	v.errors = append(v.errors, err)
}

// Utility functions
func parseKeyFields(fieldsStr string) []string {
	return strings.Fields(fieldsStr)
}


func findType(schema *ast.SchemaDocument, name string) *ast.Definition {
	for _, def := range schema.Types {
		if def.Name == name {
			return def
		}
	}
	return nil
}

func findField(fields ast.FieldList, name string) *ast.FieldDefinition {
	for _, field := range fields {
		if field.Name == name {
			return field
		}
	}
	return nil
}
