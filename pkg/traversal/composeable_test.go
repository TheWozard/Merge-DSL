package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestComposedTraversal(t *testing.T) {
	type TestData struct {
		Value string
	}

	edgeExample := traversal.EdgePointer[TestData]{Data: TestData{Value: "edge"}}
	objectExample := traversal.MapPointer[TestData]{Data: TestData{Value: "object"}, Nodes: map[string]traversal.TraversalPointer[TestData]{
		"extra": traversal.EdgePointer[TestData]{Data: TestData{Value: "data"}},
	}}
	arrayExample := traversal.SlicePointer[TestData]{Data: TestData{Value: "slice"}, Items: []traversal.TraversalPointer[TestData]{
		traversal.EdgePointer[TestData]{Data: TestData{Value: "details"}},
	}}

	testCases := []struct {
		desc  string
		input traversal.TraversalPointer[TestData]
		get   func(pointer traversal.TraversalPointer[TestData]) interface{}
		value interface{}
	}{
		// Edge
		{
			desc:  "string_value",
			input: edgeExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.Value()
			},
			value: TestData{Value: "edge"},
		},
		{
			desc:  "string_isEdge",
			input: edgeExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.IsEdge()
			},
			value: true,
		},
		{
			desc:  "string_getKey",
			input: edgeExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.GetKey("anything")
			},
			value: nil,
		},
		{
			desc:  "string_getItems",
			input: edgeExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[TestData]{},
		},

		// Object
		{
			desc:  "object_value",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.Value()
			},
			value: TestData{Value: "object"},
		},
		{
			desc:  "object_isEdge",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.IsEdge()
			},
			value: false,
		},
		{
			desc:  "object_getKey_bad",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.GetKey("bad")
			},
			value: nil,
		},
		{
			desc:  "object_getKey_good",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.GetKey("extra")
			},
			value: traversal.EdgePointer[TestData]{Data: TestData{Value: "data"}},
		},
		{
			desc:  "object_getItems",
			input: objectExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[TestData]{},
		},

		// Array
		{
			desc:  "array_value",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.Value()
			},
			value: TestData{Value: "slice"},
		},
		{
			desc:  "array_isEdge",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.IsEdge()
			},
			value: false,
		},
		{
			desc:  "array_getKey",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.GetKey("anything")
			},
			value: nil,
		},
		{
			desc:  "array_getItems",
			input: arrayExample,
			get: func(pointer traversal.TraversalPointer[TestData]) interface{} {
				return pointer.GetItems()
			},
			value: []traversal.TraversalPointer[TestData]{
				traversal.EdgePointer[TestData]{Data: TestData{Value: "details"}},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.value, tC.get(tC.input))
		})
	}
}
