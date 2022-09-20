package pkg

import (
	"merge-dsl/pkg/reference"
	"merge-dsl/pkg/resolution"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	importer = reference.NewDefaultResolver(
		"../examples", "../resources/schemas", time.Second*5, "example-test",
	)
	compiler = reference.Compiler{
		Importer:  importer,
		Validator: reference.NewSchemaValidator(importer),
	}

	testCases = []struct {
		description   string
		transformRef  string
		documentsRefs []string
		output        interface{}
	}{
		{
			description:   "empty_edge",
			transformRef:  "file://transform/edge.yaml",
			documentsRefs: []string{},
			output:        nil,
		},
		{
			description:  "single_value_edge",
			transformRef: "file://transform/edge.yaml",
			documentsRefs: []string{
				"file://documents/raw_string.yaml",
			},
			output: "raw",
		},

		{
			description:  "business",
			transformRef: "file://transform/business.yaml",
			documentsRefs: []string{
				"file://documents/business_review.yaml",
				"file://documents/business_details.yaml",
				"file://documents/business_products.yaml",
			},
			output: map[string]interface{}{
				"business_id": 1,
				"name":        "Greendale Community College",
				"address": map[string]interface{}{
					"city":   "Greendale",
					"county": "Greendale County",
					"state":  "Colorado",
				},
				"reviews": map[string]interface{}{
					"count":     3000,
					"negatives": 0,
					"stars":     1.3,
				},
				"products": []interface{}{
					map[string]interface{}{
						"name": "Ladders",
						"reviews": map[string]interface{}{
							"count": 5,
						},
					},
					map[string]interface{}{
						"name": "Whose the boss?",
						"reviews": map[string]interface{}{
							"count": 0,
						},
					},
					map[string]interface{}{
						"name": "Can I fry that?",
						"reviews": map[string]interface{}{
							"count": 0,
						},
					},
				},
			},
		},
	}
)

func TestExamples(t *testing.T) {
	for _, tC := range testCases {
		t.Run(tC.description, func(t *testing.T) {
			def, err := compiler.CompileReference(tC.transformRef)
			require.Nil(t, err)
			docs, rules, err := resolution.Resolve(tC.documentsRefs, compiler.Importer)
			require.Nil(t, err)
			output := def.Resolve(docs, rules)
			require.Equal(t, tC.output, output)
		})
	}
}

func BenchmarkExamples(b *testing.B) {
	for _, tC := range testCases {
		def, _ := compiler.CompileReference(tC.transformRef)
		docs, rules, _ := resolution.Resolve(tC.documentsRefs, compiler.Importer)
		b.Run(tC.description, func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				value := def.Resolve(docs, rules)
				require.NotNil(b, value)
			}
		})
	}
}
