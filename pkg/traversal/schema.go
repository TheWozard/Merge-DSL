package traversal

const (
	SchemaTypeKey       = "type"
	SchemaPropertiesKey = "properties"
	SchemaItemsKey      = "items"
	SchemaIdKey         = "id"

	SchemaObjectType = "object"
	SchemaArrayType  = "array"
	SchemaEdgeType   = "edge"
)

type (
	schemaObjectTraversalPointer struct {
		data map[string]interface{}
	}
	schemaArrayTraversalPointer struct {
		data map[string]interface{}
	}
	schemaEdgeTraversalPointer struct {
		data map[string]interface{}
	}
)

// NewRawTraversal converts raw golang struct based on json data into a traversable tree
func NewSchemaTraversal(raw map[string]interface{}) Pointer[map[string]interface{}] {
	if typ, ok := raw[SchemaTypeKey]; ok {
		switch typ {
		case SchemaObjectType:
			return schemaObjectTraversalPointer{data: raw}
		case SchemaArrayType:
			return schemaArrayTraversalPointer{data: raw}
		case SchemaEdgeType:
			return schemaEdgeTraversalPointer{data: raw}
		}
	}
	return nil
}

/*
 * schemaObjectTraversalPointer
 */

func (o schemaObjectTraversalPointer) IsEdge() bool {
	return false
}

func (o schemaObjectTraversalPointer) Value() map[string]interface{} {
	return o.data
}

func (o schemaObjectTraversalPointer) GetKey(key string) Pointer[map[string]interface{}] {
	if properties, ok := o.data[SchemaPropertiesKey].(map[string]interface{}); ok {
		if data, ok := properties[key].(map[string]interface{}); ok {
			return NewSchemaTraversal(data)
		}
	}
	return nil
}

func (o schemaObjectTraversalPointer) GetItems() []Pointer[map[string]interface{}] {
	return []Pointer[map[string]interface{}]{}
}

/*
 * schemaArrayTraversalPointer
 */

func (a schemaArrayTraversalPointer) IsEdge() bool {
	return false
}

func (a schemaArrayTraversalPointer) Value() map[string]interface{} {
	return a.data
}

func (a schemaArrayTraversalPointer) GetKey(key string) Pointer[map[string]interface{}] {
	return nil
}

func (a schemaArrayTraversalPointer) GetItems() []Pointer[map[string]interface{}] {
	results := []Pointer[map[string]interface{}]{}
	if items, ok := a.data[SchemaItemsKey].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				results = append(results, NewSchemaTraversal(data))
			}
		}
	}
	return results
}

/*
 * schemaEdgeTraversalPointer
 */

func (e schemaEdgeTraversalPointer) IsEdge() bool {
	return true
}

func (e schemaEdgeTraversalPointer) Value() map[string]interface{} {
	return e.data
}

func (e schemaEdgeTraversalPointer) GetKey(key string) Pointer[map[string]interface{}] {
	return nil
}

func (e schemaEdgeTraversalPointer) GetItems() []Pointer[map[string]interface{}] {
	return []Pointer[map[string]interface{}]{}
}
