package resolution

import (
	"merge-dsl/pkg/cursor"
	"merge-dsl/pkg/reference"
)

type (
	DataSet  = cursor.CursorSet[cursor.RawData]
	RulesSet = cursor.CursorSet[cursor.SchemaData]
)

// Resolve converts a passed list of references into their relevant parts.
func Resolve(references []string, importer reference.Resolver) (DataSet, RulesSet, error) {
	data := DataSet{}
	rules := RulesSet{}
	for _, ref := range references {
		resolution, err := importer.ImportInterface(ref)
		if err != nil {
			return data, rules, err
		}
		switch resolution.Info.Type {
		default:
			data = append(data, cursor.NewRawCursor(resolution.Data))
		}
	}
	return data, rules, nil
}
