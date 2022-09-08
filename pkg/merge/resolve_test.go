package merge_test

import (
	"merge-dsl/pkg/cursor"
	. "merge-dsl/pkg/merge"
	"merge-dsl/pkg/reference"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolve(t *testing.T) {
	importer := reference.Resolver{
		reference.SchemaPrefix: (&reference.FileClient{Root: "../resources/schemas"}).Import,
	}
	compiler := Compiler{
		Importer:  importer,
		Validator: reference.NewSchemaValidator(importer),
	}

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
			docs:   cursor.CursorSet[cursor.RawData]{},
			rules:  cursor.CursorSet[cursor.SchemaData]{},
			output: nil,
		},
		{
			desc: "single value edge",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs: cursor.CursorSet[cursor.RawData]{
				cursor.NewRawCursor("some data"),
			},
			rules:  cursor.CursorSet[cursor.SchemaData]{},
			output: "some data",
		},
		{
			desc: "skips nil value edges",
			def: map[string]interface{}{
				"type": "edge",
			},
			docs: cursor.CursorSet[cursor.RawData]{
				cursor.NewRawCursor(nil),
				cursor.NewRawCursor("low priority"),
				cursor.NewRawCursor("not used"),
			},
			rules:  cursor.CursorSet[cursor.SchemaData]{},
			output: "low priority",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			def, err := compiler.Compile(tC.def)
			require.Nil(t, err)
			output, err := def.Resolve(tC.docs, tC.rules)
			assert.Nil(t, err)
			assert.Equal(t, tC.output, output)
		})
	}
}
