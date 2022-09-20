package merge

import (
	"merge-dsl/pkg/cursor"
	"merge-dsl/pkg/result"
)

type (
	DocumentCursorSet = cursor.Set[interface{}]
	// TODO: Turn into a struct for RulesData
	RulesCursorSet = cursor.Set[cursor.SchemaData]
)

func (d Traversal) Resolve(documents DocumentCursorSet, rules RulesCursorSet) interface{} {
	final, ref := result.NewResult(nil)
	d.step.resolve(documents, rules, ref)
	return *final
}

func (o *objectStep) resolve(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) {
	m := ref.Map(o.AllowEmpty, o.AllowNull)
	for key, traversal := range o.nodeSteps {
		ref := m.Key(key)
		traversal.resolve(documents.GetKey(key), rules.GetKey(key), ref)
	}
	applyRules(documents, rules, ref)
}

func (a *arrayStep) resolve(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) {
	s := ref.Slice(a.AllowEmpty, a.AllowNull)
	grouped_nodes := a.sort(documents.GetGroupedItems(cursor.DefaultRawGrouper))
	// Do we care about alerting about rules without ids? Should we do something about this?
	// In general we want to succeed when we can.
	rules_index, _ := rules.GetIndexedItems(cursor.DefaultSchemaGrouper)
	rules_default := rules.GetDefault()
	for _, nodes := range grouped_nodes {
		id := cursor.DefaultRawGrouper(nodes[0])
		traversal := a.defaultStep
		rules := rules_default
		if (id == nil && a.RequireId) || (id != nil && a.ExcludeId) {
			continue
		}
		if id != nil && !a.ExcludeId {
			// Update rules/traversal with ID data
			if id_traversal, ok := a.idStep[id]; ok {
				traversal = id_traversal
			}
			if rules_id, ok := rules_index[id]; ok {
				rules = append(rules, rules_id...)
			}
		}
		ref := s.Append()
		traversal.resolve(nodes, rules, ref)
	}
	applyRules(documents, rules, ref)
}

// Sorts the items into the expected order of the step
func (a *arrayStep) sort(items []cursor.Set[interface{}]) []cursor.Set[interface{}] {
	// TODO: this
	return items
}

func (e *edgeStep) resolve(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) {
	value, _ := documents.Value(cursor.ValidateNonNil)
	if e.Default != nil && value == nil {
		ref.Update(e.Default)
	} else {
		ref.Update(value)
	}
	applyRules(documents, rules, ref)
}

func applyRules(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) {
	// TODO: This
}
