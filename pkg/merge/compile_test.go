package merge_test

import (
	"merge-dsl/pkg/internal"
	"merge-dsl/pkg/merge"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompile(t *testing.T) {
	testCases := []struct {
		desc   string
		input  interface{}
		output *merge.Definition
		err    error
	}{
		{
			desc:   "",
			input:  nil,
			output: &merge.Definition{},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			output, err := merge.Compile(tC.input)
			internal.ErrorsMatch(t, tC.err, err)
			assert.Equal(t, tC.output, output)
		})
	}
}
