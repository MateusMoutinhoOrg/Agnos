package databaseio

import (
	api "{{.Module}}/sandbox/api"
	database "{{.Module}}/sandbox/deps/database"
)

// databaseio is to a database package what routeio is to a route: the code
// every generated methods.go shares, held in a third package because a
// database may not import another one.
//
// Every reader here converts one stored value in the comma-ok form. A field
// that was never written reads as the zero value of its type — the declaration
// says which fields an insert must carry, and one it does not carry is absent,
// not broken — while a value of the wrong Go type is an error. Nothing here
// asserts a type without checking, so a malformed record can never panic a
// generated method.

// Fail turns one failure the database reported into an error the sandbox
// carries. A nil *database.Error is success and answers nil.
func Fail(sandbox *api.Sandbox, failure *database.Error) error {
	if failure == nil {
		return nil
	}
	if failure.Key != "" {
		return sandbox.Deps.Std.Errorf("%s: %s", failure.Key, failure.Message)
	}
	return sandbox.Deps.Std.Errorf("%s", failure.Message)
}

// Schema resolves one collection of a handle by name. A name the Props does
// not declare is an error rather than a nil instance: a generated method names
// a table its own specs.yaml declared, so this only fires on a handle built
// from another declaration.
func Schema(sandbox *api.Sandbox, handle database.DatabaseHandle, name string) (database.SchemaInstance, error) {
	schema, ok := handle.GetSchema(name)
	if !ok {
		return schema, sandbox.Deps.Std.Errorf("this database declares no table %q", name)
	}
	return schema, nil
}

// ReadString reads one Key or String field of a record.
func ReadString(sandbox *api.Sandbox, item database.SchemaItem, field string) (string, error) {
	raw, failure := read(sandbox, item, field)
	if failure != nil {
		return "", failure
	}
	if raw == nil {
		return "", nil
	}
	value, ok := raw.(string)
	if !ok {
		return "", mistyped(sandbox, field, "text")
	}
	return value, nil
}

// ReadInt reads one Int or Link field of a record.
func ReadInt(sandbox *api.Sandbox, item database.SchemaItem, field string) (int64, error) {
	raw, failure := read(sandbox, item, field)
	if failure != nil {
		return 0, failure
	}
	if raw == nil {
		return 0, nil
	}
	value, ok := raw.(int64)
	if !ok {
		return 0, mistyped(sandbox, field, "a whole number")
	}
	return value, nil
}

// ReadFloat reads one Float field of a record.
func ReadFloat(sandbox *api.Sandbox, item database.SchemaItem, field string) (float64, error) {
	raw, failure := read(sandbox, item, field)
	if failure != nil {
		return 0, failure
	}
	if raw == nil {
		return 0, nil
	}
	value, ok := raw.(float64)
	if !ok {
		return 0, mistyped(sandbox, field, "a number")
	}
	return value, nil
}

// read is the one call every reader shares: the stored value, or nil when the
// record carries none for that field.
func read(sandbox *api.Sandbox, item database.SchemaItem, field string) (any, error) {
	raw, failure := item.Get(field)
	if failure == nil {
		return raw, nil
	}
	if failure.Type == database.NotFound {
		return nil, nil
	}
	return nil, Fail(sandbox, failure)
}

// mistyped words the one failure a reader reports: a stored value that is not
// what the declaration says the field holds.
func mistyped(sandbox *api.Sandbox, field string, expected string) error {
	return sandbox.Deps.Std.Errorf("field %q does not hold %s", field, expected)
}

// TextMatches is the filter a generated <T>Filtrage applies to one text field:
// an empty needle passes everything, so a zero value turns the filter off.
func TextMatches(sandbox *api.Sandbox, value string, starts_with string, equals string) bool {
	if equals != "" && value != equals {
		return false
	}
	if starts_with != "" && !sandbox.Deps.Stringsdeps.HasPrefix(value, starts_with) {
		return false
	}
	return true
}

// IntInRange is TextMatches for a whole-number field: a zero bound is no bound.
func IntInRange(value int64, min int64, max int64) bool {
	if min != 0 && value < min {
		return false
	}
	if max != 0 && value > max {
		return false
	}
	return true
}

// FloatInRange is TextMatches for a floating-point field.
func FloatInRange(value float64, min float64, max float64) bool {
	if min != 0 && value < min {
		return false
	}
	if max != 0 && value > max {
		return false
	}
	return true
}
