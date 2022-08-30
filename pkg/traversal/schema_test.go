package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"
)

func TestNewSchemaTraversal(t *testing.T) {
	testCases := []struct {
		desc  string
		input map[string]interface{}
		cases StandardCases[map[string]interface{}]
	}{
		{
			desc: "object",
			input: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"extra": map[string]interface{}{
						"type": "edge",
					},
				},
			},
			cases: StandardCases[map[string]interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"extra": map[string]interface{}{
							"type": "edge",
						},
					},
				},
				Keys: map[string]traversal.Pointer[map[string]interface{}]{
					"extra": traversal.NewSchemaTraversal(map[string]interface{}{
						"type": "edge",
					}),
					"anything": nil,
				},
				Items: []traversal.Pointer[map[string]interface{}]{},
			},
		},
		{
			desc: "object_without_properties",
			input: map[string]interface{}{
				"type": "object",
			},
			cases: StandardCases[map[string]interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"type": "object",
				},
				Keys: map[string]traversal.Pointer[map[string]interface{}]{
					"anything": nil,
				},
				Items: []traversal.Pointer[map[string]interface{}]{},
			},
		},
		{
			desc: "array",
			input: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"id":   0,
						"type": "edge",
					},
					map[string]interface{}{
						"id":   1,
						"type": "edge",
					},
				},
			},
			cases: StandardCases[map[string]interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"type": "array",
					"items": []interface{}{
						map[string]interface{}{
							"id":   0,
							"type": "edge",
						},
						map[string]interface{}{
							"id":   1,
							"type": "edge",
						},
					},
				},
				Keys: map[string]traversal.Pointer[map[string]interface{}]{
					"anything": nil,
				},
				Items: []traversal.Pointer[map[string]interface{}]{
					traversal.NewSchemaTraversal(map[string]interface{}{
						"id":   0,
						"type": "edge",
					}),
					traversal.NewSchemaTraversal(map[string]interface{}{
						"id":   1,
						"type": "edge",
					}),
				},
			},
		},
		{
			desc: "array_no_items",
			input: map[string]interface{}{
				"type": "array",
			},
			cases: StandardCases[map[string]interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"type": "array",
				},
				Keys: map[string]traversal.Pointer[map[string]interface{}]{
					"anything": nil,
				},
				Items: []traversal.Pointer[map[string]interface{}]{},
			},
		},
		{
			desc: "edge",
			input: map[string]interface{}{
				"type": "edge",
			},
			cases: StandardCases[map[string]interface{}]{
				IsEdge: true, Value: map[string]interface{}{
					"type": "edge",
				},
				Keys: map[string]traversal.Pointer[map[string]interface{}]{
					"anything": nil,
				},
				Items: []traversal.Pointer[map[string]interface{}]{},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			PointerTestSuite(t, traversal.NewSchemaTraversal(tC.input), tC.cases)
		})
	}
}
