package reference

import (
	"errors"
	"fmt"

	"github.com/xeipuuv/gojsonschema"
)

// IsValidByReference validates that the passed document is valid according to the json schema loaded by reference.
// Returns nil when valid.
func IsValidByReference(document interface{}, schemaReference string) error {
	schema, err := ImportSchemaReference(schemaReference)
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
func ImportSchemaReference(schemaReference string) (*gojsonschema.Schema, error) {
	if _, ok := schemaCache[schemaReference]; !ok {
		raw, err := ImportReference[interface{}](schemaReference)
		if err != nil {
			return nil, fmt.Errorf("failed to load raw schema document: %w", err)
		}
		schema, err := gojsonschema.NewSchemaLoader().Compile(gojsonschema.NewGoLoader(raw))
		if err != nil {
			return nil, fmt.Errorf("failed to compile schema document: %w", err)
		}
		schemaCache[schemaReference] = schema
	}
	return schemaCache[schemaReference], nil
}
