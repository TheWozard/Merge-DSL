package merge

import (
	"fmt"
	"merge-dsl/pkg/cursor"
	"merge-dsl/pkg/reference"

	"github.com/mitchellh/mapstructure"
)

const (
	// TODO: what would moving this to the compiler look like
	TransformSchemaReference = "schema://transform_schema.yaml"

	TypeKey    = "type"
	ObjectType = "object"
	ArrayType  = "array"
	EdgeType   = "edge"
)

type DocumentCursor = cursor.Cursor[cursor.SchemaData]

type Compiler struct {
	Importer  reference.Resolver
	Validator *reference.SchemaValidator
}

// CompileReference imports the passed reference and passes it to Compile.
func (c *Compiler) CompileReference(ref string) (*Definition, error) {
	resolved, err := c.Importer.ImportMap(ref)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve reference for compiling: %w", err)
	}
	return c.Compile(resolved.Data)
}

// Compiles the passed golang structure into a ready to use Definition.
func (c *Compiler) Compile(document map[string]interface{}) (*Definition, error) {
	err := c.Validator.IsValidByReference(document, TransformSchemaReference)
	if err != nil {
		return nil, fmt.Errorf("failed to validated document: %w", err)
	}
	current := cursor.NewSchemaCursor(document)
	if current == nil {
		return nil, fmt.Errorf("cannot compile nil cursor")
	}
	traversal, err := compile(current)
	if err != nil {
		return nil, err
	}
	return &Definition{
		traversal: traversal,
	}, nil
}

// compile uses the type defined on the document to call the right definition compile function
// used recursively to compile the entire tree.
func compile(current DocumentCursor) (traversal, error) {
	value := current.Value()
	if typ, ok := value[TypeKey].(string); ok {
		switch typ {
		case ObjectType:
			def := &objectTraversal{}
			return def, def.compile(current, value)
		case ArrayType:
			def := &arrayTraversal{}
			return def, def.compile(current, value)
		case EdgeType:
			def := &edgeTraversal{}
			return def, def.compile(current, value)
		default:
			// Should be caught by validation, just in case.
			return nil, fmt.Errorf("unknown compile type '%s'", typ)
		}
	}
	return nil, fmt.Errorf("failed to locate type in definition")
}

func (o *objectTraversal) compile(current DocumentCursor, data map[string]interface{}) error {
	o.nodeTraversals = map[string]traversal{}
	for _, key := range current.GetKeys() {
		if nextCursor := current.GetKey(key); nextCursor != nil {
			compiled, err := compile(nextCursor)
			if err != nil {
				return fmt.Errorf("failed to compile node '%s': %w", key, err)
			}
			o.nodeTraversals[key] = compiled
		}
	}
	return mapstructure.Decode(data, o)
}

func (a *arrayTraversal) compile(current DocumentCursor, data map[string]interface{}) error {
	if def := current.GetDefault(); def != nil {
		traversal, err := compile(def)
		if err != nil {
			return fmt.Errorf("failed to compile default: %w", err)
		}
		a.defaultTraversal = traversal
	}
	a.idTraversals = map[interface{}]traversal{}
	index, extra := cursor.IndexCursors(current.GetItems(), cursor.DefaultSchemaGrouper)
	if len(extra) > 0 {
		return fmt.Errorf("unexpected non-id node during array compile")
	}
	for id, items := range index {
		if len(items) != 1 {
			return fmt.Errorf("found %d instances of the id '%v'", len(items), id)
		}
		traversal, err := compile(items[0])
		if err != nil {
			return fmt.Errorf("failed to compile id traversal: %w", err)
		}
		a.idTraversals[id] = traversal
	}
	return mapstructure.Decode(data, a)
}

func (e *edgeTraversal) compile(current DocumentCursor, data map[string]interface{}) error {
	return mapstructure.Decode(data, e)
}
