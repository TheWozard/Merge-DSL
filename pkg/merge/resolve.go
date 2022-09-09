package merge

import (
	"fmt"
	"merge-dsl/pkg/cursor"
	"merge-dsl/pkg/cursor/validator"
)

type (
	DocumentCursorSet = cursor.CursorSet[cursor.RawData]
	// Turn into a struct for RulesData
	RulesCursorSet = cursor.CursorSet[cursor.SchemaData]
)

func (d Definition) Resolve(documents DocumentCursorSet, rules RulesCursorSet) (interface{}, error) {
	return d.traversal.resolve(documents, rules)
}

func (o *objectTraversal) resolve(documents DocumentCursorSet, rules RulesCursorSet) (interface{}, error) {
	result := map[string]interface{}{}
	for key, traversal := range o.nodeTraversals {
		value, err := traversal.resolve(documents.GetKey(key), rules.GetKey(key))
		if err != nil {
			return nil, fmt.Errorf("%s.%w", key, err)
		}
		// TODO: Run Rules
		if o.allowNull || value != nil {
			result[key] = value
		}
	}
	if o.allowEmpty || len(result) > 0 {
		return result, nil
	}
	return nil, nil
}

func (a *arrayTraversal) resolve(documents DocumentCursorSet, rules RulesCursorSet) (interface{}, error) {
	result := []interface{}{}
	index, order, extra := documents.GetIdsAndExtra(cursor.DefaultRawIndexer)
	rules_index, _, _ := rules.GetIdsAndExtra(cursor.DefaultSchemaIndexer)
	rules_default := rules.GetDefault()
	if !a.excludeId {
		for _, id := range order {
			set := index[id]
			id_rules, ok := rules_index[id]
			if !ok {
				id_rules = rules_default
			}
			traversal, ok := a.idTraversals[id]
			if !ok {
				traversal = a.defaultTraversal
			}
			if traversal != nil {
				value, err := traversal.resolve(set, id_rules)
				if err != nil {
					return nil, fmt.Errorf("[%s].%w", id, err)
				}
				// TODO: Run Rules
				if a.allowNull || value != nil {
					result = append(result, value)
				}
			}
		}
	}
	if !a.requireId {
		for i, set := range extra {
			value, err := a.defaultTraversal.resolve(set, rules_default)
			if err != nil {
				return nil, fmt.Errorf("[%d].%w", i, err)
			}
			// TODO: Run Rules
			if a.allowNull || value != nil {
				result = append(result, value)
			}

		}
	}
	if a.allowEmpty || len(result) > 0 {
		return result, nil
	}
	return nil, nil
}

func (e *edgeTraversal) resolve(documents DocumentCursorSet, rules RulesCursorSet) (interface{}, error) {
	value := documents.Value(validator.NonNil)
	if e.Default != nil && value == nil {
		return e.Default, nil
	}
	return value, nil
}
