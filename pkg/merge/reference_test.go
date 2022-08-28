package merge_test

import (
	"fmt"
	"merge-dsl/pkg/internal"
	"merge-dsl/pkg/merge"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportReference(t *testing.T) {
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
			result, err := merge.ImportReference[interface{}](tC.reference)
			internal.ErrorsMatch(t, tC.err, err)
			assert.Equal(t, tC.result, result)
		})
	}
}

func TestImportToCustomStruct(t *testing.T) {
	type TestStruct struct {
		Data string `json:"data"`
	}
	result, err := merge.ImportReference[TestStruct](`json:{"data": "expected"}`)
	assert.Nil(t, err)
	assert.Equal(t, TestStruct{Data: "expected"}, result)
}

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
			result: map[string]interface{}{
				"example": "value",
			},
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
			result, err := merge.ImportValidatedReference[interface{}](tC.reference, tC.schema)
			internal.ErrorsMatch(t, tC.err, err)
			assert.Equal(t, tC.result, result)
		})
	}
}
