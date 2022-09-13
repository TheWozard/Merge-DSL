package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"
)

func TestNewRawCursor(t *testing.T) {
	parent := cursor.NewRawCursor("parent")

	testCases := []struct {
		desc  string
		input cursor.Cursor[interface{}]
		cases CursorCase[interface{}]
	}{
		{
			desc:  "edge",
			input: cursor.NewRawCursorFrom("edge", parent),
			cases: CursorCase[interface{}]{
				Parent: parent, HasChildren: false,
				Value: "edge",
				Key: map[string]cursor.Cursor[interface{}]{
					"anything": nil,
				},
				Keys:    []string{},
				Items:   []cursor.Cursor[interface{}]{},
				Default: nil,
			},
		},
		{
			desc:  "map_empty",
			input: cursor.NewRawCursorFrom(map[string]interface{}{}, parent),
			cases: CursorCase[interface{}]{
				Parent: parent, HasChildren: false,
				Value: map[string]interface{}{},
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
			input: cursor.NewRawCursorFrom(map[string]interface{}{
				"text": "data",
			}, parent),
			cases: CursorCase[interface{}]{
				Parent: parent, HasChildren: true,
				Value: map[string]interface{}{
					"text": "data",
				},
				Key: map[string]cursor.Cursor[interface{}]{
					"text": cursor.NewRawCursorFrom("data", cursor.NewRawCursorFrom(map[string]interface{}{
						"text": "data",
					}, parent)),
					"anything": nil,
				},
				Keys:    []string{"text"},
				Items:   []cursor.Cursor[interface{}]{},
				Default: nil,
			},
		},
		{
			desc:  "slice_empty",
			input: cursor.NewRawCursorFrom([]interface{}{}, parent),
			cases: CursorCase[interface{}]{
				Parent: parent, HasChildren: false,
				Value: []interface{}{},
				Key: map[string]cursor.Cursor[interface{}]{
					"anything": nil,
				},
				Keys:    []string{},
				Items:   []cursor.Cursor[interface{}]{},
				Default: nil,
			},
		},
		{
			desc: "slice",
			input: cursor.NewRawCursorFrom([]interface{}{
				map[string]interface{}{"id": 0},
				"stringy",
			}, parent),
			cases: CursorCase[interface{}]{
				Parent: parent, HasChildren: true,
				Value: []interface{}{
					map[string]interface{}{"id": 0},
					"stringy",
				},
				Key: map[string]cursor.Cursor[interface{}]{
					"anything": nil,
				},
				Keys: []string{},
				Items: []cursor.Cursor[interface{}]{
					cursor.NewRawCursorFrom(map[string]interface{}{"id": 0},
						cursor.NewRawCursorFrom([]interface{}{
							map[string]interface{}{"id": 0},
							"stringy",
						}, parent),
					),
					cursor.NewRawCursorFrom("stringy",
						cursor.NewRawCursorFrom([]interface{}{
							map[string]interface{}{"id": 0},
							"stringy",
						}, parent),
					),
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
