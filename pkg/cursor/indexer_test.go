package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexCursorsById(t *testing.T) {
	t.Run("nil_input", func(t *testing.T) {
		index := map[interface{}][]cursor.RawCursor{}
		var cursors []cursor.RawCursor
		extra := cursor.PopulateIndexCursorsById(cursors, nil, index)
		assert.Equal(t, map[interface{}][]cursor.RawCursor{}, index)
		assert.Equal(t, []cursor.RawCursor{}, extra)
	})

	t.Run("nil_parser", func(t *testing.T) {
		index := map[interface{}][]cursor.RawCursor{}
		extra := cursor.PopulateIndexCursorsById([]cursor.RawCursor{
			cursor.NewRawCursor(map[string]interface{}{"id": 1}),
		}, nil, index)
		assert.Equal(t, map[interface{}][]cursor.RawCursor{}, index)
		assert.Equal(t, []cursor.RawCursor{
			cursor.NewRawCursor(map[string]interface{}{"id": 1}),
		}, extra)
	})
}

func TestRawCursorIdParser(t *testing.T) {
	testCases := []struct {
		desc  string
		input cursor.RawCursor
		index map[interface{}][]cursor.RawCursor
		extra []cursor.RawCursor
	}{
		{
			desc:  "non-array_input",
			input: cursor.NewRawCursor("not an array :)"),
			index: map[interface{}][]cursor.RawCursor{},
			extra: []cursor.RawCursor{},
		},
		{
			desc: "single_item_with_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{"id": 1},
			}),
			index: map[interface{}][]cursor.RawCursor{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1})},
			},
			extra: []cursor.RawCursor{},
		},
		{
			desc: "multiple_items_with_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
				map[string]interface{}{"id": 4},
			}),
			index: map[interface{}][]cursor.RawCursor{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1})},
				2: {cursor.NewRawCursor(map[string]interface{}{"id": 2})},
				3: {cursor.NewRawCursor(map[string]interface{}{"id": 3})},
				4: {cursor.NewRawCursor(map[string]interface{}{"id": 4})},
			},
			extra: []cursor.RawCursor{},
		},
		{
			desc: "single_item_without_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{},
			}),
			index: map[interface{}][]cursor.RawCursor{},
			extra: []cursor.RawCursor{
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
			index: map[interface{}][]cursor.RawCursor{},
			extra: []cursor.RawCursor{
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
			index: map[interface{}][]cursor.RawCursor{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1})},
				2: {cursor.NewRawCursor(map[string]interface{}{"id": 2})},
				3: {cursor.NewRawCursor(map[string]interface{}{"id": 3})},
				4: {cursor.NewRawCursor(map[string]interface{}{"id": 4})},
			},
			extra: []cursor.RawCursor{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "B"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "C"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "D"}),
			},
		},
		{
			desc: "grouped_id",
			input: cursor.NewRawCursor([]interface{}{
				map[string]interface{}{"id": 1, "value": "A"},
				map[string]interface{}{"id": 1, "value": "B"},
				map[string]interface{}{"id": 2, "value": "A"},
				map[string]interface{}{"id": 2, "value": "B"},
			}),
			index: map[interface{}][]cursor.RawCursor{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1, "value": "A"}), cursor.NewRawCursor(map[string]interface{}{"id": 1, "value": "B"})},
				2: {cursor.NewRawCursor(map[string]interface{}{"id": 2, "value": "A"}), cursor.NewRawCursor(map[string]interface{}{"id": 2, "value": "B"})},
			},
			extra: []cursor.RawCursor{},
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
