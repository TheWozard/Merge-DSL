package merge

import (
	"fmt"
	"merge-dsl/pkg/internal"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
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
			err:    fmt.Errorf("failed to locate type in definition"),
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
			desc: "leaf",
			input: map[string]interface{}{
				"type": "leaf",
			},
			output: &Definition{
				traversal: &edgeTraversal{},
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, err := Compile(tC.input)
			internal.ErrorsMatch(t, tC.err, err)
			assert.Equal(t, tC.output, output)
		})
	}
}
