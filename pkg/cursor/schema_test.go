package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"
)

func TestNewSchemaCursor(t *testing.T) {
	testCases := []struct {
		desc  string
		input map[string]interface{}
		cases CursorCase[map[string]interface{}]
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
			cases: CursorCase[map[string]interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"extra": map[string]interface{}{
							"type": "edge",
						},
					},
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"extra": cursor.NewSchemaCursor(map[string]interface{}{
						"type": "edge",
					}),
					"anything": nil,
				},
				Keys:    []string{"extra"},
				Items:   []cursor.Cursor[map[string]interface{}]{},
				Default: nil,
			},
		},
		{
			desc: "object_without_properties",
			input: map[string]interface{}{
				"type": "object",
			},
			cases: CursorCase[map[string]interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"type": "object",
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"anything": nil,
				},
				Keys:    []string{},
				Items:   []cursor.Cursor[map[string]interface{}]{},
				Default: nil,
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
				"default": map[string]interface{}{
					"type": "edge",
				},
			},
			cases: CursorCase[map[string]interface{}]{
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
					"default": map[string]interface{}{
						"type": "edge",
					},
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"anything": nil,
				},
				Keys: []string{},
				Items: []cursor.Cursor[map[string]interface{}]{
					cursor.NewSchemaCursor(map[string]interface{}{
						"id":   0,
						"type": "edge",
					}),
					cursor.NewSchemaCursor(map[string]interface{}{
						"id":   1,
						"type": "edge",
					}),
				},
				Default: cursor.NewSchemaCursor(map[string]interface{}{
					"type": "edge",
				}),
			},
		},
		{
			desc: "array_no_items",
			input: map[string]interface{}{
				"type": "array",
			},
			cases: CursorCase[map[string]interface{}]{
				IsEdge: false, Value: map[string]interface{}{
					"type": "array",
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"anything": nil,
				},
				Keys:    []string{},
				Items:   []cursor.Cursor[map[string]interface{}]{},
				Default: nil,
			},
		},
		{
			desc: "edge",
			input: map[string]interface{}{
				"type": "edge",
			},
			cases: CursorCase[map[string]interface{}]{
				IsEdge: true, Value: map[string]interface{}{
					"type": "edge",
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"anything": nil,
				},
				Keys:    []string{},
				Items:   []cursor.Cursor[map[string]interface{}]{},
				Default: nil,
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			CursorTestSuite(t, cursor.NewSchemaCursor(tC.input), tC.cases)
		})
	}
}
