package merge

type Action = func()

type Actions struct {
	actions []Action
}

func (a *Actions) Do() error {
	for _, action := range a.actions {
		action()
	}
	a.actions = []Action{}
	return nil
}

func (a *Actions) Add(action Action) {
	if a.actions == nil {
		a.actions = []Action{}
	}
	a.actions = append(a.actions, action)
}
