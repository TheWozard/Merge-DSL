package pkg

import (
	"merge-dsl/pkg/merge"
	"merge-dsl/pkg/reference"
	"merge-dsl/pkg/resolution"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestExamples(t *testing.T) {
	importer := reference.NewDefaultResolver(
		"../examples", "./resources/schemas", time.Second*5, "example-test",
	)
	compiler := merge.Compiler{
		Importer:  importer,
		Validator: reference.NewSchemaValidator(importer),
	}

	testCases := []struct {
		description    string
		transform_ref  string
		documents_refs []string
		output         interface{}
	}{
		{
			description:    "empty_edge",
			transform_ref:  "file://transform/edge.yaml",
			documents_refs: []string{},
			output:         nil,
		},
		{
			description:   "single_value_edge",
			transform_ref: "file://transform/edge.yaml",
			documents_refs: []string{
				"file://documents/raw_string.yaml",
			},
			output: "raw",
		},

		{
			description:   "business",
			transform_ref: "file://transform/business.yaml",
			documents_refs: []string{
				"file://documents/business_review.yaml",
				"file://documents/business_details.yaml",
			},
			output: map[string]interface{}{
				"business_id": 1,
				"name":        "Greendale Community College",
			},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			def, err := compiler.CompileReference(tC.transform_ref)
			require.Nil(t, err)
			docs, rules, err := resolution.Resolve(tC.documents_refs, compiler.Importer)
			require.Nil(t, err)
			output, err := def.Resolve(docs, rules)
			require.Nil(t, err)
			require.Equal(t, tC.output, output)
		})
	}
}
