package merge

import (
	"fmt"
	"merge-dsl/pkg/cursor"
)

type (
	DocumentCursorSet = cursor.Set[interface{}]
	// TODO: Turn into a struct for RulesData
	RulesCursorSet = cursor.Set[cursor.SchemaData]
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
	grouped_nodes := documents.GetGroupedItems(cursor.DefaultRawGrouper)
	rules_index, _ := rules.GetIndexedItems(cursor.DefaultSchemaGrouper)
	rules_default := rules.GetDefault()
	for _, nodes := range grouped_nodes {
		id := cursor.DefaultRawGrouper(nodes[0])
		traversal := a.defaultTraversal
		rules := rules_default
		if (id == nil && a.requireId) || (id != nil && a.excludeId) {
			continue
		}
		if id != nil && !a.excludeId {
			if id_traversal, ok := a.idTraversals[id]; ok {
				traversal = id_traversal
			}
			if rules_id, ok := rules_index[id]; ok {
				rules = append(rules, rules_id...)
			}
		}
		value, err := traversal.resolve(nodes, rules)
		if err != nil {
			return nil, err
		}
		result = append(result, value)
	}
	if a.allowEmpty || len(result) > 0 {
		return result, nil
	}
	return nil, nil
}

func (e *edgeTraversal) resolve(documents DocumentCursorSet, rules RulesCursorSet) (interface{}, error) {
	value, _ := documents.Value(cursor.ValidateNonNil)
	if e.Default != nil && value == nil {
		return e.Default, nil
	}
	return value, nil
}
