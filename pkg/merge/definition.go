package merge

import "merge-dsl/pkg/result"

type (
	Definition struct {
		traversal traversal
	}

	traversal interface {
		resolve(documents DocumentCursorSet, rules RulesCursorSet, ref *result.Ref) error
	}

	objectTraversal struct {
		nodeTraversals map[string]traversal
		allowNull      bool
		allowEmpty     bool
	}

	arrayTraversal struct {
		defaultTraversal traversal
		idTraversals     map[interface{}]traversal
		allowEmpty       bool
		allowNull        bool
		excludeId        bool
		requireId        bool
	}

	edgeTraversal struct {
		Default interface{} `mapstructure:"default"`
	}
)
