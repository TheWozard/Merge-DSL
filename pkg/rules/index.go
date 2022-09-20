package rules

import (
	"merge-dsl/pkg/cursor"
	"merge-dsl/pkg/result"
)

type Rule interface {
	Execute(global State, local State, ref *result.Ref)
}

type RuleFactory func(map[string]interface{}) Rule

type State struct {
	Documents cursor.Set[interface{}]
	Rules     cursor.Set[cursor.SchemaData]
}
