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
	final, ref := result.NewRef(nil)
	state := NewRootState(documents, rules, ref)
	d.step.resolve(state)
	state.DelayedActions.Do()
	return *final
}

func (o *objectStep) resolve(state *State) {
	m := state.Ref.Map(o.AllowEmpty, o.AllowNull)
	for key, traversal := range o.nodeSteps {
		ref := m.Key(key)
		traversal.resolve(state.New(state.Documents.GetKey(key), state.Rules.GetKey(key), ref))
	}
	applyRules(state)
}

func (a *arrayStep) resolve(state *State) {
	s := state.Ref.Slice(a.AllowEmpty, a.AllowNull)
	grouped_nodes := a.sort(state.Documents.GetGroupedItems(cursor.DefaultRawGrouper))
	// Do we care about alerting about rules without ids? Should we do something about this?
	// In general we want to succeed when we can.
	rules_index, _ := state.Rules.GetIndexedItems(cursor.DefaultSchemaGrouper)
	rules_default := state.Rules.GetDefault()
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
		traversal.resolve(state.New(nodes, rules, ref))
	}
	applyRules(state)
}

// Sorts the items into the expected order of the step
func (a *arrayStep) sort(items []cursor.Set[interface{}]) []cursor.Set[interface{}] {
	// TODO: this
	return items
}

func (e *edgeStep) resolve(state *State) {
	value, _ := state.Documents.Value(cursor.ValidateNonNil)
	if e.Default != nil && value == nil {
		state.Ref.Update(e.Default)
	} else {
		state.Ref.Update(value)
	}
	applyRules(state)
}

func (c *calculatedStep) resolve(state *State) {
	state.DelayedActions.Add(func() {
		c.Action.Do(state)
	})
}

func applyRules(state *State) {
	// TODO: This
}
