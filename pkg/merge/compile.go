package merge

import (
	"fmt"
)

const (
	TransformSchemaReference = "file://./schemas/transform_schema.yaml"
)

func CompileReference(reference string) (*Definition, error) {
	data, err := ImportValidatedReference[interface{}](reference, TransformSchemaReference)
	if err != nil {
		return nil, fmt.Errorf("failed to load reference for compilation: %w", err)
	}
	return compile(data)
}

func Compile(document interface{}) (*Definition, error) {
	err := IsValidByReference(document, TransformSchemaReference)
	if err != nil {
		return nil, fmt.Errorf("failed to validated document: %w", err)
	}
	return compile(document)
}

func compile(document interface{}) (*Definition, error) {
	return &Definition{}, nil
}

type Definition struct {
}
