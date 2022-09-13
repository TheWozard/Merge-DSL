package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"
)

func TestNewSchemaCursor(t *testing.T) {
	parent := cursor.NewSchemaCursor(map[string]interface{}{
		"data": "parent",
	})

	testCases := []struct {
		desc  string
		input cursor.Cursor[cursor.SchemaData]
		cases CursorCase[cursor.SchemaData]
	}{
		{
			desc: "object_missing_properties",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type": "object",
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: false,
				Value: map[string]interface{}{
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
			desc: "object_empty_properties",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: false,
				Value: map[string]interface{}{
					"type":       "object",
					"properties": map[string]interface{}{},
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
			desc: "object",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"extra": map[string]interface{}{
						"type": "edge",
					},
				},
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: true,
				Value: map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"extra": map[string]interface{}{
							"type": "edge",
						},
					},
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"extra": cursor.NewSchemaCursorFrom(map[string]interface{}{
						"type": "edge",
					}, cursor.NewSchemaCursorFrom(map[string]interface{}{
						"type": "object",
						"properties": map[string]interface{}{
							"extra": map[string]interface{}{
								"type": "edge",
							},
						},
					}, parent)),
					"anything": nil,
				},
				Keys:    []string{"extra"},
				Items:   []cursor.Cursor[map[string]interface{}]{},
				Default: nil,
			},
		},
		{
			desc: "array_no_items",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type": "array",
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: false,
				Value: map[string]interface{}{
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
			desc: "array_empty_items",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type":  "array",
				"items": []interface{}{},
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: false,
				Value: map[string]interface{}{
					"type":  "array",
					"items": []interface{}{},
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
			desc: "array_default",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type": "array",
				"default": map[string]interface{}{
					"type": "edge",
				},
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: true,
				Value: map[string]interface{}{
					"type": "array",
					"default": map[string]interface{}{
						"type": "edge",
					},
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"anything": nil,
				},
				Keys:  []string{},
				Items: []cursor.Cursor[map[string]interface{}]{},
				Default: cursor.NewSchemaCursorFrom(map[string]interface{}{
					"type": "edge",
				}, cursor.NewSchemaCursorFrom(map[string]interface{}{
					"type": "array",
					"default": map[string]interface{}{
						"type": "edge",
					},
				}, parent)),
			},
		},
		{
			desc: "array_items",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"id":   0,
						"type": "edge",
					},
				},
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: true,
				Value: map[string]interface{}{
					"type": "array",
					"items": []interface{}{
						map[string]interface{}{
							"id":   0,
							"type": "edge",
						},
					},
				},
				Key: map[string]cursor.Cursor[map[string]interface{}]{
					"anything": nil,
				},
				Keys: []string{},
				Items: []cursor.Cursor[map[string]interface{}]{
					cursor.NewSchemaCursorFrom(map[string]interface{}{
						"id":   0,
						"type": "edge",
					}, cursor.NewSchemaCursorFrom(map[string]interface{}{
						"type": "array",
						"items": []interface{}{
							map[string]interface{}{
								"id":   0,
								"type": "edge",
							},
						},
					}, parent)),
				},
			},
		},
		{
			desc: "edge",
			input: cursor.NewSchemaCursorFrom(map[string]interface{}{
				"type": "edge",
			}, parent),
			cases: CursorCase[map[string]interface{}]{
				Parent: parent, HasChildren: false,
				Value: map[string]interface{}{
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
			CursorTestSuite(t, tC.input, tC.cases)
		})
	}
}
