package merge_test

import (
	"merge-dsl/pkg/cursor"
	. "merge-dsl/pkg/merge"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	testCases := []struct {
		desc   string
		def    map[string]interface{}
		docs   cursor.Set[interface{}]
		rules  cursor.Set[cursor.SchemaData]
		output interface{}
	}{
		// Edges
		{
			desc: "empty edge",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs:   cursor.Set[interface{}]{},
			rules:  cursor.Set[cursor.SchemaData]{},
			output: nil,
		},
		{
			desc: "single value edge",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs: cursor.Set[interface{}]{
				cursor.NewRawCursor("some data"),
			},
			rules:  cursor.Set[cursor.SchemaData]{},
			output: "some data",
		},
		{
			desc: "skips nil value edges",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs: cursor.Set[interface{}]{
				cursor.NewRawCursor(nil),
				cursor.NewRawCursor("low priority"),
				cursor.NewRawCursor("not used"),
			},
			rules:  cursor.Set[cursor.SchemaData]{},
			output: "low priority",
		},
		{
			desc: "default edge value",
			def: map[string]interface{}{
				"type":    "edge",
				"default": 0,
			},
			docs:   cursor.Set[interface{}]{},
			rules:  cursor.Set[cursor.SchemaData]{},
			output: 0,
		},
		{
			desc: "object empty",
			def: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"A": map[string]interface{}{
						"type": "edge",
					},
					"B": map[string]interface{}{
						"type": "edge",
					},
					"C": map[string]interface{}{
						"type": "edge",
					},
				},
			},
			docs:   cursor.Set[interface{}]{},
			rules:  cursor.Set[cursor.SchemaData]{},
			output: nil,
		},
		{
			desc: "object with values",
			def: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"A": map[string]interface{}{
						"type": "edge",
					},
					"B": map[string]interface{}{
						"type": "edge",
					},
					"C": map[string]interface{}{
						"type": "edge",
					},
				},
			},
			docs: cursor.Set[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{
					"A": 1,
				}),
				cursor.NewRawCursor(map[string]interface{}{
					"B": 2,
				}),
				cursor.NewRawCursor(map[string]interface{}{
					"C": 3,
				}),
			},
			rules: cursor.Set[cursor.SchemaData]{},
			output: map[string]interface{}{
				"A": 1,
				"B": 2,
				"C": 3,
			},
		},
		{
			desc: "array",
			def: map[string]interface{}{
				"type": "array",
				"default": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name": map[string]interface{}{
							"type": "edge",
						},
						"points": map[string]interface{}{
							"type":    "edge",
							"default": 0,
						},
					},
				},
				"items": []interface{}{
					map[string]interface{}{
						"id":   "favorite",
						"type": "object",
						"properties": map[string]interface{}{
							"name": map[string]interface{}{
								"type": "edge",
							},
							"points": map[string]interface{}{
								"type":    "edge",
								"default": 0,
							},
							"reason": map[string]interface{}{
								"type": "edge",
							},
						},
					},
				},
			},
			docs: cursor.Set[interface{}]{
				cursor.NewRawCursor([]interface{}{
					"skip",
					map[string]interface{}{
						"name": "A",
					},
				}),
				cursor.NewRawCursor([]interface{}{
					map[string]interface{}{
						"id":     "favorite",
						"name":   "B",
						"points": 10,
						"reason": "its cool",
					},
				}),
				cursor.NewRawCursor([]interface{}{
					map[string]interface{}{
						"name":   "C",
						"points": 1,
					},
					map[string]interface{}{
						"name": "D",
					},
				}),
			},
			rules: cursor.Set[cursor.SchemaData]{},
			output: []interface{}{
				map[string]interface{}{
					"points": 0, // TODO: This is awkward
				},
				map[string]interface{}{
					"name":   "A",
					"points": 0,
				},
				map[string]interface{}{
					"name":   "B",
					"points": 10,
					"reason": "its cool",
				},
				map[string]interface{}{
					"name":   "C",
					"points": 1,
				},
				map[string]interface{}{
					"name":   "D",
					"points": 0,
				},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			def, err := Compile(tC.def)
			require.Nil(t, err)
			output := def.Resolve(tC.docs, tC.rules)
			assert.Equal(t, tC.output, output)
		})
	}
}
