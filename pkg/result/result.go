package result

// TODO: Changed BenchmarkExamples/business-8 from 4.6k ns/op to 8.2k ns/op.
// Overall a very reasonable performance loss, but see if we can slim it down any.

// NewResult returns a new pointer to init value and a Ref for updating the pointer
// The goal is to be able to independently updated specific nodes on the result tree through stored Refs
// This give us the ability to store a Ref in a structure and update the value that Ref points to later
// during calculated fields.
func NewResult(init interface{}) (*interface{}, *Ref) {
	var value *interface{} = &init
	return value, &Ref{
		Update: func(i interface{}) {
			*value = i
		},
		Get: func() (interface{}, bool) {
			return *value, true
		},
		Clear: func() {
			*value = nil
		},
	}
}

// Ref defines an update-able section of data
type Ref struct {
	Update func(interface{})
	Get    func() (interface{}, bool)
	Clear  func()
}

// Map converts the Ref into a new Map reference generator.
func (r *Ref) Map(allowEmpty, allowNil bool) *Map {
	ref := &Map{
		ref:        r,
		allowEmpty: allowEmpty,
		allowNil:   allowNil,
		value:      map[string]interface{}{},
	}
	ref.sync()
	return ref
}

// Slice converts the Ref into a new Slice reference generator.
func (r *Ref) Slice(allowEmpty, allowNil bool) *mapSlice {
	ref := &mapSlice{
		ref:        r,
		allowEmpty: allowEmpty,
		allowNil:   allowNil,
		value:      []interface{}{},
		lookup:     []int{},
	}
	ref.sync()
	return ref
}

type Map struct {
	ref        *Ref
	allowEmpty bool
	allowNil   bool
	value      map[string]interface{}
}

func (m *Map) Key(key string) *Ref {
	clear := func() {
		delete(m.value, key)
		m.sync()
	}
	return &Ref{
		Update: func(value interface{}) {
			if m.allowNil || value != nil {
				m.value[key] = value
				m.sync()
			} else {
				clear()
			}
		},
		Get: func() (interface{}, bool) {
			value, ok := m.value[key]
			return value, ok
		},
		Clear: clear,
	}
}

func (m *Map) sync() {
	if m.allowEmpty || len(m.value) > 0 {
		m.ref.Update(m.value)
	} else {
		m.ref.Update(nil)
	}
}

type mapSlice struct {
	ref        *Ref
	allowEmpty bool
	allowNil   bool
	value      []interface{}
	lookup     []int
}

func (s *mapSlice) Append() *Ref {
	instance := len(s.lookup)
	s.lookup = append(s.lookup, -1)
	return &Ref{
		Update: func(value interface{}) {
			if s.allowNil || value != nil {
				// Need to update
				pos := s.lookup[instance]
				if pos < 0 {
					// We need to add the value and update the index
					target := 0
					// Find the first instance that is instantiated and pointing to an index.
					// Fall back to prepending
					for i := instance - 1; i >= 0; i-- {
						if s.lookup[i] >= 0 {
							target = s.lookup[i] + 1
							break
						}
					}
					// Add in instance to value
					s.lookup[instance] = target
					if target == len(s.value) {
						s.value = append(s.value, value)
					} else {
						s.value = append(s.value[:target+1], s.value[target:])
						s.value[target] = value
					}
					// Update all instances after the one being added.
					for i := instance + 1; i < len(s.lookup); i++ {
						s.lookup[i] = s.lookup[i] + 1
					}
					s.sync()
				} else {
					s.value[pos] = value
					s.sync()
				}
			} else {
				// Invalid value so clear
				s.clear(instance)
			}
		},
		Get: func() (interface{}, bool) {
			pos := s.lookup[instance]
			if pos < 0 {
				return nil, false
			}
			return s.value[pos], true
		},
		Clear: func() {
			s.clear(instance)
		},
	}
}

func (s *mapSlice) clear(instance int) {
	pos := s.lookup[instance]
	if pos > -1 {
		s.value = append(s.value[:pos], s.value[pos+1:])
		s.lookup[instance] = -1
		for i := instance + 1; i < len(s.lookup); i++ {
			s.lookup[i] = s.lookup[i] - 1
		}
		s.sync()
	}
}

func (s *mapSlice) sync() {
	if s.allowEmpty || len(s.value) > 0 {
		s.ref.Update(s.value)
	} else {
		s.ref.Update(nil)
	}
}
