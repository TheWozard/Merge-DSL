package merge_test

import (
	"fmt"
	"merge-dsl/pkg/internal"
	"merge-dsl/pkg/merge"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImportValidatedReference(t *testing.T) {
	testCases := []struct {
		desc      string
		reference string
		schema    string
		result    interface{}
		err       error
	}{
		{
			desc:      "Object_Positive",
			reference: `json:{"example":"value"}`,
			schema:    `json:{"type":"object"}`,
		},
		{
			desc:      "Object_Failure",
			reference: `json:{"example":"value"}`,
			schema:    `json:{"type":"array"}`,
			err:       fmt.Errorf("failed to validate document: Invalid type. Expected: array, given: object"),
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data, err := merge.ImportReference[interface{}](tC.reference)
			require.Nil(t, err)
			err = merge.IsValidByReference(data, tC.schema)
			internal.ErrorsMatch(t, tC.err, err)
		})
	}
}