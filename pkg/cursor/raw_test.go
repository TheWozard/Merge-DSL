package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"
)

func TestNewRawCursor(t *testing.T) {
	testCases := []struct {
		desc  string
		input interface{}
		cases CursorCase[interface{}]
	}{
		{
			desc:  "string",
			input: "data",
			cases: CursorCase[interface{}]{
				IsEdge: true, Value: "data",
				Key: map[string]cursor.Cursor[interface{}]{
					"anything": nil,
				},
				Keys:    []string{},
				Items:   []cursor.Cursor[interface{}]{},
				Default: nil,
			},
		},
		{
			desc: "map",
			input: map[string]interface{}{
				"text": "data",
			},
			cases: CursorCase[interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"text": "data",
				},
				Key: map[string]cursor.Cursor[interface{}]{
					"text":     cursor.NewRawCursor("data"),
					"anything": nil,
				},
				Keys:    []string{"text"},
				Items:   []cursor.Cursor[interface{}]{},
				Default: nil,
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
			cases: CursorCase[interface{}]{
				IsEdge: false, Value: []interface{}{
					map[string]interface{}{"id": 0},
					"stringy",
				},
				Key: map[string]cursor.Cursor[interface{}]{
					"anything": nil,
				},
				Keys: []string{},
				Items: []cursor.Cursor[interface{}]{
					cursor.NewRawCursor(map[string]interface{}{"id": 0}),
					cursor.NewRawCursor("stringy"),
				},
				Default: nil,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			CursorTestSuite(t, cursor.NewRawCursor(tC.input), tC.cases)
		})
	}
}
