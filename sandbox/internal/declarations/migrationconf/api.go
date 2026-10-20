package migrationconf

// The operations a version may declare, each one what an editing command of a
// database.yaml changed about the data already stored. They are the
// OpinionatedAgnosDatabase lib's Op kinds, which runs them; a `kind` outside
// this list is refused by the parser.
const (
	// OpRenameField moves every stored value of a field to its new name.
	OpRenameField = "rename-field"
	// OpRetypeField converts every stored value of a field to its new type.
	OpRetypeField = "retype-field"
	// OpRequireField fills a field that became required with its default, and
	// leaves the records that still hold none for the check after the last
	// step to refuse.
	OpRequireField = "require-field"
	// OpRemoveField deletes every stored value of a field.
	OpRemoveField = "remove-field"
	// OpRemoveTable deletes every record of a table.
	OpRemoveTable = "remove-table"
)

// OpKinds is every kind an op may declare.
var OpKinds = []string{OpRenameField, OpRetypeField, OpRequireField, OpRemoveField, OpRemoveTable}

// Op is one operation of a version. Parent names the object field the field
// sits in, empty for a field of the table itself; Type is the field's type
// when the operation runs — the old one, for a rename or a retype — and To the
// new name of a rename and the new type of a retype. Default is the value a
// require-field fills, as typed on the command line.
type Op struct {
	Kind    string
	Table   string
	Parent  string
	Field   string
	Type    string
	To      string
	Default string
}

// MigrationConf is the parsed form of
// sandbox/internal/databases/<db>/migrations/v<N>.ops.yaml: what brings the
// data of version N-1 to version N, run in order. It is written by the
// editing command that opened version N, never by hand.
type MigrationConf struct {
	Ops []Op

	// Render serializes the declaration back to the v<N>.ops.yaml shape.
	Render func() string
}
