package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewRawTraversal(t *testing.T) {
	stringExample := "data"
	objectExample := map[string]interface{}{
		"text": "data",
	}
	arrayExample := []interface{}{
		map[string]interface{}{
			"id": 0,
		},
		"stringy",
	}

	testCases := []struct {
		desc  string
		input interface{}
		get   func(pointer traversal.TraversalPointer[interface{}]) interface{}
		value interface{}
	}{
		// Edge
		{
			desc:  "string_value",
			input: stringExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.Value()
			},
			value: "data",
		},
		{
			desc:  "string_isEdge",
			input: stringExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.IsEdge()
			},
			value: true,
		},
		{
			desc:  "string_getKey",
			input: stringExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("test")
			},
			value: nil,
		},
		{
			desc:  "string_getItems",
			input: stringExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[interface{}]{},
		},

		// Object
		{
			desc:  "object_value",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.Value()
			},
			value: map[string]interface{}{
				"text": "data",
			},
		},
		{
			desc:  "object_isEdge",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.IsEdge()
			},
			value: false,
		},
		{
			desc:  "object_getKey",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("text").Value()
			},
			value: "data",
		},
		{
			desc:  "object_getItems",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[interface{}]{},
		},

		// Array
		{
			desc:  "array_value",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.Value()
			},
			value: []interface{}{
				map[string]interface{}{
					"id": 0,
				},
				"stringy",
			},
		},
		{
			desc:  "array_isEdge",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.IsEdge()
			},
			value: false,
		},
		{
			desc:  "array_getKey",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetKey("text")
			},
			value: nil,
		},
		{
			desc:  "array_getItems",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[interface{}]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[interface{}]{
				traversal.NewRawTraversal(map[string]interface{}{
					"id": 0,
				}),
				traversal.NewRawTraversal("stringy"),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			root := traversal.NewRawTraversal(tC.input)
			require.Equal(t, tC.value, tC.get(root))
		})
	}
}
