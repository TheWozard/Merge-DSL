package cursor

// Set makes handling a list of cursors easier.
type Set[T any] []Cursor[T]

// Value gets the first value in the set that passes the validator.
// if not value is found the second value returns false.
func (s Set[T]) Value(validator Validator[T]) (T, bool) {
	for _, cursor := range s {
		if value := cursor.Value(); validator(value) {
			return value, true
		}
	}
	var empty T
	return empty, false
}

// GetKey creates a new Set based on the results of traversing all cursors in the current set
// to the given key. Any cursors that become nil are dropped.
func (s Set[T]) GetKey(key string) Set[T] {
	cursors := []Cursor[T]{}
	for _, cursor := range s {
		if next := cursor.GetKey(key); next != nil {
			cursors = append(cursors, next)
		}
	}
	return Set[T](cursors)
}

// GetItems creates a slice of all items for all cursors in order.
func (s Set[T]) GetItems() []Cursor[T] {
	cursors := []Cursor[T]{}
	for _, cursor := range s {
		cursors = append(cursors, cursor.GetItems()...)
	}
	return cursors
}

// GetGroupedItems groups GetItems into a slice Sets of grouped cursors based on the grouper
func (s Set[T]) GetGroupedItems(grouper Grouper[T]) []Set[T] {
	final := []Set[T]{}
	groups := GroupCursors(s.GetItems(), grouper)
	for _, group := range groups {
		final = append(final, Set[T](group))
	}
	return final
}

func (s Set[T]) GetIndexedItems(grouper Grouper[T]) (map[interface{}]Set[T], []Set[T]) {
	final_extra := []Set[T]{}
	final_index := map[interface{}]Set[T]{}
	index, extra := IndexCursors(s.GetItems(), grouper)
	for key, cursors := range index {
		final_index[key] = Set[T](cursors)
	}
	for _, extra := range extra {
		final_extra = append(final_extra, Set[T]{extra})
	}
	return final_index, final_extra
}

func (s Set[T]) GetDefault() Set[T] {
	result := []Cursor[T]{}
	for _, cursor := range s {
		if def := cursor.GetDefault(); def != nil {
			result = append(result, def)
		}
	}
	return Set[T](result)
}
