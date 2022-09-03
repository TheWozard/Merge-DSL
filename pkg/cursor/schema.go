package cursor

const (
	SchemaTypeKey       = "type"
	SchemaPropertiesKey = "properties"
	SchemaItemsKey      = "items"
	SchemaDefaultKey    = "default"

	SchemaObjectType = "object"
	SchemaArrayType  = "array"
	SchemaEdgeType   = "edge"
)

type (
	schemaObjectCursorPointer struct {
		data map[string]interface{}
	}
	schemaArrayCursorPointer struct {
		data map[string]interface{}
	}
	schemaEdgeCursorPointer struct {
		data map[string]interface{}
	}
)

// NewRawCursor converts raw golang struct based on json data into a traversable tree
func NewSchemaCursor(raw map[string]interface{}) SchemaCursor {
	if typ, ok := raw[SchemaTypeKey]; ok {
		switch typ {
		case SchemaObjectType:
			return schemaObjectCursorPointer{data: raw}
		case SchemaArrayType:
			return schemaArrayCursorPointer{data: raw}
		case SchemaEdgeType:
			return schemaEdgeCursorPointer{data: raw}
		}
	}
	return nil
}

/*
 * schemaObjectCursorPointer
 */

func (o schemaObjectCursorPointer) IsEdge() bool {
	return false
}

func (o schemaObjectCursorPointer) Value() map[string]interface{} {
	return o.data
}

func (o schemaObjectCursorPointer) GetKey(key string) SchemaCursor {
	if properties, ok := o.data[SchemaPropertiesKey].(map[string]interface{}); ok {
		if data, ok := properties[key].(map[string]interface{}); ok {
			return NewSchemaCursor(data)
		}
	}
	return nil
}

func (o schemaObjectCursorPointer) GetKeys() []string {
	results := []string{}
	if properties, ok := o.data[SchemaPropertiesKey].(map[string]interface{}); ok {
		for key := range properties {
			results = append(results, key)
		}
	}
	return results
}

func (o schemaObjectCursorPointer) GetItems() []SchemaCursor {
	return []SchemaCursor{}
}

func (o schemaObjectCursorPointer) GetDefault() SchemaCursor {
	return nil
}

/*
 * schemaArrayCursorPointer
 */

func (a schemaArrayCursorPointer) IsEdge() bool {
	return false
}

func (a schemaArrayCursorPointer) Value() map[string]interface{} {
	return a.data
}

func (a schemaArrayCursorPointer) GetKey(key string) SchemaCursor {
	return nil
}

func (a schemaArrayCursorPointer) GetKeys() []string {
	return []string{}
}

func (a schemaArrayCursorPointer) GetItems() []SchemaCursor {
	results := []SchemaCursor{}
	if items, ok := a.data[SchemaItemsKey].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(map[string]interface{}); ok {
				results = append(results, NewSchemaCursor(data))
			}
		}
	}
	return results
}

func (a schemaArrayCursorPointer) GetDefault() SchemaCursor {
	if data, ok := a.data[SchemaDefaultKey].(map[string]interface{}); ok {
		return NewSchemaCursor(data)
	}
	return nil
}

/*
 * schemaEdgeCursorPointer
 */

func (e schemaEdgeCursorPointer) IsEdge() bool {
	return true
}

func (e schemaEdgeCursorPointer) Value() map[string]interface{} {
	return e.data
}

func (e schemaEdgeCursorPointer) GetKey(key string) SchemaCursor {
	return nil
}

func (e schemaEdgeCursorPointer) GetKeys() []string {
	return []string{}
}

func (e schemaEdgeCursorPointer) GetItems() []SchemaCursor {
	return []SchemaCursor{}
}

func (e schemaEdgeCursorPointer) GetDefault() SchemaCursor {
	return nil
}
