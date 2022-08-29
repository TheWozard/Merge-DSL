package merge

type (
	Definition struct {
		traversal traversal
	}

	traversal interface {
	}

	objectTraversal struct {
		nodeTraversals map[string]traversal
	}

	arrayTraversal struct {
		defaultTraversal traversal
		idTraversals     map[interface{}]traversal
	}

	leafTraversal struct {
		// allowNull bool
	}
)
