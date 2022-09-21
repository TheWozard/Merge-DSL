package cursor

// Cursor points to a node in a tree during traversal.
// Cursors can be traversed as all other cursors, but will return nil when cursor moves outside of the tree.
type Cursor[T any] interface {
	// General

	// HasChildren returns if the current node has any children. How to access the children depends on the node.
	HasChildren() bool
	// Value returns the current value of the cursor.
	Value() T

	// Objects

	// GetKey returns a new Cursor pointing to a child node based on the passed key value.
	GetKey(key string) Cursor[T]
	// GetKeys returns a slice of all possible keys of child nodes.
	GetKeys() []string

	// Arrays

	// GetItems returns a slice of all item nodes.
	GetItems() []Cursor[T]
	// GetDefault returns the cursor of the default item node.
	GetDefault() Cursor[T]
}
