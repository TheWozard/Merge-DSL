package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"
)

func TestSafeTraversal(t *testing.T) {
	testCases := []struct {
		desc  string
		input traversal.Pointer[interface{}]
		cases StandardCases[interface{}]
	}{
		// nil
		{
			desc:  "nil",
			input: nil,
			cases: StandardCases[interface{}]{
				IsEdge: true, Value: nil,
				Keys: map[string]traversal.Pointer[interface{}]{
					"anything": traversal.Safe[interface{}](nil),
				},
				Items: []traversal.Pointer[interface{}]{},
			},
		},
		{
			desc: "object",
			input: traversal.NewRawTraversal(map[string]interface{}{
				"example": "data",
			}),
			cases: StandardCases[interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"example": "data",
				},
				Keys: map[string]traversal.Pointer[interface{}]{
					"example":  traversal.Safe(traversal.NewRawTraversal("data")),
					"anything": traversal.Safe[interface{}](nil),
				},
				Items: []traversal.Pointer[interface{}]{},
			},
		},
		{
			desc:  "array",
			input: traversal.NewRawTraversal([]interface{}{"example", "data"}),
			cases: StandardCases[interface{}]{
				IsEdge: false, Value: []interface{}{"example", "data"},
				Keys: map[string]traversal.Pointer[interface{}]{
					"anything": traversal.Safe[interface{}](nil),
				},
				Items: []traversal.Pointer[interface{}]{
					traversal.Safe(traversal.NewRawTraversal("example")),
					traversal.Safe(traversal.NewRawTraversal("data")),
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			PointerTestSuite(t, traversal.Safe(tC.input), tC.cases)
		})
	}
}
