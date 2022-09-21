package merge

import "merge-dsl/pkg/result"

type (
	// Traversal defines a compiled merge operation that is ready to be executed.
	Traversal struct {
		step step
	}

	// step defines a part of a resolution process.
	step interface {
		// resolve applies the step to the result based on passed data and rules.
		resolve(state *State)
	}

	State struct {
		Parent         *State
		RootDocuments  DocumentCursorSet
		RootRules      RulesCursorSet
		Documents      DocumentCursorSet
		Rules          RulesCursorSet
		Ref            *result.Ref
		DelayedActions *Actions
	}

	objectStep struct {
		nodeSteps  map[string]step
		AllowEmpty bool `mapstructure:"allow_empty"`
		AllowNull  bool `mapstructure:"allow_null"`
	}

	arrayStep struct {
		defaultStep step
		idStep      map[interface{}]step
		AllowEmpty  bool `mapstructure:"allow_empty"`
		AllowNull   bool `mapstructure:"allow_null"`
		ExcludeId   bool `mapstructure:"exclude_id"`
		RequireId   bool `mapstructure:"require_id"`
	}

	edgeStep struct {
		Default interface{} `mapstructure:"default"`
	}

	calculatedStep struct {
		Action Operation
	}
)
