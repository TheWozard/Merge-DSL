package traversal

// PointerSet makes handling a list of pointers easier.
type PointerSet[T any] struct {
	Set       []Pointer[T]
	Validator func(T) bool
}

func (s PointerSet[T]) IsEdge() bool {
	if len(s.Set) > 0 {
		return s.Set[0].IsEdge()
	}
	return true
}

func (s PointerSet[T]) Value() T {
	for _, pointer := range s.Set {
		if value := pointer.Value(); s.Validator(value) {
			return value
		}
	}
	var empty T
	return empty
}

func (s PointerSet[T]) AllValues() []T {
	values := []T{}
	for _, pointer := range s.Set {
		if value := pointer.Value(); s.Validator(value) {
			values = append(values, value)
		}
	}
	return values
}

func (s PointerSet[T]) GetKey(key string) PointerSet[T] {
	result := PointerSet[T]{
		Set:       []Pointer[T]{},
		Validator: s.Validator,
	}
	for _, pointer := range s.Set {
		if next := pointer.GetKey(key); next != nil {
			result.Set = append(result.Set, next)
		}
	}
	return result
}

func (s PointerSet[T]) GetItems() []PointerSet[T] {
	result := []PointerSet[T]{}
	for _, pointer := range s.Set {
		for _, child := range pointer.GetItems() {
			result = append(result, PointerSet[T]{Set: []Pointer[T]{child}, Validator: s.Validator})
		}
	}
	return result
}
