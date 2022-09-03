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

	edgeTraversal struct {
		// allowNull bool
	}
)
