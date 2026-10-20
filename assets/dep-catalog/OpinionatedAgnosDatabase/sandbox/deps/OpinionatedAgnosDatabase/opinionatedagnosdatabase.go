package opinionatedagnosdatabase

import (
	"{{.Module}}/sandbox/deps/databasedeps"
	"{{.Module}}/sandbox/deps/iodeps"
)

// This package is the contract of an *opinionated* lib: unlike every other dep,
// which restates a library's raw capability and nothing more, it carries the
// agnos database mechanic itself — what every generated.methods.go shares:
// resolving a table, reading a stored value in the type its declaration names,
// turning a store failure into an error, and the filters a <T>Filter
// applies.
//
// Every reader converts one stored value in the comma-ok form. A field that
// was never written reads as the zero value of its type, while a value of the
// wrong Go type is an error, so a malformed record can never panic a
// generated method. What stays in the project is each database's database.yaml,
// the generated.api.go, generated.new.go and generated.methods.go generated from it, and its
// methods_custom.go.
//
// It also carries the migration engine: a database.yaml declares the version
// of the data its code expects, each version past the first carries the
// operations that bring the data to it, and Migrate runs them on a copy of the
// database's folder, swapping the copy in only once every step and check
// passed. A generated.migrations.go per database supplies the data; nothing of
// the engine is generated.

// The operations a version runs, the kinds of Op.Kind. Each one is what one
// editing command of a database.yaml changed about the stored data.
const (
	// OpRenameField moves every stored value of Field to the name To. A
	// nested collection (Type object) is moved whole, ids and all.
	OpRenameField = "rename-field"
	// OpRetypeField converts every stored value of Field from Type to the
	// type To; a value that does not convert fails the migration.
	OpRetypeField = "retype-field"
	// OpRequireField gives Default to every record holding no value for
	// Field, when Default is set. A record still holding none fails the
	// check that follows the last step.
	OpRequireField = "require-field"
	// OpRemoveField deletes every stored value of Field, or every record of
	// the nested collection it names.
	OpRemoveField = "remove-field"
	// OpRemoveTable deletes every record of Table.
	OpRemoveTable = "remove-table"
)

// The field types an Op names, spelled the way database.yaml spells them.
const (
	// TypeKey is a unique, indexed text field.
	TypeKey = "key"
	// TypeString is a plain text field.
	TypeString = "string"
	// TypeInteger is a whole-number field.
	TypeInteger = "integer"
	// TypeNumber is a floating-point field.
	TypeNumber = "number"
	// TypeBytes is a binary field.
	TypeBytes = "bytes"
	// TypeLink is the id of a record of another table.
	TypeLink = "link"
	// TypeObject is a collection nested under each record.
	TypeObject = "object"
)

// The names the engine reserves inside the --database folder.
const (
	// VersionTable is the table every database records the version of its
	// data in, under its own key-prefix, so the version travels with the data.
	VersionTable = "agnos-version"
	// MigratingDir is the folder a migration works in: the copy each step
	// runs on, swapped in only once every step and check passed.
	MigratingDir = "agnos-migrating"
	// PreMigrationDir is the folder each migrated database's previous data is
	// kept in, as <key-prefix>-v<version>: undoing a migration is moving it
	// back.
	PreMigrationDir = "agnos-pre-migration"
)

// Op is one operation of a version, as migrations/v<N>.ops.yaml declares it.
// Parent names the object field the field sits in, empty for a field of the
// table itself.
type Op struct {
	// Kind is one of the Op* constants.
	Kind string
	// Table is the table the operation reaches.
	Table string
	// Parent is the nested collection holding Field, "" for the table's own.
	Parent string
	// Field is the field the operation reaches, "" for OpRemoveTable.
	Field string
	// Type is the field's type when the operation runs: the old one, for a
	// rename or a retype.
	Type string
	// To is the new name of OpRenameField and the new type of OpRetypeField.
	To string
	// Default is the value OpRequireField fills, as text; converted to Type.
	Default string
}

// Version is one version of a database: the schema its data has once it
// reached it, and what brings the data of the version before to it.
type Version struct {
	// Version is the number, 1 for the first.
	Version int
	// Schemas is the schema of the data at this version.
	Schemas []databasedeps.Schema
	// Ops is what brings the data of the version before to this one, run in
	// order; empty on version 1.
	Ops []Op
	// Custom is the hand-written step of this version, run after Ops, with
	// from opened on the schema of the version before and to on this one —
	// the same folder read two ways. Nil when the version has none.
	Custom func(from databasedeps.Database, to databasedeps.Database) error
}

// Migrations is one database as the engine reads it, built by its
// generated.migrations.go on every call so Dir follows --database.
type Migrations struct {
	// Name is the database's own name, for the messages.
	Name string
	// Dir is the folder every database lives under (sandbox.Config.DatabaseDir).
	Dir string
	// KeyPrefix is the folder inside Dir the database's records live in.
	KeyPrefix string
	// Versions is every version, the first one first; the last is the one the
	// code expects.
	Versions []Version
}

// MigrateProps is everything Migrate and CheckVersions read. The generated
// code builds it per call, so nothing the lib holds goes stale.
type MigrateProps struct {
	// Databases opens every database the engine reads or writes.
	Databases databasedeps.Databases
	// IoDeps copies, renames and removes the folders of a migration.
	IoDeps iodeps.Contract
	// Logf reports each step taken.
	Logf func(format string, a ...any) (n int, err error)
	// Command is the command line that migrates, named by CheckVersions'
	// failures ("<project> migrate-database").
	Command string
	// Migrations is every database to check or migrate.
	Migrations []Migrations
}

// Contract is the database lib injected whole as the
// Deps.OpinionatedAgnosDatabase field.
type Contract struct {
	// Fail turns one failure the database reported into an error the
	// sandbox carries. A nil *databasedeps.Error is success and answers nil.
	Fail func(failure *databasedeps.Error) error

	// Collection resolves one collection of a database by name. A name the
	// database does not declare, or a database Databases.New refused, is an
	// error rather than a zero collection.
	Collection func(database databasedeps.Database, name string) (databasedeps.Collection, error)

	// ReadString reads one Key or String field of a record.
	ReadString func(record databasedeps.Record, field string) (string, error)

	// ReadInt reads one Int or Link field of a record.
	ReadInt func(record databasedeps.Record, field string) (int64, error)

	// ReadFloat reads one Float field of a record.
	ReadFloat func(record databasedeps.Record, field string) (float64, error)

	// ReadBytes reads one Bytes field of a record, byte for byte.
	ReadBytes func(record databasedeps.Record, field string) ([]byte, error)

	// TextMatches is the filter a generated <T>Filter applies to one text
	// field: an empty needle passes everything, so a zero value turns the
	// filter off.
	TextMatches func(value string, startsWith string, equals string) bool

	// BytesMatches is TextMatches for a binary field, compared byte for
	// byte: an empty needle passes everything.
	BytesMatches func(value []byte, startsWith []byte, equals []byte) bool

	// IntInRange is TextMatches for a whole-number field: a zero bound is no
	// bound.
	IntInRange func(value int64, min int64, max int64) bool

	// FloatInRange is TextMatches for a floating-point field.
	FloatInRange func(value float64, min float64, max float64) bool

	// Migrate brings the data of every database of props up to its last
	// version. Each one is copied into MigratingDir, every step from the
	// version its data records runs on the copy, the copy is repaired and
	// checked against the last schema, and only then is the live folder
	// moved into PreMigrationDir and the copy moved in its place. A failure
	// before that move leaves the live folder as it was; a crash during it is
	// finished by the next call.
	Migrate func(props MigrateProps) error

	// CheckVersions refuses a database whose data is at another version than
	// its code expects, or that a migration left half done, naming
	// props.Command. A database holding no data at all is stamped with the
	// last version instead; one holding data but no version is at version 1,
	// written before versions existed.
	CheckVersions func(props MigrateProps) error
}
