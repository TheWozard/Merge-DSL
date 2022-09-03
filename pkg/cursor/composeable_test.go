package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"
)

func TestComposedCursor(t *testing.T) {
	type TestData struct {
		Value string
	}

	testCases := []struct {
		desc  string
		input cursor.Cursor[TestData]
		cases CursorCase[TestData]
	}{
		{
			desc:  "edge",
			input: cursor.EdgePointer[TestData]{Data: TestData{Value: "edge"}},
			cases: CursorCase[TestData]{
				IsEdge: true, Value: TestData{Value: "edge"},
				Key: map[string]cursor.Cursor[TestData]{
					"anything": nil,
				},
				Keys:    []string{},
				Items:   []cursor.Cursor[TestData]{},
				Default: nil,
			},
		},
		{
			desc: "map",
			input: cursor.MapPointer[TestData]{Data: TestData{Value: "object"}, Nodes: map[string]cursor.Cursor[TestData]{
				"extra": cursor.EdgePointer[TestData]{Data: TestData{Value: "data"}},
			}},
			cases: CursorCase[TestData]{
				IsEdge: false, Value: TestData{Value: "object"},
				Key: map[string]cursor.Cursor[TestData]{
					"extra": cursor.EdgePointer[TestData]{Data: TestData{Value: "data"}},
					"other": nil,
				},
				Keys:    []string{"extra"},
				Items:   []cursor.Cursor[TestData]{},
				Default: nil,
			},
		},
		{
			desc: "slice",
			input: cursor.SlicePointer[TestData]{Data: TestData{Value: "slice"}, Items: []cursor.Cursor[TestData]{
				cursor.EdgePointer[TestData]{Data: TestData{Value: "details"}},
			}},
			cases: CursorCase[TestData]{
				IsEdge: false, Value: TestData{Value: "slice"},
				Key: map[string]cursor.Cursor[TestData]{
					"anything": nil,
				},
				Keys: []string{},
				Items: []cursor.Cursor[TestData]{
					cursor.EdgePointer[TestData]{Data: TestData{Value: "details"}},
				},
				Default: nil,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			CursorTestSuite(t, tC.input, tC.cases)
		})
	}
}
