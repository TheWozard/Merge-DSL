package merge

type (
	Definition struct {
		traversal traversal
	}

	traversal interface {
		resolve(documents DocumentCursorSet, rules RulesCursorSet) (interface{}, error)
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
		// allowNull bool
	}
)
