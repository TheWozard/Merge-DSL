package cursor

// Cursor points to a node in a tree during traversal.
// Cursors can be traversed as all other cursors, but will return nil when cursor moves outside of the tree.
type Cursor[T any] interface {
	// IsEdge returns true if the cursor is pointing to an edge node that has no possible children.
	IsEdge() bool
	// Value returns the current value of the cursor.
	Value() T
	// GetKey returns a new Cursor pointing to a child node based on the passed key value.
	GetKey(key string) Cursor[T]
	// GetKeys returns a slice of all possible keys of child nodes.
	GetKeys() []string
	// GetItems returns a slice of all item nodes.
	// Use with IndexCursorsById or an IdParser.Index to get items indexed
	GetItems() []Cursor[T]
	// GetDefault returns the cursor of the default item node.
	GetDefault() Cursor[T]
}

// IdParser can convert a cursor to an Id to be indexed by IndexCursorsById.
type IdParser[T any] interface {
	// Parse extracts the Id from a cursor
	Parse(item Cursor[T]) interface{}
	// Index convenience function for calling IndexCursorsById with the passed cursors and this parser
	Index(cursors []Cursor[T]) (map[interface{}][]Cursor[T], []Cursor[T])
}

type RawData = interface{}
type RawCursor = Cursor[RawData]

type SchemaData = map[string]interface{}
type SchemaCursor = Cursor[SchemaData]
