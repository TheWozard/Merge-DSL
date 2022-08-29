package merge

import (
	"fmt"
)

const (
	TransformSchemaReference = "file://./schemas/transform_schema.yaml"

	TypeKey       = "type"
	PropertiesKey = "properties"
	ItemsKey      = "items"
	IdsKey        = "ids"
	IdKey         = "id"

	ObjectType = "object"
	ArrayType  = "array"
	LeafType   = "leaf"
)

// CompileReference imports the passed reference and passes it to Compile.
func CompileReference(reference string) (*Definition, error) {
	document, err := ImportReference[map[string]interface{}](reference)
	if err != nil {
		return nil, fmt.Errorf("failed to load document: %w", err)
	}
	return Compile(document)
}

// Compiles the passed golang structure into a ready to use Definition.
func Compile(document map[string]interface{}) (*Definition, error) {
	err := IsValidByReference(document, TransformSchemaReference)
	if err != nil {
		return nil, fmt.Errorf("failed to validated document: %w", err)
	}
	traversal, err := compile(document)
	if err != nil {
		return nil, err
	}
	return &Definition{
		traversal: traversal,
	}, nil
}

// compile uses the type defined on the document to call the right definition compile function
// used recursively to compile the entire tree.
func compile(document map[string]interface{}) (traversal, error) {
	if typ, ok := document[TypeKey].(string); ok {
		switch typ {
		case ObjectType:
			def := &objectTraversal{}
			return def, def.compile(document)
		case ArrayType:
			def := &arrayTraversal{}
			return def, def.compile(document)
		case LeafType:
			def := &leafTraversal{}
			return def, def.compile(document)
		default:
			// Should be caught by validation, just in case.
			return nil, fmt.Errorf("unknown compile type '%s'", typ)
		}
	}
	return nil, fmt.Errorf("failed to locate type in definition")
}

func (o *objectTraversal) compile(info map[string]interface{}) error {
	o.nodeTraversals = map[string]traversal{}
	if properties, ok := info[PropertiesKey].(map[string]interface{}); ok {
		for key, property := range properties {
			if typed, ok := property.(map[string]interface{}); ok {
				traversal, err := compile(typed)
				if err != nil {
					return fmt.Errorf("failed to compile node '%s': %w", key, err)
				}
				o.nodeTraversals[key] = traversal
			}
		}
	}
	return nil
}

func (a *arrayTraversal) compile(info map[string]interface{}) error {
	if items, ok := info[PropertiesKey].(map[string]interface{}); ok {
		traversal, err := compile(items)
		if err != nil {
			return fmt.Errorf("failed to compile items: %w", err)
		}
		a.defaultTraversal = traversal
	}
	a.idTraversals = map[interface{}]traversal{}
	// We unfortunately have to store these as an array as object keys
	// in json can only be provided as strings.
	if ids, ok := info[IdsKey].([]interface{}); ok {
		for _, idData := range ids {
			if typed, ok := idData.(map[string]interface{}); ok {
				if id, ok := typed[IdKey]; ok {
					traversal, err := compile(typed)
					if err != nil {
						return fmt.Errorf("failed to compile node '%v': %w", id, err)
					}
					a.idTraversals[id] = traversal
				}
			}
		}
	}
	return nil
}

func (l *leafTraversal) compile(info map[string]interface{}) error {
	return nil
}
