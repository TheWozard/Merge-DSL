package result

// Ref
//------------------------------------------------------------------------------

// NewRef returns a new pointer to init value and a Ref for updating the pointer
func NewRef(init interface{}) (*interface{}, *Ref) {
	return &init, &Ref{&valueUpdater{value: &init}}
}

// Ref defines an update-able reference to a section of an overall object
type Ref struct {
	data Updater
}

// Update defines a nil safe update of a value
func (r *Ref) Update(value interface{}) {
	if r != nil {
		r.data.Update(value)
	}
}

// Get defines a nil safe get of the value
func (r *Ref) Get() interface{} {
	if r != nil {
		return r.data.Get()
	}
	return nil
}

// Clean defines a nil safe clear of the value
func (r *Ref) Clear() {
	if r != nil {
		r.data.Clear()
	}
}

// Map converts the Ref into a new Map reference generator
func (r *Ref) Map(allowEmpty, allowNil bool) *Object {
	if r == nil {
		return nil
	}
	ref := &Object{
		Ref:        r,
		AllowEmpty: allowEmpty,
		AllowNil:   allowNil,
	}
	if prev, ok := r.data.Get().(map[string]interface{}); ok {
		ref.value = prev
	} else {
		ref.value = map[string]interface{}{}
		ref.Sync()
	}
	return ref
}

// Slice converts the Ref into a new Slice reference generator
func (r *Ref) Slice(allowEmpty, allowNil bool) *Array {
	if r == nil {
		return nil
	}
	ref := &Array{
		Ref:        r,
		AllowEmpty: allowEmpty,
		AllowNil:   allowNil,
		lookup:     []int{},
	}
	if prev, ok := r.data.Get().([]interface{}); ok {
		ref.value = prev
	} else {
		ref.value = []interface{}{}
		ref.Sync()
	}
	ref.Sync()
	return ref
}

// Intermediaries
//------------------------------------------------------------------------------

type Object struct {
	Ref        *Ref
	AllowEmpty bool
	AllowNil   bool
	value      map[string]interface{}
}

// Key creates a new reference for a key
func (o *Object) Key(key string) *Ref {
	return &Ref{&objectUpdater{key, o}}
}

func (o *Object) Sync() {
	if o.AllowEmpty || len(o.value) > 0 {
		o.Ref.Update(o.value)
	} else {
		o.Ref.Update(nil)
	}
}

type Array struct {
	Ref        *Ref
	AllowEmpty bool
	AllowNil   bool
	value      []interface{}
	lookup     []int
}

func (a *Array) Append() *Ref {
	instance := len(a.lookup)
	a.lookup = append(a.lookup, -1)
	return &Ref{&arrayUpdater{instance, a}}
}

func (a *Array) Index(i int) *Ref {
	if i >= len(a.value) {
		return nil
	}
	for instance, index := range a.lookup {
		if index == i {
			return &Ref{&arrayUpdater{instance, a}}
		}
		if index > i {
			a.lookup = append(append(a.lookup[:instance], index), a.lookup[instance:]...)
			return &Ref{&arrayUpdater{instance, a}}
		}
	}
	instance := len(a.lookup)
	a.lookup = append(a.lookup, i)
	return &Ref{&arrayUpdater{instance, a}}
}

func (a *Array) Sync() {
	if a.AllowEmpty || len(a.value) > 0 {
		a.Ref.Update(a.value)
	} else {
		a.Ref.Update(nil)
	}
}

// Updaters
//------------------------------------------------------------------------------

type Updater interface {
	Update(interface{})
	Get() interface{}
	Clear()
}

// valueUpdater updater for a pointer to a value
type valueUpdater struct {
	value *interface{}
}

func (vu *valueUpdater) Update(value interface{}) {
	*vu.value = value
}

func (vu *valueUpdater) Get() interface{} {
	return *vu.value
}

func (vu *valueUpdater) Clear() {
	vu.Update(nil)
}

// objectUpdater updater for a key on an Object
type objectUpdater struct {
	key    string
	object *Object
}

func (ou *objectUpdater) Update(value interface{}) {
	if ou.object.AllowNil || value != nil {
		ou.object.value[ou.key] = value
		ou.object.Sync()
	} else {
		ou.Clear()
	}
}

func (ou *objectUpdater) Get() interface{} {
	return ou.object.value[ou.key]
}

func (ou *objectUpdater) Clear() {
	delete(ou.object.value, ou.key)
	ou.object.Sync()
}

// arrayUpdater updater for a key on an instance on an Array
type arrayUpdater struct {
	instance int
	array    *Array
}

func (au *arrayUpdater) Update(value interface{}) {
	if au.array.AllowNil || value != nil {
		// Need to update
		pos := au.array.lookup[au.instance]
		if pos < 0 {
			// We need to add the value and update the index
			target := 0
			// Find the first instance that is instantiated and pointing to an index
			// Fall back to prepending
			for i := au.instance - 1; i >= 0; i-- {
				if au.array.lookup[i] >= 0 {
					target = au.array.lookup[i] + 1
					break
				}
			}
			// Add in instance to value
			au.array.lookup[au.instance] = target
			if target == len(au.array.value) {
				au.array.value = append(au.array.value, value)
			} else {
				au.array.value = append(au.array.value[:target+1], au.array.value[target:])
				au.array.value[target] = value
			}
			// Update all instances after the one being added
			for i := au.instance + 1; i < len(au.array.lookup); i++ {
				au.array.lookup[i] = au.array.lookup[i] + 1
			}
			au.array.Sync()
		} else {
			au.array.value[pos] = value
			au.array.Sync()
		}
	} else {
		au.Clear()
	}
}

func (au *arrayUpdater) Get() interface{} {
	pos := au.array.lookup[au.instance]
	if pos < 0 {
		return nil
	}
	return au.array.value[pos]
}

func (au *arrayUpdater) Clear() {
	pos := au.array.lookup[au.instance]
	if pos > -1 {
		au.array.value = append(au.array.value[:pos], au.array.value[pos+1:])
		// Update lookup
		au.array.lookup[au.instance] = -1
		for i := au.instance + 1; i < len(au.array.lookup); i++ {
			au.array.lookup[i] = au.array.lookup[i] - 1
		}
		au.array.Sync()
	}
}
