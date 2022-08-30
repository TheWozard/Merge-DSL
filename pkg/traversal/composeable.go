package traversal

type (
	MapPointer[T any] struct {
		Data  T
		Nodes map[string]Pointer[T]
	}
	SlicePointer[T any] struct {
		Data  T
		Items []Pointer[T]
	}
	EdgePointer[T any] struct {
		Data T
	}
)

/*
 * ComposableMapTraversalPointer
 */

func (m MapPointer[T]) IsEdge() bool {
	return false
}

func (m MapPointer[T]) Value() T {
	return m.Data
}

func (m MapPointer[T]) GetKey(key string) Pointer[T] {
	if node, ok := m.Nodes[key]; ok {
		return node
	}
	return nil
}

func (m MapPointer[T]) GetItems() []Pointer[T] {
	return []Pointer[T]{}
}

/*
 * ComposableSliceTraversalPointer
 */

func (c SlicePointer[T]) IsEdge() bool {
	return false
}

func (c SlicePointer[T]) Value() T {
	return c.Data
}

func (c SlicePointer[T]) GetKey(key string) Pointer[T] {
	return nil
}

func (c SlicePointer[T]) GetItems() []Pointer[T] {
	return c.Items
}

/*
 * ComposableEdgeTraversalPointer
 */

func (e EdgePointer[T]) IsEdge() bool {
	return true
}

func (e EdgePointer[T]) Value() T {
	return e.Data
}

func (e EdgePointer[T]) GetKey(key string) Pointer[T] {
	return nil
}

func (e EdgePointer[T]) GetItems() []Pointer[T] {
	return []Pointer[T]{}
}
