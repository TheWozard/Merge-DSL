package merge

import (
	"fmt"
	"merge-dsl/pkg/cursor"
	"strings"

	"github.com/mitchellh/mapstructure"
)

const (
	TypeKey        = "type"
	ObjectType     = "object"
	ArrayType      = "array"
	EdgeType       = "edge"
	CalculatedType = "calculated"

	OperationKey = "operation"
)

// Compiles the passed golang structure into a ready to use Traversal.
func Compile(document map[string]interface{}) (*Traversal, error) {
	step, err := compile(cursor.NewSchemaCursor(document))
	if err != nil {
		return nil, err
	}
	return &Traversal{
		step: step,
	}, nil
}

// compile uses the type defined on the document to call the right definition compile function
// used recursively to compile the entire tree.
func compile(current cursor.Cursor[cursor.SchemaData]) (step, error) {
	value := current.Value()
	if typ, ok := value[TypeKey].(string); ok {
		switch strings.ToLower(typ) {
		case ObjectType:
			def := &objectStep{}
			return def, def.compile(current)
		case ArrayType:
			def := &arrayStep{}
			return def, def.compile(current)
		case EdgeType:
			def := &edgeStep{}
			return def, def.compile(current)
		case CalculatedType:
			def := &calculatedStep{}
			return def, def.compile(current)
		default:
			return nil, FailedLookup(typ, "compilable types")
		}
	}
	return nil, MissingFieldError(TypeKey)
}

func (o *objectStep) compile(current cursor.Cursor[cursor.SchemaData]) error {
	o.nodeSteps = map[string]step{}
	for _, key := range current.GetKeys() {
		if nextCursor := current.GetKey(key); nextCursor != nil {
			compiled, err := compile(nextCursor)
			if err != nil {
				return fmt.Errorf("failed to compile node '%s': %w", key, err)
			}
			o.nodeSteps[key] = compiled
		}
	}
	return mapstructure.Decode(current.Value(), o)
}

func (a *arrayStep) compile(current cursor.Cursor[cursor.SchemaData]) error {
	// Default Step
	if def := current.GetDefault(); def != nil {
		traversal, err := compile(def)
		if err != nil {
			return fmt.Errorf("failed to compile default: %w", err)
		}
		a.defaultStep = traversal
	}
	// ID Steps
	a.idStep = map[interface{}]step{}
	index, extra := cursor.IndexCursors(current.GetItems(), cursor.DefaultSchemaGrouper)
	if len(extra) > 0 {
		return fmt.Errorf("unexpected non-id node during array compile, all items are expected to contain an id")
	}
	for id, items := range index {
		if len(items) != 1 {
			return fmt.Errorf("found %d instances of the id '%v'", len(items), id)
		}
		traversal, err := compile(items[0])
		if err != nil {
			return fmt.Errorf("failed to compile id traversal: %w", err)
		}
		a.idStep[id] = traversal
	}
	return mapstructure.Decode(current.Value(), a)
}

func (e *edgeStep) compile(current cursor.Cursor[cursor.SchemaData]) error {
	return mapstructure.Decode(current.Value(), e)
}

func (c *calculatedStep) compile(current cursor.Cursor[cursor.SchemaData]) error {
	value := current.Value()
	if name, ok := value[OperationKey].(string); ok {
		name = strings.ToLower(name)
		if op := GetOperation(name, value); op != nil {
			c.Action = op
			return nil
		}
		return FailedLookup(name, "GetOperation")
	}
	return MissingFieldError(OperationKey)
}

func MissingFieldError(field string) error {
	return fmt.Errorf("failed to locate field '%s'", field)
}

func FailedLookup(key, index string) error {
	return fmt.Errorf("failed to find '%s' in %s", key, index)
}
