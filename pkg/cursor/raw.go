package cursor

type (
	rawMapCursorPointer struct {
		nodes map[string]interface{}
	}
	rawSliceCursorPointer struct {
		items []interface{}
	}
	rawEdgeCursorPointer struct {
		data interface{}
	}
)

// NewRawCursor converts raw golang struct based on json data into a traversable tree
func NewRawCursor(raw interface{}) RawCursor {
	switch typed := raw.(type) {
	case map[string]interface{}:
		return rawMapCursorPointer{nodes: typed}
	case []interface{}:
		return rawSliceCursorPointer{items: typed}
	}
	return rawEdgeCursorPointer{data: raw}
}

/*
 * rawMapCursorPointer
 */

func (m rawMapCursorPointer) IsEdge() bool {
	return false
}

func (m rawMapCursorPointer) Value() interface{} {
	return m.nodes
}

func (m rawMapCursorPointer) GetKey(key string) RawCursor {
	if data, ok := m.nodes[key]; ok {
		return NewRawCursor(data)
	}
	return nil
}

func (m rawMapCursorPointer) GetKeys() []string {
	result := []string{}
	for key := range m.nodes {
		result = append(result, key)
	}
	return result
}

func (m rawMapCursorPointer) GetItems() []RawCursor {
	return []RawCursor{}
}

func (m rawMapCursorPointer) GetDefault() RawCursor {
	return nil
}

/*
 * rawSliceCursorPointer
 */

func (s rawSliceCursorPointer) IsEdge() bool {
	return false
}

func (s rawSliceCursorPointer) Value() interface{} {
	return s.items
}

func (s rawSliceCursorPointer) GetKey(key string) RawCursor {
	return nil
}

func (s rawSliceCursorPointer) GetKeys() []string {
	return []string{}
}

func (s rawSliceCursorPointer) GetItems() []RawCursor {
	items := []RawCursor{}
	for _, item := range s.items {
		items = append(items, NewRawCursor(item))
	}
	return items
}

func (s rawSliceCursorPointer) GetDefault() RawCursor {
	return nil
}

/*
 * rawEdgeCursorPointer
 */

func (e rawEdgeCursorPointer) IsEdge() bool {
	return true
}

func (e rawEdgeCursorPointer) Value() interface{} {
	return e.data
}

func (e rawEdgeCursorPointer) GetKey(key string) RawCursor {
	return nil
}

func (e rawEdgeCursorPointer) GetKeys() []string {
	return []string{}
}

func (e rawEdgeCursorPointer) GetItems() []RawCursor {
	return []RawCursor{}
}

func (e rawEdgeCursorPointer) GetDefault() RawCursor {
	return nil
}
