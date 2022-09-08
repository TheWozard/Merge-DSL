package reference_test

import (
	"merge-dsl/pkg/internal"
	"merge-dsl/pkg/reference"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	defaultTestImporter = reference.Resolver{
		reference.JSONPrefix: (&reference.EmbeddedClient{Format: reference.JSONFormat}).Import,
		reference.YAMLPrefix: (&reference.EmbeddedClient{Format: reference.YAMLFormat}).Import,
	}
)

func TestImporter(t *testing.T) {

	testCases := []struct {
		desc      string
		reference string
		result    interface{}
		err       error
	}{
		{
			desc:      "JSON_Object",
			reference: `json:{"example":"value"}`,
			result: map[string]interface{}{
				"example": "value",
			},
		},
		{
			desc:      "JSON_Array",
			reference: `json:["example", "value"]`,
			result: []interface{}{
				"example", "value",
			},
		},
		{
			desc:      "YAML_Object",
			reference: `yaml:example: value`,
			result: map[string]interface{}{
				"example": "value",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result, err := reference.Import[interface{}](defaultTestImporter, tC.reference)
			internal.ErrorsMatch(t, tC.err, err)
			assert.Equal(t, tC.result, result.Data)
		})
	}
}

func TestImportToCustomStruct(t *testing.T) {
	type TestStruct struct {
		Data string `json:"data"`
	}
	result, err := reference.Import[TestStruct](defaultTestImporter, `json:{"data": "expected"}`)
	assert.Nil(t, err)
	assert.Equal(t, TestStruct{Data: "expected"}, result.Data)
}
