package traversal

type safeTraversalPointer[T any] struct {
	pointer TraversalPointer[T]
}

// Wraps a given pointer in a safety net that allows it to traverse outside of the root tree safely.
// It is generally preferred to check if a pointer is nil and drop any nil pointers as it is more efficient.
func Safe[T any](pointer TraversalPointer[T]) TraversalPointer[T] {
	return safeTraversalPointer[T]{pointer: pointer}
}

func (s safeTraversalPointer[T]) IsEdge() bool {
	if s.pointer != nil {
		return s.pointer.IsEdge()
	}
	return true
}

func (s safeTraversalPointer[T]) Value() T {
	if s.pointer != nil {
		return s.pointer.Value()
	}
	var empty T
	return empty
}

func (s safeTraversalPointer[T]) GetKey(key string) TraversalPointer[T] {
	result := safeTraversalPointer[T]{}
	if s.pointer != nil {
		if next := s.pointer.GetKey(key); next != nil {
			result.pointer = next
		}
	}
	return result
}

func (s safeTraversalPointer[T]) GetItems() []TraversalPointer[T] {
	final := []TraversalPointer[T]{}
	if s.pointer != nil {
		for _, item := range s.pointer.GetItems() {
			final = append(final, safeTraversalPointer[T]{pointer: item})
		}
	}
	return final
}
