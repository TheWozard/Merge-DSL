package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexCursorsById(t *testing.T) {
	t.Run("nil_input", func(t *testing.T) {
		var cursors []cursor.RawCursor
		index, extra := cursor.IndexCursorsById(cursors, nil)
		assert.Equal(t, map[interface{}]cursor.Cursor[interface{}]{}, index)
		assert.Equal(t, []cursor.Cursor[interface{}]{}, extra)
	})

	t.Run("nil_parser", func(t *testing.T) {
		index, extra := cursor.IndexCursorsById([]cursor.RawCursor{
			cursor.NewRawCursor(map[string]interface{}{"id": 1}),
		}, nil)
		assert.Equal(t, map[interface{}]cursor.Cursor[interface{}]{}, index)
		assert.Equal(t, []cursor.Cursor[interface{}]{
			cursor.NewRawCursor(map[string]interface{}{"id": 1}),
		}, extra)
	})
}

func TestRawCursorIdParser(t *testing.T) {
	testCases := []struct {
		desc  string
		input cursor.Cursor[interface{}]
		index map[interface{}]cursor.Cursor[interface{}]
		extra []cursor.Cursor[interface{}]
	}{
		{
			desc:  "non-array_input",
			input: cursor.NewRawCursor("not an array :)"),
			index: map[interface{}]cursor.Cursor[interface{}]{},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "single_item_with_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{"id": 1},
			}),
			index: map[interface{}]cursor.Cursor[interface{}]{
				1: cursor.NewRawCursor(map[string]interface{}{"id": 1}),
			},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "multiple_items_with_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
				map[string]interface{}{"id": 4},
			}),
			index: map[interface{}]cursor.Cursor[interface{}]{
				1: cursor.NewRawCursor(map[string]interface{}{"id": 1}),
				2: cursor.NewRawCursor(map[string]interface{}{"id": 2}),
				3: cursor.NewRawCursor(map[string]interface{}{"id": 3}),
				4: cursor.NewRawCursor(map[string]interface{}{"id": 4}),
			},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "single_item_without_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{},
			}),
			index: map[interface{}]cursor.Cursor[interface{}]{},
			extra: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{}),
			},
		},
		{
			desc: "multiple_items_without_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{"value": "A"},
				map[string]interface{}{"value": "B"},
				map[string]interface{}{"value": "C"},
				map[string]interface{}{"value": "D"},
			}),
			index: map[interface{}]cursor.Cursor[interface{}]{},
			extra: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "B"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "C"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "D"}),
			},
		},
		{
			desc: "multi_mixed",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{"value": "A"},
				map[string]interface{}{"id": 1},
				map[string]interface{}{"value": "B"},
				map[string]interface{}{"value": "C"},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
				map[string]interface{}{"value": "D"},
				map[string]interface{}{"id": 4},
			}),
			index: map[interface{}]cursor.Cursor[interface{}]{
				1: cursor.NewRawCursor(map[string]interface{}{"id": 1}),
				2: cursor.NewRawCursor(map[string]interface{}{"id": 2}),
				3: cursor.NewRawCursor(map[string]interface{}{"id": 3}),
				4: cursor.NewRawCursor(map[string]interface{}{"id": 4}),
			},
			extra: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "B"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "C"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "D"}),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			index, extra := cursor.DefaultRawIndexer.Index(tC.input.GetItems())
			assert.Equal(t, tC.index, index)
			assert.Equal(t, tC.extra, extra)
		})
	}
}
