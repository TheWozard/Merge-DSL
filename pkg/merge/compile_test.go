package merge

import (
	"fmt"
	"merge-dsl/pkg/internal"
	"merge-dsl/pkg/reference"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	importer := reference.Resolver{
		reference.SchemaPrefix: (&reference.FileClient{Root: "../resources/schemas"}).Import,
	}
	compiler := Compiler{
		Importer:  importer,
		Validator: reference.NewSchemaValidator(importer),
	}

	testCases := []struct {
		desc   string
		input  map[string]interface{}
		output *Definition
		err    error
	}{
		{
			desc:   "empty input",
			input:  map[string]interface{}{},
			output: nil,
			err:    fmt.Errorf("cannot compile nil cursor"),
		},
		{
			desc: "object",
			input: map[string]interface{}{
				"type": "object",
			},
			output: &Definition{
				traversal: &objectTraversal{
					nodeTraversals: map[string]traversal{},
				},
			},
		},
		{
			desc: "array",
			input: map[string]interface{}{
				"type": "array",
			},
			output: &Definition{
				traversal: &arrayTraversal{
					idTraversals: map[interface{}]traversal{},
				},
			},
		},
		{
			desc: "edge",
			input: map[string]interface{}{
				"type": "edge",
			},
			output: &Definition{
				traversal: &edgeTraversal{},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, err := compiler.Compile(tC.input)
			internal.ErrorsMatch(t, tC.err, err)
			assert.Equal(t, tC.output, output)
		})
	}
}
