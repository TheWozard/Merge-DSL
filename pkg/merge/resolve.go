package merge

import (
	"fmt"
	"merge-dsl/pkg/cursor"
	"merge-dsl/pkg/result"
)

type (
	DocumentCursorSet = cursor.Set[interface{}]
	// TODO: Turn into a struct for RulesData
	RulesCursorSet = cursor.Set[cursor.SchemaData]
)

func (d Definition) Resolve(documents DocumentCursorSet, rules RulesCursorSet) (interface{}, error) {
	final, ref := result.NewResult(nil)
	err := d.traversal.resolve(documents, rules, ref)
	return *final, err
}

func (o *objectTraversal) resolve(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) error {
	m := ref.Map(o.allowEmpty, o.allowNull)
	for key, traversal := range o.nodeTraversals {
		ref := m.Key(key)
		err := traversal.resolve(documents.GetKey(key), rules.GetKey(key), ref)
		if err != nil {
			return fmt.Errorf("%s.%w", key, err)
		}
		// TODO: Run Rules
	}
	return nil
}

func (a *arrayTraversal) resolve(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) error {
	s := ref.Slice(a.allowEmpty, a.allowNull)
	grouped_nodes := documents.GetGroupedItems(cursor.DefaultRawGrouper)
	rules_index, _ := rules.GetIndexedItems(cursor.DefaultSchemaGrouper)
	rules_default := rules.GetDefault()
	// Sort?
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
		ref := s.Append()
		err := traversal.resolve(nodes, rules, ref)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *edgeTraversal) resolve(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) error {
	value, _ := documents.Value(cursor.ValidateNonNil)
	if e.Default != nil && value == nil {
		ref.Update(e.Default)
	} else {
		ref.Update(value)
	}
	return nil
}
