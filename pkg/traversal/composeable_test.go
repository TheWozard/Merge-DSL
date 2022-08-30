package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"
)

func TestComposedTraversal(t *testing.T) {
	type TestData struct {
		Value string
	}

	testCases := []struct {
		desc  string
		input traversal.Pointer[TestData]
		cases StandardCases[TestData]
	}{
		{
			desc:  "edge",
			input: traversal.EdgePointer[TestData]{Data: TestData{Value: "edge"}},
			cases: StandardCases[TestData]{
				IsEdge: true, Value: TestData{Value: "edge"},
				Keys: map[string]traversal.Pointer[TestData]{
					"anything": nil,
				},
				Items: []traversal.Pointer[TestData]{},
			},
		},
		{
			desc: "map",
			input: traversal.MapPointer[TestData]{Data: TestData{Value: "object"}, Nodes: map[string]traversal.Pointer[TestData]{
				"extra": traversal.EdgePointer[TestData]{Data: TestData{Value: "data"}},
			}},
			cases: StandardCases[TestData]{
				IsEdge: false, Value: TestData{Value: "object"},
				Keys: map[string]traversal.Pointer[TestData]{
					"extra": traversal.EdgePointer[TestData]{Data: TestData{Value: "data"}},
					"other": nil,
				},
				Items: []traversal.Pointer[TestData]{},
			},
		},
		{
			desc: "slice",
			input: traversal.SlicePointer[TestData]{Data: TestData{Value: "slice"}, Items: []traversal.Pointer[TestData]{
				traversal.EdgePointer[TestData]{Data: TestData{Value: "details"}},
			}},
			cases: StandardCases[TestData]{
				IsEdge: false, Value: TestData{Value: "slice"},
				Keys: map[string]traversal.Pointer[TestData]{
					"anything": nil,
				},
				Items: []traversal.Pointer[TestData]{
					traversal.EdgePointer[TestData]{Data: TestData{Value: "details"}},
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			PointerTestSuite(t, tC.input, tC.cases)
		})
	}
}
