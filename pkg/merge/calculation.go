package merge

import "github.com/mitchellh/mapstructure"

type Operation interface {
	Do(*State)
}

func GetOperation(name string, data map[string]interface{}) Operation {
	if factory, ok := OperationLookup[name]; ok {
		op := factory()
		err := mapstructure.Decode(data, op)
		if err != nil {
			return nil
		}
		return op
	}
	return nil
}

var OperationLookup = map[string]func() Operation{
	"add":     func() Operation { return &AddOperation{} },
	"average": func() Operation { return &AverageOperation{} },
}

type AddOperation struct {
	Keys []string `mapstructure:"keys"`
}

func (a *AddOperation) Do(local *State) {
	sum := 0
	parent := local.Parent.Ref.Get()
	if typed, ok := parent.(map[string]interface{}); ok {
		for _, key := range a.Keys {
			if count, ok := typed[key].(int); ok {
				sum += count
			}
		}
	}
	local.Ref.Update(sum)
}

type AverageOperation struct {
	Points map[string]int `mapstructure:"points"`
}

func (a *AverageOperation) Do(local *State) {
	sum := 0
	total := 0
	parent := local.Parent.Ref.Get()
	if typed, ok := parent.(map[string]interface{}); ok {
		for key, points := range a.Points {
			if count, ok := typed[key].(int); ok {
				sum += count * points
				total += count
			}
		}
	}
	local.Ref.Update(float64(sum) / float64(total))
}
