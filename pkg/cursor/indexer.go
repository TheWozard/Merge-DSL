package cursor

var (
	DefaultRawIndexer    = RawCursorIdParser[interface{}]{Key: "id"}
	DefaultSchemaIndexer = SchemaCursorIdParser{Key: "id"}
)

// PopulateIndexCursorsById indexes a slice of cursors based on the parser.
// Returns a list of the order of the ids as received from the cursors. Will not return ids that already existed in the index.
// Returns any cursors that do not contain an id in an extra slice.
func PopulateIndexCursorsById[T any](cursors []Cursor[T], parser IdParser[T], index map[interface{}][]Cursor[T]) ([]interface{}, []Cursor[T]) {
	extras := []Cursor[T]{}
	// We need to preserve order of ids despite
	ids := []interface{}{}
	// Panic safety
	if len(cursors) == 0 {
		return ids, extras
	}
	if parser == nil {
		return ids, cursors
	}
	// Indexing
	for _, cursor := range cursors {
		if cursor != nil {
			if id := parser.Parse(cursor); id != nil {
				sets, ok := index[id]
				if !ok {
					sets = []Cursor[T]{}
					ids = append(ids, id)
				}
				sets = append(sets, cursor)
				index[id] = sets
			} else {
				extras = append(extras, cursor)
			}
		}
	}
	return ids, extras
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

func (k RawCursorIdParser[T]) Index(cursors []Cursor[T]) (map[interface{}][]Cursor[T], []Cursor[T]) {
	index := map[interface{}][]Cursor[T]{}
	_, extra := PopulateIndexCursorsById[T](cursors, k, index)
	return index, extra
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

func (s SchemaCursorIdParser) Index(cursors []SchemaCursor) (map[interface{}][]SchemaCursor, []SchemaCursor) {
	index := map[interface{}][]SchemaCursor{}
	_, extra := PopulateIndexCursorsById[SchemaData](cursors, s, index)
	return index, extra
}
