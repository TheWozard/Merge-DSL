package traversal

type IdIndexer[T any] interface {
	Parse(item Pointer[T]) interface{}
	Index(pointer Pointer[T]) (map[interface{}]Pointer[T], []Pointer[T])
}

type Pointer[T any] interface {
	IsEdge() bool
	Value() T
	GetKey(key string) Pointer[T]
	GetItems() []Pointer[T]
}
