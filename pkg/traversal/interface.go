package traversal

type IdIndexer[T any] interface {
	Parse(item TraversalPointer[T]) interface{}
	Index(pointer TraversalPointer[T]) (map[interface{}]TraversalPointer[T], []TraversalPointer[T])
}

type TraversalPointer[T any] interface {
	IsEdge() bool
	Value() T
	GetKey(key string) TraversalPointer[T]
	GetItems() []TraversalPointer[T]
}
