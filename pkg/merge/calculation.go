package merge

import "github.com/mitchellh/mapstructure"

type Operation interface {
	Do(*State)
}

func GetOperation(name string, data map[string]interface{}) Operation {
	if factory, ok := OperationLookup[name]; ok {
		op := factory()
		mapstructure.Decode(data, op)
		return op
	}
	return nil
}

var OperationLookup = map[string]func() Operation{
	"add": func() Operation { return &AddOperation{} },
}

type AddOperation struct {
	Keys []string `mapstructure:"keys"`
}

func (a *AddOperation) Do(local *State) {
	sum := 0
	if parent, ok := local.Parent.Ref.Get(); ok {
		if typed, ok := parent.(map[string]interface{}); ok {
			for _, key := range a.Keys {
				if count, ok := typed[key].(int); ok {
					sum += count
				}
			}
		}
	}
	local.Ref.Update(sum)
}
