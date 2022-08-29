package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSafeTraversal(t *testing.T) {
	testCases := []struct {
		desc  string
		input traversal.TraversalPointer[interface{}]
		get   func(pointer traversal.TraversalPointer[interface{}]) interface{}
		value interface{}
	}{
		// nil
		{
			desc:  "nil",
			input: nil,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("foobar").Value()
			},
			value: nil,
		},
		{
			desc:  "nil_value",
			input: nil,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.Value()
			},
			value: nil,
		},
		{
			desc:  "nil_isEdge",
			input: nil,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.IsEdge()
			},
			value: true,
		},
		{
			desc:  "nil_getKey",
			input: nil,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("foobar").Value()
			},
			value: nil,
		},
		{
			desc:  "nil_getItems",
			input: nil,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[interface{}]{},
		},

		// object
		{
			desc:  "out_of_object",
			input: traversal.NewRawTraversal(map[string]interface{}{}),
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("foobar").GetKey("next").Value()
			},
			value: nil,
		},
		{
			desc: "out_of_object_key",
			input: traversal.NewRawTraversal(map[string]interface{}{
				"foobar": "bad",
			}),
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("foobar").GetKey("next").Value()
			},
			value: nil,
		},
		{
			desc: "object_value",
			input: traversal.NewRawTraversal(map[string]interface{}{
				"foobar": "success",
			}),
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("foobar").Value()
			},
			value: "success",
		},
		{
			desc: "object_isEdge",
			input: traversal.NewRawTraversal(map[string]interface{}{
				"foobar": "success",
			}),
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.IsEdge()
			},
			value: false,
		},

		// array
		{
			desc:  "array_getItems",
			input: traversal.NewRawTraversal([]interface{}{"example", "data"}),
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[interface{}]{
				traversal.Safe(traversal.NewRawTraversal("example")),
				traversal.Safe(traversal.NewRawTraversal("data")),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.value, tC.get(traversal.Safe(tC.input)))
		})
	}
}
