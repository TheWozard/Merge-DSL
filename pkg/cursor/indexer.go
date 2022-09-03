package cursor

var (
	DefaultRawIndexer    = RawCursorIdParser[interface{}]{Key: "id"}
	DefaultSchemaIndexer = SchemaCursorIdParser{Key: "id"}
)

// IndexCursorsById indexes a slice of cursors based on the parser. Returns any cursors that do not contain an id
// in an extra slice.
func IndexCursorsById[T any](cursors []Cursor[T], parser IdParser[T]) (map[interface{}]Cursor[T], []Cursor[T]) {
	// TODO: Would it be better to return an index of lists. That way all matches could be merged.
	// Order would have to be preserved to keep merge consistent.
	index := map[interface{}]Cursor[T]{}
	extras := []Cursor[T]{}
	// Panic safety
	if len(cursors) == 0 {
		return index, extras
	}
	if parser == nil {
		return index, cursors
	}
	// Indexing
	for _, cursor := range cursors {
		if cursor != nil {
			if id := parser.Parse(cursor); id != nil {
				index[id] = cursor
			} else {
				extras = append(extras, cursor)
			}
		}
	}
	return index, extras
}

// RawCursorIdParser gets an id from a cursor based on GetKey.
// By default requires the node be an edge node.
type RawCursorIdParser[T any] struct {
	// Key to get the id from.
	Key string
	// If true, will allow ids to be non edge nodes.
	IgnoreIsEdge bool
}

func (k RawCursorIdParser[T]) Parse(cursor Cursor[T]) interface{} {
	if node := cursor.GetKey(k.Key); node != nil {
		if k.IgnoreIsEdge || node.IsEdge() {
			return node.Value()
		}
	}
	return nil
}

func (k RawCursorIdParser[T]) Index(cursors []Cursor[T]) (map[interface{}]Cursor[T], []Cursor[T]) {
	return IndexCursorsById[T](cursors, k)
}

// SchemaCursorIdParser gets an id from the value of the passed cursor.
type SchemaCursorIdParser struct {
	// The key to pull the id from in the index
	Key string
}

func (s SchemaCursorIdParser) Parse(cursor SchemaCursor) interface{} {
	if value := cursor.Value(); value != nil {
		if key, ok := value[s.Key]; ok {
			return key
		}
	}
	return nil
}

func (s SchemaCursorIdParser) Index(cursors []SchemaCursor) (map[interface{}]SchemaCursor, []SchemaCursor) {
	return IndexCursorsById[SchemaData](cursors, s)
}
