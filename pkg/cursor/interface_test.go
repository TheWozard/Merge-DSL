package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"

	"github.com/stretchr/testify/assert"
)

type CursorCase[T any] struct {
	IsEdge  bool
	Value   T
	Key     map[string]cursor.Cursor[T]
	Keys    []string
	Items   []cursor.Cursor[T]
	Default cursor.Cursor[T]
}

func CursorTestSuite[T any](t *testing.T, cursor cursor.Cursor[T], expectations CursorCase[T]) {
	assert.Equal(t, expectations.IsEdge, cursor.IsEdge())
	assert.Equal(t, expectations.Value, cursor.Value())
	for key, value := range expectations.Key {
		assert.Equal(t, value, cursor.GetKey(key), key)
	}
	assert.Equal(t, expectations.Keys, cursor.GetKeys())
	assert.Equal(t, expectations.Items, cursor.GetItems())
	assert.Equal(t, expectations.Default, cursor.GetDefault())
}
