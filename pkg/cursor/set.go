package cursor

// CursorSet makes handling a list of cursors easier.
type CursorSet[T any] struct {
	Set       []Cursor[T]
	Validator func(T) bool
}

func (s CursorSet[T]) NewSet(set ...Cursor[T]) CursorSet[T] {
	return CursorSet[T]{
		Set:       set,
		Validator: s.Validator,
	}
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
	sets := []Cursor[T]{}
	for _, cursor := range s.Set {
		if next := cursor.GetKey(key); next != nil {
			sets = append(sets, next)
		}
	}
	return s.NewSet(sets...)
}

func (s CursorSet[T]) GetItems() []CursorSet[T] {
	result := []CursorSet[T]{}
	for _, cursor := range s.Set {
		for _, child := range cursor.GetItems() {
			result = append(result, s.NewSet(child))
		}
	}
	return result
}

func (s CursorSet[T]) GetIdsAndExtra(parser IdParser[T]) (map[interface{}]CursorSet[T], []CursorSet[T]) {
	extra := []CursorSet[T]{}
	index := map[interface{}][]Cursor[T]{}
	for _, cursor := range s.Set {
		for _, extra_cursor := range PopulateIndexCursorsById(cursor.GetItems(), parser, index) {
			extra = append(extra, s.NewSet(extra_cursor))
		}
	}
	result := make(map[interface{}]CursorSet[T], len(index))
	for id, set := range index {
		result[id] = s.NewSet(set...)
	}
	return result, extra
}

func (s CursorSet[T]) GetDefault() CursorSet[T] {
	result := []Cursor[T]{}
	for _, cursor := range s.Set {
		if def := cursor.GetDefault(); def != nil {
			result = append(result, def)
		}
	}
	return s.NewSet(result...)
}
