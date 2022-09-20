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
		output *Traversal
		err    error
	}{
		{
			desc:   "empty input",
			input:  map[string]interface{}{},
			output: nil,
			err:    fmt.Errorf("failed to locate type in definition"),
		},
		{
			desc: "unknown type",
			input: map[string]interface{}{
				"type": "bad",
			},
			output: nil,
			err:    fmt.Errorf("unknown compile type 'bad'"),
		},

		// Positive
		{
			desc: "object",
			input: map[string]interface{}{
				"type": "object",
			},
			output: &Traversal{
				step: &objectStep{
					nodeSteps: map[string]step{},
				},
			},
		},
		{
			desc: "object properties",
			input: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"a": map[string]interface{}{
						"type": "edge",
					},
				},
			},
			output: &Traversal{
				step: &objectStep{
					nodeSteps: map[string]step{
						"a": &edgeStep{},
					},
				},
			},
		},
		{
			desc: "array",
			input: map[string]interface{}{
				"type": "array",
			},
			output: &Traversal{
				step: &arrayStep{
					idStep: map[interface{}]step{},
				},
			},
		},
		{
			desc: "array default",
			input: map[string]interface{}{
				"type": "array",
				"default": map[string]interface{}{
					"type": "edge",
				},
			},
			output: &Traversal{
				step: &arrayStep{
					defaultStep: &edgeStep{},
					idStep:      map[interface{}]step{},
				},
			},
		},
		{
			desc: "array id",
			input: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"id":   "example",
						"type": "edge",
					},
				},
			},
			output: &Traversal{
				step: &arrayStep{
					idStep: map[interface{}]step{
						"example": &edgeStep{},
					},
				},
			},
		},
		{
			desc: "edge",
			input: map[string]interface{}{
				"type": "edge",
			},
			output: &Traversal{
				step: &edgeStep{},
			},
		},

		// Negative
		{
			desc: "object bad sub object",
			input: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"invalid": map[string]interface{}{
						"type": "bad",
					},
				},
			},
			output: nil,
			err:    fmt.Errorf("failed to compile node 'invalid': unknown compile type 'bad'"),
		},
		{
			desc: "array default bad",
			input: map[string]interface{}{
				"type": "array",
				"default": map[string]interface{}{
					"type": "bad",
				},
			},
			output: nil,
			err:    fmt.Errorf("failed to compile default: unknown compile type 'bad'"),
		},
		{
			desc: "array no id item",
			input: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"type": "edge",
					},
				},
			},
			output: nil,
			err:    fmt.Errorf("unexpected non-id node during array compile, all items are expected to contain an id"),
		},
		{
			desc: "array bad item",
			input: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"id":   "example",
						"type": "bad",
					},
				},
			},
			output: nil,
			err:    fmt.Errorf("failed to compile id traversal: unknown compile type 'bad'"),
		},
		{
			desc: "array multiple instances",
			input: map[string]interface{}{
				"type": "array",
				"items": []interface{}{
					map[string]interface{}{
						"id":   "example",
						"type": "edge",
					},
					map[string]interface{}{
						"id":   "example",
						"type": "edge",
					},
				},
			},
			output: nil,
			err:    fmt.Errorf("found 2 instances of the id 'example'"),
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
