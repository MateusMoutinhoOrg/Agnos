package opinatedagnosdatabase

import (
	"fmt"
	"strings"

	"{{.Module}}/sandbox/deps"
	opinatedagnosdatabase "{{.Module}}/sandbox/deps/OpinatedAgnosDatabase"
	"{{.Module}}/sandbox/deps/database"
)

// Bind fills deps.Deps.OpinatedAgnosDatabase with the readers and filters
// every generated methods.go shares. Nothing here holds a dep: every function
// is handed the handle, the record or the value it reads.
func Bind(deps *deps.Deps) {
	deps.OpinatedAgnosDatabase = opinatedagnosdatabase.Sandbox{
		Fail:         fail,
		Schema:       schema,
		ReadString:   readString,
		ReadInt:      readInt,
		ReadFloat:    readFloat,
		TextMatches:  textMatches,
		IntInRange:   intInRange,
		FloatInRange: floatInRange,
	}
}

// fail turns one failure the database reported into an error the sandbox
// carries. A nil *database.Error is success and answers nil.
func fail(failure *database.Error) error {
	if failure == nil {
		return nil
	}
	if failure.Key != "" {
		return fmt.Errorf("%s: %s", failure.Key, failure.Message)
	}
	return fmt.Errorf("%s", failure.Message)
}

// schema resolves one collection of a handle by name. A name the Props does
// not declare is an error rather than a nil instance: a generated method names
// a table its own specs.yaml declared, so this only fires on a handle built
// from another declaration.
func schema(handle database.DatabaseHandle, name string) (database.SchemaInstance, error) {
	schema, ok := handle.GetSchema(name)
	if !ok {
		return schema, fmt.Errorf("this database declares no table %q", name)
	}
	return schema, nil
}

// readString reads one Key or String field of a record.
func readString(item database.SchemaItem, field string) (string, error) {
	raw, failure := read(item, field)
	if failure != nil {
		return "", failure
	}
	if raw == nil {
		return "", nil
	}
	value, ok := raw.(string)
	if !ok {
		return "", mistyped(field, "text")
	}
	return value, nil
}

// readInt reads one Int or Link field of a record.
func readInt(item database.SchemaItem, field string) (int64, error) {
	raw, failure := read(item, field)
	if failure != nil {
		return 0, failure
	}
	if raw == nil {
		return 0, nil
	}
	value, ok := raw.(int64)
	if !ok {
		return 0, mistyped(field, "a whole number")
	}
	return value, nil
}

// readFloat reads one Float field of a record.
func readFloat(item database.SchemaItem, field string) (float64, error) {
	raw, failure := read(item, field)
	if failure != nil {
		return 0, failure
	}
	if raw == nil {
		return 0, nil
	}
	value, ok := raw.(float64)
	if !ok {
		return 0, mistyped(field, "a number")
	}
	return value, nil
}

// read is the one call every reader shares: the stored value, or nil when the
// record carries none for that field.
func read(item database.SchemaItem, field string) (any, error) {
	raw, failure := item.Get(field)
	if failure == nil {
		return raw, nil
	}
	if failure.Type == database.NotFound {
		return nil, nil
	}
	return nil, fail(failure)
}

// mistyped words the one failure a reader reports: a stored value that is not
// what the declaration says the field holds.
func mistyped(field string, expected string) error {
	return fmt.Errorf("field %q does not hold %s", field, expected)
}

// textMatches is the filter a generated <T>Filtrage applies to one text field:
// an empty needle passes everything, so a zero value turns the filter off.
func textMatches(value string, starts_with string, equals string) bool {
	if equals != "" && value != equals {
		return false
	}
	if starts_with != "" && !strings.HasPrefix(value, starts_with) {
		return false
	}
	return true
}

// intInRange is textMatches for a whole-number field: a zero bound is no bound.
func intInRange(value int64, min int64, max int64) bool {
	if min != 0 && value < min {
		return false
	}
	if max != 0 && value > max {
		return false
	}
	return true
}

// floatInRange is textMatches for a floating-point field.
func floatInRange(value float64, min float64, max float64) bool {
	if min != 0 && value < min {
		return false
	}
	if max != 0 && value > max {
		return false
	}
	return true
}
