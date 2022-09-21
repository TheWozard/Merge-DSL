package merge

import "merge-dsl/pkg/result"

func NewRootState(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) *State {
	return &State{
		Parent:         nil,
		RootDocuments:  documents,
		RootRules:      rules,
		Documents:      documents,
		Rules:          rules,
		Ref:            ref,
		DelayedActions: &Actions{},
	}
}

func (c *State) New(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) *State {
	return &State{
		Parent:         c,
		RootDocuments:  c.RootDocuments,
		RootRules:      c.RootRules,
		Documents:      documents,
		Rules:          rules,
		Ref:            ref,
		DelayedActions: c.DelayedActions,
	}
}
