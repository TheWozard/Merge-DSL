package reference

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// NewSchemaValidator creates a new SchemaValidator with an empty starting cache
func NewSchemaValidator(importer Resolver) *SchemaValidator {
	return &SchemaValidator{
		Cache:    map[string]*gojsonschema.Schema{},
		Importer: importer,
	}
}

// Used to validate documents against a schema
type SchemaValidator struct {
	Cache    map[string]*gojsonschema.Schema
	Importer Resolver
}

// IsValidByReference validates that the passed document is valid according to the json schema loaded by reference.
// Returns nil when valid.
func (sv *SchemaValidator) IsValidByReference(document interface{}, schemaReference string) error {
	schema, err := sv.ImportSchemaReference(schemaReference)
	if err != nil {
		return fmt.Errorf("failed to load schema: %w", err)
	}
	result, err := schema.Validate(gojsonschema.NewGoLoader(document))
	if err != nil {
		return fmt.Errorf("failed to validate document: %w", err)
	}
	if !result.Valid() {
		err := errors.New("failed to validate document")
		for _, schemaError := range result.Errors() {
			err = fmt.Errorf("%w: %s", err, schemaError.Description())
		}
		return err
	}
	return nil
}

// ImportSchemaReference load a reference into memory as a *gojsonschema.Schema
func (sv *SchemaValidator) ImportSchemaReference(schemaReference string) (*gojsonschema.Schema, error) {
	if _, ok := sv.Cache[schemaReference]; !ok {
		raw, err := sv.Importer.ImportInterface(schemaReference)
		if err != nil {
			return nil, fmt.Errorf("failed to load raw schema document: %w", err)
		}
		schema, err := gojsonschema.NewSchemaLoader().Compile(gojsonschema.NewGoLoader(raw.Data))
		if err != nil {
			return nil, fmt.Errorf("failed to compile schema document: %w", err)
		}
		sv.Cache[schemaReference] = schema
	}
	return sv.Cache[schemaReference], nil
}
