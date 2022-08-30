package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"
)

func TestNewRawTraversal(t *testing.T) {
	testCases := []struct {
		desc  string
		input interface{}
		cases StandardCases[interface{}]
	}{
		{
			desc:  "string",
			input: "data",
			cases: StandardCases[interface{}]{
				IsEdge: true, Value: "data",
				Keys: map[string]traversal.Pointer[interface{}]{
					"anything": nil,
				},
				Items: []traversal.Pointer[interface{}]{},
			},
		},
		{
			desc: "map",
			input: map[string]interface{}{
				"text": "data",
			},
			cases: StandardCases[interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"text": "data",
				},
				Keys: map[string]traversal.Pointer[interface{}]{
					"text":     traversal.NewRawTraversal("data"),
					"anything": nil,
				},
				Items: []traversal.Pointer[interface{}]{},
			},
		},
		{
			desc: "slice",
			input: []interface{}{
				map[string]interface{}{
					"id": 0,
				},
				"stringy",
			},
			cases: StandardCases[interface{}]{
				IsEdge: false, Value: []interface{}{
					map[string]interface{}{"id": 0},
					"stringy",
				},
				Keys: map[string]traversal.Pointer[interface{}]{
					"anything": nil,
				},
				Items: []traversal.Pointer[interface{}]{
					traversal.NewRawTraversal(map[string]interface{}{"id": 0}),
					traversal.NewRawTraversal("stringy"),
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			PointerTestSuite(t, traversal.NewRawTraversal(tC.input), tC.cases)
		})
	}
}
