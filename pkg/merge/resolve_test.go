package merge_test

import (
	"merge-dsl/pkg/cursor"
	"merge-dsl/pkg/cursor/validator"
	. "merge-dsl/pkg/merge"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	testCases := []struct {
		desc   string
		def    map[string]interface{}
		docs   cursor.CursorSet[cursor.RawData]
		rules  cursor.CursorSet[cursor.SchemaData]
		output interface{}
	}{
		{
			desc: "empty edge",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs:   cursor.CursorSet[cursor.RawData]{Set: []cursor.RawCursor{}},
			rules:  cursor.CursorSet[cursor.SchemaData]{Set: []cursor.SchemaCursor{}},
			output: nil,
		},
		{
			desc: "single value edge",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs: cursor.CursorSet[cursor.RawData]{Set: []cursor.RawCursor{
				cursor.NewRawCursor("some data"),
			}, Validator: validator.NonNil},
			rules:  cursor.CursorSet[cursor.SchemaData]{Set: []cursor.SchemaCursor{}},
			output: "some data",
		},
		{
			desc: "skips nil value edges",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs: cursor.CursorSet[cursor.RawData]{Set: []cursor.RawCursor{
				cursor.NewRawCursor(nil),
				cursor.NewRawCursor("low priority"),
				cursor.NewRawCursor("not used"),
			}, Validator: validator.NonNil},
			rules:  cursor.CursorSet[cursor.SchemaData]{Set: []cursor.SchemaCursor{}},
			output: "low priority",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			def, err := Compile(tC.def)
			require.Nil(t, err)
			output, err := def.Resolve(tC.docs, tC.rules)
			assert.Nil(t, err)
			assert.Equal(t, tC.output, output)
		})
	}
}
