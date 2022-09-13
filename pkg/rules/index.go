package rules

import "merge-dsl/pkg/cursor"

type Rule interface {
	Execute(global State, local State)
}

type RuleFactory func(map[string]interface{}) Rule

type State struct {
	Documents cursor.Set[interface{}]
	Rules     cursor.Set[cursor.SchemaData]
}
