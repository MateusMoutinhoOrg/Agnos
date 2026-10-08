package opinatedagnosdatabase

import (
	"{{.Module}}/sandbox/deps/database"
)

// This package is the contract of an *opinated* lib: unlike every other dep,
// which restates a library's raw capability and nothing more, it carries the
// agnos database mechanic itself — what every generated methods.go shares:
// resolving a table, reading a stored value in the type its declaration names,
// turning a store failure into an error, and the filters a <T>Filtrage
// applies.
//
// Every reader converts one stored value in the comma-ok form. A field that
// was never written reads as the zero value of its type, while a value of the
// wrong Go type is an error, so a malformed record can never panic a
// generated method. What stays in the project is each database's specs.yaml,
// the api.go, new.go and methods.go generated from it, and its
// methods_custom.go.

// Sandbox is the database lib injected whole as the
// Deps.OpinatedAgnosDatabase field.
type Sandbox struct {
	// Fail turns one failure the database reported into an error the
	// sandbox carries. A nil *database.Error is success and answers nil.
	Fail func(failure *database.Error) error

	// Schema resolves one collection of a handle by name. A name the handle
	// does not declare is an error rather than a nil instance.
	Schema func(handle database.DatabaseHandle, name string) (database.SchemaInstance, error)

	// ReadString reads one Key or String field of a record.
	ReadString func(item database.SchemaItem, field string) (string, error)

	// ReadInt reads one Int or Link field of a record.
	ReadInt func(item database.SchemaItem, field string) (int64, error)

	// ReadFloat reads one Float field of a record.
	ReadFloat func(item database.SchemaItem, field string) (float64, error)

	// TextMatches is the filter a generated <T>Filtrage applies to one text
	// field: an empty needle passes everything, so a zero value turns the
	// filter off.
	TextMatches func(value string, startsWith string, equals string) bool

	// IntInRange is TextMatches for a whole-number field: a zero bound is no
	// bound.
	IntInRange func(value int64, min int64, max int64) bool

	// FloatInRange is TextMatches for a floating-point field.
	FloatInRange func(value float64, min float64, max float64) bool
}
