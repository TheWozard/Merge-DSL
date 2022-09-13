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
	SchemaData = map[string]interface{}

	schemaObjectCursorPointer struct {
		parent Cursor[SchemaData]
		data   map[string]interface{}
	}
	schemaArrayCursorPointer struct {
		parent Cursor[SchemaData]
		data   map[string]interface{}
	}
	schemaEdgeCursorPointer struct {
		parent Cursor[SchemaData]
		data   map[string]interface{}
	}
)

// NewRawCursor converts raw golang struct based on json data into a traversable tree
func NewSchemaCursor(raw SchemaData) Cursor[SchemaData] {
	return NewSchemaCursorFrom(raw, nil)
}

func NewSchemaCursorFrom(raw SchemaData, parent Cursor[SchemaData]) Cursor[SchemaData] {
	if typ, ok := raw[SchemaTypeKey]; ok {
		switch typ {
		case SchemaObjectType:
			return schemaObjectCursorPointer{data: raw, parent: parent}
		case SchemaArrayType:
			return schemaArrayCursorPointer{data: raw, parent: parent}
		case SchemaEdgeType:
			return schemaEdgeCursorPointer{data: raw, parent: parent}
		}
	}
	return nil
}

/*
 * schemaObjectCursorPointer
 */

func (o schemaObjectCursorPointer) Parent() Cursor[SchemaData] {
	return o.parent
}

func (o schemaObjectCursorPointer) HasChildren() bool {
	if properties, ok := o.data[SchemaPropertiesKey].(map[string]interface{}); ok {
		return len(properties) > 0
	}
	return false
}

func (o schemaObjectCursorPointer) Value() map[string]interface{} {
	return o.data
}

func (o schemaObjectCursorPointer) GetKey(key string) Cursor[SchemaData] {
	if properties, ok := o.data[SchemaPropertiesKey].(map[string]interface{}); ok {
		if data, ok := properties[key].(SchemaData); ok {
			return NewSchemaCursorFrom(data, o)
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

func (o schemaObjectCursorPointer) GetItems() []Cursor[SchemaData] {
	return []Cursor[SchemaData]{}
}

func (o schemaObjectCursorPointer) GetDefault() Cursor[SchemaData] {
	return nil
}

/*
 * schemaArrayCursorPointer
 */

func (a schemaArrayCursorPointer) Parent() Cursor[SchemaData] {
	return a.parent
}

func (a schemaArrayCursorPointer) HasChildren() bool {
	if items, ok := a.data[SchemaItemsKey].([]interface{}); ok {
		if len(items) > 0 {
			return true
		}
	}
	_, ok := a.data[SchemaDefaultKey].(map[string]interface{})
	return ok
}

func (a schemaArrayCursorPointer) Value() map[string]interface{} {
	return a.data
}

func (a schemaArrayCursorPointer) GetKey(key string) Cursor[SchemaData] {
	return nil
}

func (a schemaArrayCursorPointer) GetKeys() []string {
	return []string{}
}

func (a schemaArrayCursorPointer) GetItems() []Cursor[SchemaData] {
	results := []Cursor[SchemaData]{}
	if items, ok := a.data[SchemaItemsKey].([]interface{}); ok {
		for _, item := range items {
			if data, ok := item.(SchemaData); ok {
				results = append(results, NewSchemaCursorFrom(data, a))
			}
		}
	}
	return results
}

func (a schemaArrayCursorPointer) GetDefault() Cursor[SchemaData] {
	if data, ok := a.data[SchemaDefaultKey].(map[string]interface{}); ok {
		return NewSchemaCursorFrom(data, a)
	}
	return nil
}

/*
 * schemaEdgeCursorPointer
 */

func (e schemaEdgeCursorPointer) Parent() Cursor[SchemaData] {
	return e.parent
}

func (e schemaEdgeCursorPointer) HasChildren() bool {
	return false
}

func (e schemaEdgeCursorPointer) Value() map[string]interface{} {
	return e.data
}

func (e schemaEdgeCursorPointer) GetKey(key string) Cursor[SchemaData] {
	return nil
}

func (e schemaEdgeCursorPointer) GetKeys() []string {
	return []string{}
}

func (e schemaEdgeCursorPointer) GetItems() []Cursor[SchemaData] {
	return []Cursor[SchemaData]{}
}

func (e schemaEdgeCursorPointer) GetDefault() Cursor[SchemaData] {
	return nil
}
