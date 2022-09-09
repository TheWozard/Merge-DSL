package cursor

// CursorSet makes handling a list of cursors easier.
type CursorSet[T any] []Cursor[T]

func (s CursorSet[T]) NewSet(set ...Cursor[T]) CursorSet[T] {
	return set
}

func (s CursorSet[T]) IsEdge() bool {
	if len(s) > 0 {
		return s[0].IsEdge()
	}
	return true
}

func (s CursorSet[T]) Value(validator func(T) bool) T {
	for _, cursor := range s {
		if value := cursor.Value(); validator(value) {
			return value
		}
	}
	var empty T
	return empty
}

func (s CursorSet[T]) GetKey(key string) CursorSet[T] {
	sets := []Cursor[T]{}
	for _, cursor := range s {
		if next := cursor.GetKey(key); next != nil {
			sets = append(sets, next)
		}
	}
	return s.NewSet(sets...)
}

func (s CursorSet[T]) GetItems() []CursorSet[T] {
	result := []CursorSet[T]{}
	for _, cursor := range s {
		for _, child := range cursor.GetItems() {
			result = append(result, s.NewSet(child))
		}
	}
	return result
}

func (s CursorSet[T]) GetIdsAndExtra(parser IdParser[T]) (map[interface{}]CursorSet[T], []interface{}, []CursorSet[T]) {
	extra := []CursorSet[T]{}
	order := []interface{}{}
	index := map[interface{}][]Cursor[T]{}
	for _, cursor := range s {
		added_order, extra_cursors := PopulateIndexCursorsById(cursor.GetItems(), parser, index)
		order = append(order, added_order...)
		for _, extra_cursor := range extra_cursors {
			extra = append(extra, s.NewSet(extra_cursor))
		}
	}
	result := make(map[interface{}]CursorSet[T], len(index))
	for id, set := range index {
		result[id] = s.NewSet(set...)
	}
	return result, order, extra
}

func (s CursorSet[T]) GetDefault() CursorSet[T] {
	result := []Cursor[T]{}
	for _, cursor := range s {
		if def := cursor.GetDefault(); def != nil {
			result = append(result, def)
		}
	}
	return s.NewSet(result...)
}
