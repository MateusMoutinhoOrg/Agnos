package import_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serializables "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/serializables"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/routeconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// formatPatterns is what --infer-format reads a string back as: the four
// formats the schema subset spells, each by the shape its values have. A
// string matching none of them is a string, which is the answer in every case
// the table does not cover.
var formatPatterns = []struct {
	Format  string
	Pattern string
}{
	{"uuid", `^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`},
	{"date-time", `^\d{4}-\d{2}-\d{2}[Tt ]\d{2}:\d{2}:\d{2}(\.\d+)?([Zz]|[+-]\d{2}:\d{2})$`},
	{"email", `^[^@\s]+@[^@\s]+\.[^@\s]+$`},
	{"uri", `^[a-zA-Z][a-zA-Z0-9+.-]*://[^\s]+$`},
}

// inferSchema reads one node of the example as the json-schema node that
// accepts it: an object becomes an object of its keys, a list an array of
// whatever its first item is, and a scalar the type it is written as. A json
// null says only that the key may be null, so it becomes a nullable string —
// the type an example cannot name is the one the person changes afterwards.
func inferSchema(sandbox *api.Sandbox, raw string, props api.RouteBodyImportProps) (*routeconf.Schema, error) {
	document, err := sandbox.Deps.Serializables.ParseJson(raw)
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("the example is not valid json: %w", err)
	}
	return inferNode(sandbox, document, props)
}

// inferNode is inferSchema over one already-parsed node.
func inferNode(sandbox *api.Sandbox, node *serializables.SerializibleObject, props api.RouteBodyImportProps) (*routeconf.Schema, error) {
	if node == nil || node.IsNull() {
		return &routeconf.Schema{Type: "string", Nullable: true}, nil
	}
	if node.IsObject() {
		return inferObject(sandbox, node, props)
	}
	if node.IsArray() {
		return inferArray(sandbox, node, props)
	}
	if node.IsBool() {
		return &routeconf.Schema{Type: "boolean"}, nil
	}
	if node.IsInt() {
		return &routeconf.Schema{Type: "integer"}, nil
	}
	if node.IsFloat() {
		return &routeconf.Schema{Type: "number"}, nil
	}

	return inferString(sandbox, node, props)
}

// inferObject declares one property per key the example object carries, in the
// alphabetical order every declaration is read back in. --required lists every
// one of them: a key the example carries is a key the payload has, which is
// the most an example can say about it.
func inferObject(sandbox *api.Sandbox, node *serializables.SerializibleObject, props api.RouteBodyImportProps) (*routeconf.Schema, error) {
	schema := &routeconf.Schema{Type: "object"}

	keys, err := node.GetKeys()
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not read the keys of the example: %w", err)
	}

	for _, key := range keys {
		item, err := node.GetObjectItem(key)
		if err != nil {
			return nil, sandbox.Deps.Std.Errorf("could not read %q out of the example: %w", key, err)
		}

		child, err := inferNode(sandbox, item, props)
		if err != nil {
			return nil, err
		}

		utils.InsertSchemaProperty(schema, key, child)
		if props.Required {
			schema.Required = utils.AppendUnique(schema.Required, []string{key})
		}
	}

	sandbox.Deps.Sortdeps.Strings(schema.Required)
	return schema, nil
}

// inferArray reads a list as an array of whatever its first item is. An empty
// list names no element type at all, so it becomes an array of strings — the
// same answer a json null gets, and for the same reason.
func inferArray(sandbox *api.Sandbox, node *serializables.SerializibleObject, props api.RouteBodyImportProps) (*routeconf.Schema, error) {
	size, err := node.GetArraySize()
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not read the length of a list in the example: %w", err)
	}
	if size == 0 {
		return &routeconf.Schema{Type: "array", Items: &routeconf.Schema{Type: "string"}}, nil
	}

	items, err := inferNode(sandbox, node.GetArrayItem(0), props)
	if err != nil {
		return nil, err
	}
	return &routeconf.Schema{Type: "array", Items: items}, nil
}

// inferString reads one text value, and with --infer-format the format it
// spells.
func inferString(sandbox *api.Sandbox, node *serializables.SerializibleObject, props api.RouteBodyImportProps) (*routeconf.Schema, error) {
	value, err := node.GetString()
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("could not read a value of the example: %w", err)
	}

	schema := &routeconf.Schema{Type: "string"}
	if !props.InferFormat {
		return schema, nil
	}

	for _, candidate := range formatPatterns {
		matched, err := sandbox.Deps.Stringsdeps.MatchPattern(candidate.Pattern, value)
		if err != nil {
			return nil, err
		}
		if matched {
			schema.Format = candidate.Format
			return schema, nil
		}
	}

	return schema, nil
}

// mergeSchema declares into `into` every property `from` carries that is not
// there yet, and reports both halves: what was added, and what was left as it
// was. A property already declared is never written over — an example says
// what a payload looks like, and a bound already declared says what it is
// allowed to be, which is the one the person typed on purpose.
func mergeSchema(sandbox *api.Sandbox, into *routeconf.Schema, from *routeconf.Schema, prefix string) ([]string, []string) {
	added := []string{}
	skipped := []string{}

	for _, property := range from.Properties {
		path := property.Name
		if prefix != "" {
			path = prefix + "." + property.Name
		}

		current := utils.SchemaPropertyOf(into, property.Name)
		if current == nil {
			utils.InsertSchemaProperty(into, property.Name, property.Schema)
			added = append(added, path)
			if utils.SchemaDemands(from, property.Name) {
				into.Required = utils.AppendUnique(into.Required, []string{property.Name})
			}
			continue
		}

		mine := utils.SchemaObjectOf(current)
		theirs := utils.SchemaObjectOf(property.Schema)
		if mine != nil && theirs != nil && mine.Type == "object" && theirs.Type == "object" {
			deeper, left := mergeSchema(sandbox, mine, theirs, path)
			added = append(added, deeper...)
			skipped = append(skipped, left...)
			continue
		}

		skipped = append(skipped, path)
	}

	sandbox.Deps.Sortdeps.Strings(into.Required)
	return added, skipped
}
