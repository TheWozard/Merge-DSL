package reference

import (
	"fmt"
	"merge-dsl/pkg/merge"
)

const (
	TransformSchemaReference = "schema://transform_schema.yaml"
)

type Compiler struct {
	Importer  Resolver
	Validator *SchemaValidator
}

// CompileReference imports the passed reference and passes it to CompileValidated.
func (c *Compiler) CompileReference(ref string) (*merge.Traversal, error) {
	resolved, err := c.Importer.ImportMap(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve reference for compiling: %w", err)
	}
	return c.CompileValidated(resolved.Data)
}

// CompileValidated validates the passed document and sends it to merge.Compile
func (c *Compiler) CompileValidated(document map[string]interface{}) (*merge.Traversal, error) {
	err := c.Validator.IsValidByReference(document, TransformSchemaReference)
	if err != nil {
		return nil, fmt.Errorf("failed to validated document: %w", err)
	}
	return merge.Compile(document)
}
