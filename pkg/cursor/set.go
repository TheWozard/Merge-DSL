package cursor

// CursorSet makes handling a list of cursors easier.
type CursorSet[T any] struct {
	Set       []Cursor[T]
	Validator func(T) bool
}

func (s CursorSet[T]) IsEdge() bool {
	if len(s.Set) > 0 {
		return s.Set[0].IsEdge()
	}
	return true
}

func (s CursorSet[T]) Value() T {
	for _, cursor := range s.Set {
		if value := cursor.Value(); s.Validator(value) {
			return value
		}
	}
	var empty T
	return empty
}

func (s CursorSet[T]) AllValues() []T {
	values := []T{}
	for _, cursor := range s.Set {
		if value := cursor.Value(); s.Validator(value) {
			values = append(values, value)
		}
	}
	return values
}

func (s CursorSet[T]) GetKey(key string) CursorSet[T] {
	result := CursorSet[T]{
		Set:       []Cursor[T]{},
		Validator: s.Validator,
	}
	for _, cursor := range s.Set {
		if next := cursor.GetKey(key); next != nil {
			result.Set = append(result.Set, next)
		}
	}
	return result
}

func (s CursorSet[T]) GetItems() []CursorSet[T] {
	result := []CursorSet[T]{}
	for _, cursor := range s.Set {
		for _, child := range cursor.GetItems() {
			result = append(result, CursorSet[T]{Set: []Cursor[T]{child}, Validator: s.Validator})
		}
	}
	return result
}
