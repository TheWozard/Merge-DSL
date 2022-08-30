package traversal

// Iterates over the pointer.getItems() for a pointer. Uses the passed parser to
// create an index of all items with a non nil parsed id value. In the case of duplicate
// id, last one wins.
// Returns both the index as well as a slice of all pointers without an id.
// nil pointer returns an empty index and an empty extra
// nil parser returns empty index and pointer.GetItems()
func IndexPointerItemsById[T any](pointer Pointer[T], parser IdIndexer[T]) (map[interface{}]Pointer[T], []Pointer[T]) {
	// TODO: Would it be better to return an index of lists. That way all matches could be merged.
	// Order would have to be preserved to keep merge consistent.
	index := map[interface{}]Pointer[T]{}
	extra := []Pointer[T]{}
	// Panic safety
	if pointer == nil {
		return index, extra
	}
	if parser == nil {
		return index, pointer.GetItems()
	}
	// Indexing
	for _, item := range pointer.GetItems() {
		if id := parser.Parse(item); id != nil {
			index[id] = item
		} else {
			extra = append(extra, item)
		}
	}
	return index, extra
}

// PointerKeyIdIndexer gets an id from a pointers key.
// By default requires the node be an edge node.
type PointerKeyIdIndexer[T any] struct {
	// Key to get the id from.
	Key string
	// If true, will allow ids to be non edge nodes.
	IgnoreIsEdge bool
}

func (k PointerKeyIdIndexer[T]) Parse(pointer Pointer[T]) interface{} {
	if node := pointer.GetKey(k.Key); node != nil {
		if k.IgnoreIsEdge || node.IsEdge() {
			return node.Value()
		}
	}
	return nil
}

func (k PointerKeyIdIndexer[T]) Index(pointer Pointer[T]) (map[interface{}]Pointer[T], []Pointer[T]) {
	return IndexPointerItemsById[T](pointer, k)
}
