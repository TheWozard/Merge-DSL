package cursor

var (
	DefaultRawGrouper    = KeyGrouper{Key: "id"}.Group
	DefaultSchemaGrouper = SchemaDataKeyGrouper{Key: "id"}.Group
)

// Grouper defines a function that extracts an ID for grouping cursors by
type Grouper[T any] func(cursor Cursor[T]) interface{}

// GroupCursors groups the passed cursor using the grouper to build group ids.
// This preserves order of cursors, but requires iteration to find an id.
// Returns a slice of groups of cursors. These are grouped based on the first occurrence in the passed slice of cursors.
func GroupCursors[T any](cursors []Cursor[T], grouper Grouper[T]) [][]Cursor[T] {
	groups := [][]Cursor[T]{}
	index := map[interface{}]int{}

	// Indexing
	for _, cursor := range cursors {
		if cursor != nil {
			if grouper != nil {
				if id := grouper(cursor); id != nil {
					i, ok := index[id]
					if !ok {
						i = len(groups)
						groups = append(groups, []Cursor[T]{})
						index[id] = i
					}
					groups[i] = append(groups[i], cursor)
					continue
				}
			}
			groups = append(groups, []Cursor[T]{cursor})
		}
	}
	return groups
}

// IndexCursors indexes the cursors using the grouper to build group ids.
// This allows for random access, but does not preserve order of the ids.
// Returns an index of ids to Cursors and a list of all cursors without an id.
func IndexCursors[T any](cursors []Cursor[T], grouper Grouper[T]) (map[interface{}][]Cursor[T], []Cursor[T]) {
	index := map[interface{}][]Cursor[T]{}
	extra := []Cursor[T]{}

	if grouper == nil {
		return index, cursors
	}

	for _, cursor := range cursors {
		if cursor != nil {
			if id := grouper(cursor); id != nil {
				group, ok := index[id]
				if !ok {
					group = []Cursor[T]{}
				}
				index[id] = append(group, cursor)
				continue
			}
			extra = append(extra, cursor)
		}
	}
	return index, extra
}

// KeyGrouper gets an id from a cursor based on GetKey.
type KeyGrouper struct {
	// Key to get the id from.
	Key string
}

func (k KeyGrouper) Group(cursor Cursor[interface{}]) interface{} {
	if node := cursor.GetKey(k.Key); node != nil {
		return node.Value()
	}
	return nil
}

// SchemaDataKeyGrouper gets an id from the value of the passed SchemaData cursor.
type SchemaDataKeyGrouper struct {
	// Key to get the id from.
	Key string
}

func (s SchemaDataKeyGrouper) Group(cursor Cursor[SchemaData]) interface{} {
	if value := cursor.Value(); value != nil {
		if key, ok := value[s.Key]; ok {
			return key
		}
	}
	return nil
}
