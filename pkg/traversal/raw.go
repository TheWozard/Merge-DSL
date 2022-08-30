package traversal

type (
	rawMapTraversalPointer struct {
		nodes map[string]interface{}
	}
	rawSliceTraversalPointer struct {
		items []interface{}
	}
	rawEdgeTraversalPointer struct {
		data interface{}
	}
)

// NewRawTraversal converts raw golang struct based on json data into a traversable tree
func NewRawTraversal(raw interface{}) Pointer[interface{}] {
	switch typed := raw.(type) {
	case map[string]interface{}:
		return rawMapTraversalPointer{nodes: typed}
	case []interface{}:
		return rawSliceTraversalPointer{items: typed}
	}
	return rawEdgeTraversalPointer{data: raw}
}

/*
 * rawMapTraversalPointer
 */

func (o rawMapTraversalPointer) IsEdge() bool {
	return false
}

func (o rawMapTraversalPointer) Value() interface{} {
	return o.nodes
}

func (o rawMapTraversalPointer) GetKey(key string) Pointer[interface{}] {
	if data, ok := o.nodes[key]; ok {
		return NewRawTraversal(data)
	}
	return nil
}

func (o rawMapTraversalPointer) GetItems() []Pointer[interface{}] {
	return []Pointer[interface{}]{}
}

/*
 * rawSliceTraversalPointer
 */

func (a rawSliceTraversalPointer) IsEdge() bool {
	return false
}

func (a rawSliceTraversalPointer) Value() interface{} {
	return a.items
}

func (a rawSliceTraversalPointer) GetKey(key string) Pointer[interface{}] {
	return nil
}

func (a rawSliceTraversalPointer) GetItems() []Pointer[interface{}] {
	items := []Pointer[interface{}]{}
	for _, item := range a.items {
		items = append(items, NewRawTraversal(item))
	}
	return items
}

/*
 * rawEdgeTraversalPointer
 */

func (e rawEdgeTraversalPointer) IsEdge() bool {
	return true
}

func (e rawEdgeTraversalPointer) Value() interface{} {
	return e.data
}

func (e rawEdgeTraversalPointer) GetKey(key string) Pointer[interface{}] {
	return nil
}

func (e rawEdgeTraversalPointer) GetItems() []Pointer[interface{}] {
	return []Pointer[interface{}]{}
}
