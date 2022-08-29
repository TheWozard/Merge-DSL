package traversal_test

import (
	"merge-dsl/pkg/traversal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexPointerItemsById(t *testing.T) {
	// All other cases are covered by tests including a parser.
	t.Run("nil_parser", func(t *testing.T) {
		index, extra := traversal.IndexPointerItemsById(traversal.NewRawTraversal([]interface{}{
			map[string]interface{}{"id": 1},
		}), nil)
		assert.Equal(t, map[interface{}]traversal.TraversalPointer[interface{}]{}, index)
		assert.Equal(t, []traversal.TraversalPointer[interface{}]{
			traversal.NewRawTraversal(map[string]interface{}{"id": 1}),
		}, extra)
	})
}

func TestPointerKeyIdIndexer(t *testing.T) {
	indexer := traversal.PointerKeyIdIndexer[interface{}]{Key: "id"}

	testCases := []struct {
		desc  string
		input traversal.TraversalPointer[interface{}]
		index map[interface{}]traversal.TraversalPointer[interface{}]
		extra []traversal.TraversalPointer[interface{}]
	}{
		{
			desc:  "nil_input",
			input: nil,
			index: map[interface{}]traversal.TraversalPointer[interface{}]{},
			extra: []traversal.TraversalPointer[interface{}]{},
		},
		{
			desc:  "non-array_input",
			input: traversal.NewRawTraversal("not an array :)"),
			index: map[interface{}]traversal.TraversalPointer[interface{}]{},
			extra: []traversal.TraversalPointer[interface{}]{},
		},
		{
			desc: "single_item_with_id",
			input: traversal.NewRawTraversal([]interface{}{
				map[string]interface{}{"id": 1},
			}),
			index: map[interface{}]traversal.TraversalPointer[interface{}]{
				1: traversal.NewRawTraversal(map[string]interface{}{"id": 1}),
			},
			extra: []traversal.TraversalPointer[interface{}]{},
		},
		{
			desc: "multiple_items_with_id",
			input: traversal.NewRawTraversal([]interface{}{
				map[string]interface{}{"id": 1},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
				map[string]interface{}{"id": 4},
			}),
			index: map[interface{}]traversal.TraversalPointer[interface{}]{
				1: traversal.NewRawTraversal(map[string]interface{}{"id": 1}),
				2: traversal.NewRawTraversal(map[string]interface{}{"id": 2}),
				3: traversal.NewRawTraversal(map[string]interface{}{"id": 3}),
				4: traversal.NewRawTraversal(map[string]interface{}{"id": 4}),
			},
			extra: []traversal.TraversalPointer[interface{}]{},
		},
		{
			desc: "single_item_without_id",
			input: traversal.NewRawTraversal([]interface{}{
				map[string]interface{}{},
			}),
			index: map[interface{}]traversal.TraversalPointer[interface{}]{},
			extra: []traversal.TraversalPointer[interface{}]{
				traversal.NewRawTraversal(map[string]interface{}{}),
			},
		},
		{
			desc: "multiple_items_without_id",
			input: traversal.NewRawTraversal([]interface{}{
				map[string]interface{}{"value": "A"},
				map[string]interface{}{"value": "B"},
				map[string]interface{}{"value": "C"},
				map[string]interface{}{"value": "D"},
			}),
			index: map[interface{}]traversal.TraversalPointer[interface{}]{},
			extra: []traversal.TraversalPointer[interface{}]{
				traversal.NewRawTraversal(map[string]interface{}{"value": "A"}),
				traversal.NewRawTraversal(map[string]interface{}{"value": "B"}),
				traversal.NewRawTraversal(map[string]interface{}{"value": "C"}),
				traversal.NewRawTraversal(map[string]interface{}{"value": "D"}),
			},
		},
		{
			desc: "multi_mixed",
			input: traversal.NewRawTraversal([]interface{}{
				map[string]interface{}{"value": "A"},
				map[string]interface{}{"id": 1},
				map[string]interface{}{"value": "B"},
				map[string]interface{}{"value": "C"},
				map[string]interface{}{"id": 2},
				map[string]interface{}{"id": 3},
				map[string]interface{}{"value": "D"},
				map[string]interface{}{"id": 4},
			}),
			index: map[interface{}]traversal.TraversalPointer[interface{}]{
				1: traversal.NewRawTraversal(map[string]interface{}{"id": 1}),
				2: traversal.NewRawTraversal(map[string]interface{}{"id": 2}),
				3: traversal.NewRawTraversal(map[string]interface{}{"id": 3}),
				4: traversal.NewRawTraversal(map[string]interface{}{"id": 4}),
			},
			extra: []traversal.TraversalPointer[interface{}]{
				traversal.NewRawTraversal(map[string]interface{}{"value": "A"}),
				traversal.NewRawTraversal(map[string]interface{}{"value": "B"}),
				traversal.NewRawTraversal(map[string]interface{}{"value": "C"}),
				traversal.NewRawTraversal(map[string]interface{}{"value": "D"}),
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			index, extra := indexer.Index(tC.input)
			assert.Equal(t, tC.index, index)
			assert.Equal(t, tC.extra, extra)
		})
	}
}
