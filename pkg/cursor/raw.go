package cursor

type (
	rawMapCursorPointer struct {
		parent Cursor[interface{}]
		nodes  map[string]interface{}
	}
	rawSliceCursorPointer struct {
		parent Cursor[interface{}]
		items  []interface{}
	}
	rawEdgeCursorPointer struct {
		parent Cursor[interface{}]
		data   interface{}
	}
)

// NewRawCursor starts a cursor pointing to an interface.
func NewRawCursor(raw interface{}) Cursor[interface{}] {
	return NewRawCursorFrom(raw, nil)
}

// NewRawCursorFrom starts a cursor pointing to an interface that came from the passed parent.
func NewRawCursorFrom(raw interface{}, parent Cursor[interface{}]) Cursor[interface{}] {
	switch typed := raw.(type) {
	case map[string]interface{}:
		return rawMapCursorPointer{nodes: typed, parent: parent}
	case []interface{}:
		return rawSliceCursorPointer{items: typed, parent: parent}
	}
	return rawEdgeCursorPointer{data: raw, parent: parent}
}

/*
 * rawMapCursorPointer
 */

func (m rawMapCursorPointer) Parent() Cursor[interface{}] {
	return m.parent
}

func (m rawMapCursorPointer) HasChildren() bool {
	return len(m.nodes) > 0
}

func (m rawMapCursorPointer) Value() interface{} {
	return m.nodes
}

func (m rawMapCursorPointer) GetKey(key string) Cursor[interface{}] {
	if data, ok := m.nodes[key]; ok {
		return NewRawCursorFrom(data, m)
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

func (m rawMapCursorPointer) GetItems() []Cursor[interface{}] {
	return []Cursor[interface{}]{}
}

func (m rawMapCursorPointer) GetDefault() Cursor[interface{}] {
	return nil
}

/*
 * rawSliceCursorPointer
 */

func (s rawSliceCursorPointer) Parent() Cursor[interface{}] {
	return s.parent
}

func (s rawSliceCursorPointer) HasChildren() bool {
	return len(s.items) > 0
}

func (s rawSliceCursorPointer) Value() interface{} {
	return s.items
}

func (s rawSliceCursorPointer) GetKey(key string) Cursor[interface{}] {
	return nil
}

func (s rawSliceCursorPointer) GetKeys() []string {
	return []string{}
}

func (s rawSliceCursorPointer) GetItems() []Cursor[interface{}] {
	items := []Cursor[interface{}]{}
	for _, item := range s.items {
		items = append(items, NewRawCursorFrom(item, s))
	}
	return items
}

func (s rawSliceCursorPointer) GetDefault() Cursor[interface{}] {
	return nil
}

/*
 * rawEdgeCursorPointer
 */

func (e rawEdgeCursorPointer) Parent() Cursor[interface{}] {
	return e.parent
}

func (e rawEdgeCursorPointer) HasChildren() bool {
	return false
}

func (e rawEdgeCursorPointer) Value() interface{} {
	return e.data
}

func (e rawEdgeCursorPointer) GetKey(key string) Cursor[interface{}] {
	return nil
}

func (e rawEdgeCursorPointer) GetKeys() []string {
	return []string{}
}

func (e rawEdgeCursorPointer) GetItems() []Cursor[interface{}] {
	return []Cursor[interface{}]{}
}

func (e rawEdgeCursorPointer) GetDefault() Cursor[interface{}] {
	return nil
}
