package databaseconf

// The field types a table may declare, each one a databasedeps.Field kind the
// Keep contract carries. They are the whole vocabulary of a database.yaml: a
// `type` outside this list is refused by the parser rather than generated into
// a method nothing can call.
const (
	// FieldKey is a unique, indexed text field. It is the only kind a
	// Find<T>By<Field> is generated for, because it is the only one the
	// database indexes.
	FieldKey = "key"
	// FieldString is a plain text field, reachable through List<T> alone.
	FieldString = "string"
	// FieldInteger is a plain integer field, carried as an int64.
	FieldInteger = "integer"
	// FieldNumber is a plain floating-point field, carried as a float64.
	FieldNumber = "number"
	// FieldLink is a reference to a record of the table its Target names,
	// stored as that record's id and resolved by Get<T><Field>.
	FieldLink = "link"
	// FieldObject is a collection nested under each record of the table.
	// It carries Fields of its own and is never read as a value.
	FieldObject = "object"
)

// FieldTypes is every type a field may declare, in the order the doc and the
// error messages spell them.
var FieldTypes = []string{FieldKey, FieldString, FieldInteger, FieldNumber, FieldLink, FieldObject}

// Field is one column of a table as database.yaml declares it. Target is filled
// on a link alone and names the table it points at; Fields is filled on a
// nested database alone, and Required is meaningless there — a nested
// collection is never provided by an insert.
type Field struct {
	Name     string
	Type     string
	Required bool
	Target   string
	Fields   []Field
}

// Table is one collection of records: the name every generated method is
// spelled after, and the fields each record can hold.
type Table struct {
	Name   string
	Fields []Field
}

// DatabaseConf is the parsed form of
// sandbox/internal/databases/<db>/database.yaml — the declarative description of
// one database, which `agnos build` turns into that package's api.go, new.go
// and methods.go. It is written by `add-database` and rewritten by
// `add-table` / `add-table-field` and their editors, never by hand.
type DatabaseConf struct {
	Name      string
	KeyPrefix string
	Tables    []Table

	// Render serializes the declaration back to the database.yaml shape.
	Render func() string
}
