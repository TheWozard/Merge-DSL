package merge

import (
	"fmt"
)

const (
	TransformSchemaReference = "file://./schemas/transform_schema.yaml"

	TypeKey = "type"

	ObjectType = "object"
	ArrayType  = "array"
	LeafType   = "leaf"
)

func CompileReference(reference string) (Definition, error) {
	data, err := ImportValidatedReference[interface{}](reference, TransformSchemaReference)
	if err != nil {
		return nil, fmt.Errorf("failed to load reference for compilation: %w", err)
	}
	return compile(data)
}

func Compile(document interface{}) (Definition, error) {
	err := IsValidByReference(document, TransformSchemaReference)
	if err != nil {
		return nil, fmt.Errorf("failed to validated document: %w", err)
	}
	return compile(document)
}

func compile(document interface{}) (Definition, error) {
	if typed, ok := document.(map[string]interface{}); ok {
		if typ, ok := typed[TypeKey].(string); ok {
			switch typ {
			case ObjectType:
				return ObjectDefinition{}, nil
			case ArrayType:
				return ArrayDefinition{}, nil
			case LeafType:
				return LeafDefinition{}, nil
			default:
				return nil, fmt.Errorf("unknown compile type '%s'", typ)
			}
		}
	}
	return nil, fmt.Errorf("failed to locate type in definition")
}
