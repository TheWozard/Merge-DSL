package reference_test

import (
	"fmt"
	"merge-dsl/pkg/internal"
	"merge-dsl/pkg/reference"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	defaultTestValidator = reference.NewSchemaValidator(defaultTestImporter)
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
			imported, err := reference.Import[interface{}](defaultTestImporter, tC.reference)
			require.Nil(t, err)
			err = defaultTestValidator.IsValidByReference(imported.Data, tC.schema)
			internal.ErrorsMatch(t, tC.err, err)
		})
	}
}
