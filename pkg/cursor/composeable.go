package cursor

type (
	MapPointer[T any] struct {
		Data  T
		Nodes map[string]Cursor[T]
	}
	SlicePointer[T any] struct {
		Data    T
		Items   []Cursor[T]
		Default Cursor[T]
	}
	EdgePointer[T any] struct {
		Data T
	}
)

/*
 * ComposableMapCursorPointer
 */

func (m MapPointer[T]) IsEdge() bool {
	return false
}

func (m MapPointer[T]) Value() T {
	return m.Data
}

func (m MapPointer[T]) GetKey(key string) Cursor[T] {
	if node, ok := m.Nodes[key]; ok {
		return node
	}
	return nil
}

func (m MapPointer[T]) GetKeys() []string {
	result := []string{}
	for key := range m.Nodes {
		result = append(result, key)
	}
	return result
}

func (m MapPointer[T]) GetItems() []Cursor[T] {
	return []Cursor[T]{}
}

func (m MapPointer[T]) GetDefault() Cursor[T] {
	return nil
}

/*
 * ComposableSliceCursorPointer
 */

func (s SlicePointer[T]) IsEdge() bool {
	return false
}

func (s SlicePointer[T]) Value() T {
	return s.Data
}

func (s SlicePointer[T]) GetKey(key string) Cursor[T] {
	return nil
}

func (s SlicePointer[T]) GetKeys() []string {
	return []string{}
}

func (s SlicePointer[T]) GetItems() []Cursor[T] {
	return s.Items
}

func (s SlicePointer[T]) GetDefault() Cursor[T] {
	return s.Default
}

/*
 * ComposableEdgeCursorPointer
 */

func (e EdgePointer[T]) IsEdge() bool {
	return true
}

func (e EdgePointer[T]) Value() T {
	return e.Data
}

func (e EdgePointer[T]) GetKey(key string) Cursor[T] {
	return nil
}

func (e EdgePointer[T]) GetKeys() []string {
	return []string{}
}

func (e EdgePointer[T]) GetItems() []Cursor[T] {
	return []Cursor[T]{}
}

func (e EdgePointer[T]) GetDefault() Cursor[T] {
	return nil
}
