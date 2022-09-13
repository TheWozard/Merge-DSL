package cursor_test

import (
	"merge-dsl/pkg/cursor"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRawCursorIdParser(t *testing.T) {
	grouper := cursor.DefaultRawGrouper

	testCases := []struct {
		desc  string
		input []cursor.Cursor[interface{}]
		group [][]cursor.Cursor[interface{}]
		index map[interface{}][]cursor.Cursor[interface{}]
		extra []cursor.Cursor[interface{}]
	}{
		{
			desc:  "empty_list",
			input: []cursor.Cursor[interface{}]{},
			group: [][]cursor.Cursor[interface{}]{},
			index: map[interface{}][]cursor.Cursor[interface{}]{},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "single_item_with_id",
			input: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"id": 1}),
			},
			group: [][]cursor.Cursor[interface{}]{
				{cursor.NewRawCursor(map[string]interface{}{"id": 1})},
			},
			index: map[interface{}][]cursor.Cursor[interface{}]{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1})},
			},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "multiple_items_with_different_id",
			input: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"id": 1}),
				cursor.NewRawCursor(map[string]interface{}{"id": 2}),
			},
			group: [][]cursor.Cursor[interface{}]{
				{cursor.NewRawCursor(map[string]interface{}{"id": 1})},
				{cursor.NewRawCursor(map[string]interface{}{"id": 2})},
			},
			index: map[interface{}][]cursor.Cursor[interface{}]{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1})},
				2: {cursor.NewRawCursor(map[string]interface{}{"id": 2})},
			},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "multiple_items_with_different_id_order_is_preserved",
			input: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"id": 2}),
				cursor.NewRawCursor(map[string]interface{}{"id": 1}),
			},
			group: [][]cursor.Cursor[interface{}]{
				{cursor.NewRawCursor(map[string]interface{}{"id": 2})},
				{cursor.NewRawCursor(map[string]interface{}{"id": 1})},
			},
			index: map[interface{}][]cursor.Cursor[interface{}]{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1})},
				2: {cursor.NewRawCursor(map[string]interface{}{"id": 2})},
			},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "multiple_items_with_same_id",
			input: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"id": 1}),
				cursor.NewRawCursor(map[string]interface{}{"id": 1}),
			},
			group: [][]cursor.Cursor[interface{}]{
				{cursor.NewRawCursor(map[string]interface{}{"id": 1}),
					cursor.NewRawCursor(map[string]interface{}{"id": 1})},
			},
			index: map[interface{}][]cursor.Cursor[interface{}]{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1}),
					cursor.NewRawCursor(map[string]interface{}{"id": 1})},
			},
			extra: []cursor.Cursor[interface{}]{},
		},
		{
			desc: "single_item_without_id",
			input: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
			},
			group: [][]cursor.Cursor[interface{}]{
				{cursor.NewRawCursor(map[string]interface{}{"value": "A"})},
			},
			index: map[interface{}][]cursor.Cursor[interface{}]{},
			extra: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
			},
		},
		{
			desc: "multiple_items_without_id",
			input: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "B"}),
			},
			group: [][]cursor.Cursor[interface{}]{
				{cursor.NewRawCursor(map[string]interface{}{"value": "A"})},
				{cursor.NewRawCursor(map[string]interface{}{"value": "B"})},
			},
			index: map[interface{}][]cursor.Cursor[interface{}]{},
			extra: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "B"}),
			},
		},
		{
			desc: "multi_mixed",
			input: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"id": 2}),
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
				cursor.NewRawCursor(map[string]interface{}{"id": 1}),
				cursor.NewRawCursor(map[string]interface{}{"id": 1, "value": "B"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "C"}),
				cursor.NewRawCursor(map[string]interface{}{"id": 2, "value": "D"}),
			},
			group: [][]cursor.Cursor[interface{}]{
				{cursor.NewRawCursor(map[string]interface{}{"id": 2}),
					cursor.NewRawCursor(map[string]interface{}{"id": 2, "value": "D"})},
				{cursor.NewRawCursor(map[string]interface{}{"value": "A"})},
				{cursor.NewRawCursor(map[string]interface{}{"id": 1}),
					cursor.NewRawCursor(map[string]interface{}{"id": 1, "value": "B"})},
				{cursor.NewRawCursor(map[string]interface{}{"value": "C"})},
			},
			index: map[interface{}][]cursor.Cursor[interface{}]{
				1: {cursor.NewRawCursor(map[string]interface{}{"id": 1}),
					cursor.NewRawCursor(map[string]interface{}{"id": 1, "value": "B"})},
				2: {cursor.NewRawCursor(map[string]interface{}{"id": 2}),
					cursor.NewRawCursor(map[string]interface{}{"id": 2, "value": "D"})},
			},
			extra: []cursor.Cursor[interface{}]{
				cursor.NewRawCursor(map[string]interface{}{"value": "A"}),
				cursor.NewRawCursor(map[string]interface{}{"value": "C"}),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			group := cursor.GroupCursors(tC.input, grouper)
			assert.Equal(t, tC.group, group)
			index, extra := cursor.IndexCursors(tC.input, grouper)
			assert.Equal(t, tC.index, index)
			assert.Equal(t, tC.extra, extra)
		})
	}
}
