package result_test

import (
	"merge-dsl/pkg/result"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testCases = []struct {
		desc       string
		operations func(r *result.Ref)
		output     interface{}
		benchmark  bool
	}{
		{
			desc:       "empty ref",
			operations: func(r *result.Ref) {},
			output:     nil,
		},
		{
			desc: "update ref",
			operations: func(r *result.Ref) {
				r.Update("just a string")
			},
			output:    "just a string",
			benchmark: true,
		},

		// Map
		{
			desc: "basic map",
			operations: func(r *result.Ref) {
				r.Map(false, false).Key("example").Update("map")
			},
			output: map[string]interface{}{
				"example": "map",
			},
			benchmark: true,
		},
		{
			desc: "map allow empty false",
			operations: func(r *result.Ref) {
				r.Map(false, false)
			},
			output: nil,
		},
		{
			desc: "map allow empty true",
			operations: func(r *result.Ref) {
				r.Map(true, false)
			},
			output: map[string]interface{}{},
		},
		{
			desc: "map allow nil false",
			operations: func(r *result.Ref) {
				r.Map(true, false).Key("test").Update(nil)
			},
			output: map[string]interface{}{},
		},
		{
			desc: "map allow nil true",
			operations: func(r *result.Ref) {
				r.Map(true, true).Key("test").Update(nil)
			},
			output: map[string]interface{}{
				"test": nil,
			},
		},
		{
			desc: "map multiple keys",
			operations: func(r *result.Ref) {
				m := r.Map(false, false)
				f1 := m.Key("f1")
				f2 := m.Key("f2")
				f1.Update(1)
				f2.Update(2)
				f1.Update(3)
			},
			output: map[string]interface{}{
				"f1": 3, "f2": 2,
			},
			benchmark: true,
		},
		{
			desc: "map stacked",
			operations: func(r *result.Ref) {
				m1 := r.Map(false, false)
				m2 := m1.Key("m2").Map(false, false)
				m1.Key("f1").Update(1)
				m2.Key("f1").Update(2)
			},
			output: map[string]interface{}{
				"m2": map[string]interface{}{
					"f1": 2,
				},
				"f1": 1,
			},
			benchmark: true,
		},

		// Slice
		{
			desc: "basic slice",
			operations: func(r *result.Ref) {
				r.Slice(false, false).Append().Update("slice")
			},
			output:    []interface{}{"slice"},
			benchmark: true,
		},
		{
			desc: "slice allow empty false",
			operations: func(r *result.Ref) {
				r.Slice(false, false)
			},
			output: nil,
		},
		{
			desc: "slice allow empty true",
			operations: func(r *result.Ref) {
				r.Slice(true, false)
			},
			output: []interface{}{},
		},
		{
			desc: "slice allow nil false",
			operations: func(r *result.Ref) {
				r.Slice(true, false).Append().Update(nil)
			},
			output: []interface{}{},
		},
		{
			desc: "slice allow nil true",
			operations: func(r *result.Ref) {
				r.Slice(true, true).Append().Update(nil)
			},
			output: []interface{}{nil},
		},
	}
)

func TestRef(t *testing.T) {
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			data, ref := result.NewRef(nil)
			tC.operations(ref)
			require.Equal(t, tC.output, *data)
		})
	}
}

func BenchmarkRef(b *testing.B) {
	for _, tC := range testCases {
		if tC.benchmark {
			b.Run(tC.desc, func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					_, ref := result.NewRef(nil)
					tC.operations(ref)
				}
			})
		}
	}
}
