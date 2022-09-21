package merge_test

import (
	"merge-dsl/pkg/merge"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCalculationFactory(t *testing.T) {
	testCases := []struct {
		desc   string
		name   string
		input  map[string]interface{}
		output merge.Operation
	}{
		{
			desc: "add",
			name: "add",
			input: map[string]interface{}{
				"keys": []interface{}{"a", "b", "c"},
			},
			output: &merge.AddOperation{Keys: []string{"a", "b", "c"}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			require.Equal(t, tC.output, merge.GetOperation(tC.name, tC.input))
		})
	}
}
