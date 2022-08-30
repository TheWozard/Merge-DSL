package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"

	"github.com/stretchr/testify/assert"
)

type StandardCases[T any] struct {
	IsEdge bool
	Value  T
	Keys   map[string]traversal.Pointer[T]
	Items  []traversal.Pointer[T]
}

func PointerTestSuite[T any](t *testing.T, pointer traversal.Pointer[T], cases StandardCases[T]) {
	assert.Equal(t, cases.IsEdge, pointer.IsEdge())
	assert.Equal(t, cases.Value, pointer.Value())
	for key, value := range cases.Keys {
		assert.Equal(t, value, pointer.GetKey(key), key)
	}
	assert.Equal(t, cases.Items, pointer.GetItems())
}
