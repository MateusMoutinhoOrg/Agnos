package opinionatedagnosdatabase

import (
	"bytes"
	"fmt"
	"strings"

	"{{.Module}}/sandbox/deps"
	opinionatedagnosdatabase "{{.Module}}/sandbox/deps/OpinionatedAgnosDatabase"
	"{{.Module}}/sandbox/deps/databasedeps"
)

// Bind fills deps.Deps.OpinionatedAgnosDatabase with the readers and filters
// every generated.methods.go shares. Nothing here holds a dep: every function
// is handed the database, the record or the value it reads.
func Bind(deps *deps.Deps) {
	deps.OpinionatedAgnosDatabase = opinionatedagnosdatabase.Contract{
		Fail:         fail,
		Collection:   collection,
		ReadString:   readString,
		ReadInt:      readInt,
		ReadFloat:    readFloat,
		ReadBytes:    readBytes,
		TextMatches:  textMatches,
		BytesMatches: bytesMatches,
		IntInRange:   intInRange,
		FloatInRange: floatInRange,
	}
}

// fail turns one failure the database reported into an error the sandbox
// carries. A nil *databasedeps.Error is success and answers nil.
func fail(failure *databasedeps.Error) error {
	if failure == nil {
		return nil
	}
	if failure.Field != "" {
		return fmt.Errorf("%s: %s", failure.Field, failure.Message)
	}
	return fmt.Errorf("%s", failure.Message)
}

// collection resolves one collection of a database by name. A name the Props
// does not declare is an error rather than a zero collection: a generated
// method names a table its own database.yaml declared, so this only fires on a
// database built from another declaration. A database Databases.New refused is
// the zero Database, whose Collection is nil, and is an error too.
func collection(database databasedeps.Database, name string) (databasedeps.Collection, error) {
	if database.Collection == nil {
		return databasedeps.Collection{}, fmt.Errorf("this database was never built: Databases.New refused its Props")
	}
	collection, ok := database.Collection(name)
	if !ok {
		return collection, fmt.Errorf("this database declares no table %q", name)
	}
	return collection, nil
}

// readString reads one Key or String field of a record.
func readString(record databasedeps.Record, field string) (string, error) {
	raw, failure := read(record, field)
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
func readInt(record databasedeps.Record, field string) (int64, error) {
	raw, failure := read(record, field)
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
func readFloat(record databasedeps.Record, field string) (float64, error) {
	raw, failure := read(record, field)
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

// readBytes reads one Bytes field of a record. The store hands back a copy
// of its own, so the slice is the caller's to write into.
func readBytes(record databasedeps.Record, field string) ([]byte, error) {
	raw, failure := read(record, field)
	if failure != nil {
		return nil, failure
	}
	if raw == nil {
		return nil, nil
	}
	value, ok := raw.([]byte)
	if !ok {
		return nil, mistyped(field, "bytes")
	}
	return value, nil
}

// read is the one call every reader shares: the stored value, or nil when the
// record carries none for that field.
func read(record databasedeps.Record, field string) (any, error) {
	raw, failure := record.Get(field)
	if failure == nil {
		return raw, nil
	}
	if failure.Type == databasedeps.NoValue {
		return nil, nil
	}
	return nil, fail(failure)
}

// mistyped words the one failure a reader reports: a stored value that is not
// what the declaration says the field holds.
func mistyped(field string, expected string) error {
	return fmt.Errorf("field %q does not hold %s", field, expected)
}

// textMatches is the filter a generated <T>Filter applies to one text field:
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

// bytesMatches is textMatches for a binary field, compared byte for byte: an
// empty needle passes everything.
func bytesMatches(value []byte, starts_with []byte, equals []byte) bool {
	if len(equals) > 0 && !bytes.Equal(value, equals) {
		return false
	}
	if len(starts_with) > 0 && !bytes.HasPrefix(value, starts_with) {
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
