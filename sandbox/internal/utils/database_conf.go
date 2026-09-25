package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// The database helpers below mirror the route ones of route_conf.go, one for
// one: a database is declared, edited and looked up exactly the way a route
// is, only against specs.yaml instead of route.yaml.

// DatabasesDir holds one declared database per sub-directory, the database
// layer's mirror of sandbox/internal/routeslist.
const DatabasesDir = "sandbox/internal/databases"

// DatabaseSpecsFile is the declaration every database directory carries, the
// whole of what its api.go, new.go and methods.go are generated from.
const DatabaseSpecsFile = "specs.yaml"

// DatabaseCustomFile is the one file of a database package agnos never writes
// and never reads: the escape hatch for a query specs.yaml cannot describe.
const DatabaseCustomFile = "methods_custom.go"

// DatabaseIdentifier normalizes a user-typed database name into its canonical
// spelling, the same alphabet a command and a route are held to.
func DatabaseIdentifier(sandbox *api.Sandbox, name string) string {
	return CommandIdentifier(sandbox, name)
}

// ValidateDatabaseName reports whether a user-typed database name normalizes
// to a usable identifier. It becomes a directory name and a Go package clause,
// so the same alphabet a route name is held to applies.
func ValidateDatabaseName(sandbox *api.Sandbox, name string) error {
	identifier := DatabaseIdentifier(sandbox, name)
	if identifier == "" {
		return sandbox.Deps.Std.Errorf("a database needs a name")
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return sandbox.Deps.Std.Errorf("invalid database name %q: a database name must start with a lowercase letter", name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return sandbox.Deps.Std.Errorf(
				"invalid database name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes the directory %s/%s and a Go package name)",
				name, DatabasesDir, DatabasePackage(sandbox, name))
		}
	}
	return nil
}

// DatabasePackage is the Go package / directory name for a database: the
// identifier with dashes turned into underscores ("app-database" ->
// "app_database").
func DatabasePackage(sandbox *api.Sandbox, name string) string {
	return sandbox.Deps.Stringsdeps.ReplaceAll(DatabaseIdentifier(sandbox, name), "-", "_")
}

// DatabaseDir is the project-relative directory holding a database package.
func DatabaseDir(sandbox *api.Sandbox, name string) string {
	return DatabasesDir + "/" + DatabasePackage(sandbox, name)
}

// DatabaseConfPath is the project-relative path of a database's specs.yaml.
func DatabaseConfPath(sandbox *api.Sandbox, name string) string {
	return DatabaseDir(sandbox, name) + "/" + DatabaseSpecsFile
}

// DatabaseCustomPath is the project-relative path of a database's hand-written
// half.
func DatabaseCustomPath(sandbox *api.Sandbox, name string) string {
	return DatabaseDir(sandbox, name) + "/" + DatabaseCustomFile
}

// LoadDatabaseConf reads and parses
// sandbox/internal/databases/<name>/specs.yaml.
func LoadDatabaseConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string) (*databaseconf.DatabaseConf, error) {
	if err := ValidateDatabaseName(sandbox, name); err != nil {
		return nil, err
	}
	content, err := io.ReadFile(DatabaseConfPath(sandbox, name))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("database %q not found in %s", DatabaseIdentifier(sandbox, name), DatabaseDir(sandbox, name))
	}
	conf, err := databaseconf.New(sandbox, string(content))
	if err != nil {
		return nil, sandbox.Deps.Std.Errorf("databases/%s/%s: %w", DatabasePackage(sandbox, name), DatabaseSpecsFile, err)
	}
	return conf, nil
}

// SaveDatabaseConf renders conf back over
// sandbox/internal/databases/<name>/specs.yaml.
func SaveDatabaseConf(sandbox *api.Sandbox, io *smartio.SmartIO, name string, conf *databaseconf.DatabaseConf) error {
	return io.WriteFileOverwrite(DatabaseConfPath(sandbox, name), []byte(conf.Render()))
}

// DatabaseName normalizes a table or field name. Both become a Go identifier
// and a key of the stored record, so they are held to the same alphabet the
// database's own name is.
func DatabaseName(sandbox *api.Sandbox, name string) string {
	return DatabaseIdentifier(sandbox, name)
}

// ValidateDatabaseMember reports whether a table or field name normalizes to
// something that can be both a key of a stored record and a piece of a
// generated Go identifier. kind names what is being declared, so the message
// says "table" or "field" rather than a word covering both.
func ValidateDatabaseMember(sandbox *api.Sandbox, kind string, name string) error {
	identifier := DatabaseName(sandbox, name)
	if identifier == "" {
		return sandbox.Deps.Std.Errorf("a %s needs a name", kind)
	}
	if identifier[0] < 'a' || identifier[0] > 'z' {
		return sandbox.Deps.Std.Errorf("invalid %s name %q: it must start with a lowercase letter", kind, name)
	}
	for _, letter := range identifier {
		valid := (letter >= 'a' && letter <= 'z') ||
			(letter >= '0' && letter <= '9') ||
			letter == '-'
		if !valid {
			return sandbox.Deps.Std.Errorf(
				"invalid %s name %q: only letters, digits, spaces, dashes and underscores are allowed (it becomes part of a generated Go identifier)",
				kind, name)
		}
	}
	return nil
}

// ExportedName is a declared name as the exported Go identifier generated from
// it: every dash- or underscore-separated part capitalized ("app-database" ->
// "AppDatabase").
func ExportedName(sandbox *api.Sandbox, raw string) string {
	parts := sandbox.Deps.Stringsdeps.FieldsFunc(raw, func(r rune) bool { return r == '-' || r == '_' })
	out := ""
	for _, part := range parts {
		if part == "" {
			continue
		}
		out += sandbox.Deps.Stringsdeps.ToUpper(part[:1]) + part[1:]
	}
	return out
}

// FindDatabaseTable returns the index of the table named name in tables, or -1.
func FindDatabaseTable(sandbox *api.Sandbox, tables []databaseconf.Table, name string) int {
	key := DatabaseName(sandbox, name)
	for i, table := range tables {
		if table.Name == key {
			return i
		}
	}
	return -1
}

// FindDatabaseField returns the index of the field named name in fields, or -1.
func FindDatabaseField(sandbox *api.Sandbox, fields []databaseconf.Field, name string) int {
	key := DatabaseName(sandbox, name)
	for i, field := range fields {
		if field.Name == key {
			return i
		}
	}
	return -1
}

// RemoveDatabaseField drops the field at index from fields.
func RemoveDatabaseField(fields []databaseconf.Field, index int) []databaseconf.Field {
	out := make([]databaseconf.Field, 0, len(fields)-1)
	out = append(out, fields[:index]...)
	out = append(out, fields[index+1:]...)
	return out
}

// RemoveDatabaseTable drops the table at index from tables.
func RemoveDatabaseTable(tables []databaseconf.Table, index int) []databaseconf.Table {
	out := make([]databaseconf.Table, 0, len(tables)-1)
	out = append(out, tables[:index]...)
	out = append(out, tables[index+1:]...)
	return out
}

// DatabaseFieldType normalizes a user-typed field type, refusing one no field
// may declare.
func DatabaseFieldType(sandbox *api.Sandbox, raw string) (string, error) {
	kind := sandbox.Deps.Stringsdeps.ToLower(sandbox.Deps.Stringsdeps.TrimSpace(raw))
	if kind == "" {
		return databaseconf.FieldString, nil
	}
	if !databaseconf.IsFieldType(kind) {
		return "", sandbox.Deps.Std.Errorf("unknown field type %q (use one of %s)",
			raw, sandbox.Deps.Stringsdeps.Join(databaseconf.FieldTypes, ", "))
	}
	return kind, nil
}

// DatabaseFieldsAt is the field list one --table, or one --parent collection
// inside it, declares — plus the label error messages name that place by. It
// is how every editor of a specs.yaml finds the list it is about to change, so
// "a field of a table" and "a field of a nested collection" are one code path.
func DatabaseFieldsAt(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf, table string, parent string) ([]databaseconf.Field, string, error) {
	name := DatabaseName(sandbox, table)
	index := FindDatabaseTable(sandbox, conf.Tables, name)
	if index < 0 {
		return nil, "", sandbox.Deps.Std.Errorf("this database declares no table %q", name)
	}

	if DatabaseName(sandbox, parent) == "" {
		return conf.Tables[index].Fields, name, nil
	}

	nested := DatabaseName(sandbox, parent)
	at := FindDatabaseField(sandbox, conf.Tables[index].Fields, nested)
	if at < 0 {
		return nil, "", sandbox.Deps.Std.Errorf("table %q declares no field %q", name, nested)
	}
	if conf.Tables[index].Fields[at].Type != databaseconf.FieldDatabase {
		return nil, "", sandbox.Deps.Std.Errorf(
			"%s.%s is a %s field, not a nested collection: only a `database` field holds fields of its own",
			name, nested, conf.Tables[index].Fields[at].Type)
	}

	return conf.Tables[index].Fields[at].Fields, name + "." + nested, nil
}

// SetDatabaseFieldsAt writes a field list back where DatabaseFieldsAt read it.
func SetDatabaseFieldsAt(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf, table string, parent string, fields []databaseconf.Field) error {
	name := DatabaseName(sandbox, table)
	index := FindDatabaseTable(sandbox, conf.Tables, name)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("this database declares no table %q", name)
	}

	if DatabaseName(sandbox, parent) == "" {
		conf.Tables[index].Fields = fields
		return nil
	}

	at := FindDatabaseField(sandbox, conf.Tables[index].Fields, parent)
	if at < 0 {
		return sandbox.Deps.Std.Errorf("table %q declares no field %q", name, DatabaseName(sandbox, parent))
	}
	conf.Tables[index].Fields[at].Fields = fields
	return nil
}

// NewDatabaseField builds a databaseconf.Field from the raw values typed on
// the command line, holding it to the rules the generator cannot render
// around: a link names a table of this database, a nested collection only
// exists one level deep and is never required, and nothing else carries a
// target.
func NewDatabaseField(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf, props api.DatabaseFieldProps) (databaseconf.Field, error) {
	field := databaseconf.Field{
		Name:     DatabaseName(sandbox, props.Name),
		Required: props.Required,
		Target:   DatabaseName(sandbox, props.Target),
		Fields:   []databaseconf.Field{},
	}

	if err := ValidateDatabaseMember(sandbox, "field", props.Name); err != nil {
		return field, err
	}

	kind, err := DatabaseFieldType(sandbox, props.Type)
	if err != nil {
		return field, err
	}
	field.Type = kind

	if err := CheckDatabaseField(sandbox, conf, field, props.Parent != ""); err != nil {
		return field, err
	}

	return field, nil
}

// CheckDatabaseField is the rule set both add-table-field and set-table-field
// hold a field to, so the two commands refuse the same declarations for the
// same reasons.
func CheckDatabaseField(sandbox *api.Sandbox, conf *databaseconf.DatabaseConf, field databaseconf.Field, nested bool) error {
	if field.Name == databaseIdFieldName {
		return sandbox.Deps.Std.Errorf("a field cannot be named %q: every record already carries its permanent Id", databaseIdFieldName)
	}

	if field.Type == databaseconf.FieldLink {
		if field.Target == "" {
			return sandbox.Deps.Std.Errorf("a link needs --target, the table it points at")
		}
		if FindDatabaseTable(sandbox, conf.Tables, field.Target) < 0 {
			return sandbox.Deps.Std.Errorf("--target %q is not a table of this database", field.Target)
		}
	} else if field.Target != "" {
		return sandbox.Deps.Std.Errorf("--target belongs to `--type link`: a %s field points at no table", field.Type)
	}

	if field.Type == databaseconf.FieldDatabase {
		if nested {
			return sandbox.Deps.Std.Errorf("a nested collection holds plain fields only: only one level of `--type database` is generated")
		}
		if field.Required {
			return sandbox.Deps.Std.Errorf("a nested collection cannot be required: it is never written by an insert")
		}
	}

	return nil
}

// databaseIdFieldName is the one field name a table may not declare, because
// every generated record already carries the stored id under it.
const databaseIdFieldName = "id"

// DatabaseFieldClearKeys is every key --clear may take off a declared field.
var DatabaseFieldClearKeys = []string{"required", "target"}

// DatabaseFieldEditEmpty reports a set-table-field that was given nothing to
// change, which would otherwise rewrite the declaration to exactly what it
// already says.
func DatabaseFieldEditEmpty(sandbox *api.Sandbox, props api.DatabaseFieldEditProps) bool {
	return sandbox.Deps.Stringsdeps.TrimSpace(props.Rename) == "" &&
		sandbox.Deps.Stringsdeps.TrimSpace(props.Type) == "" &&
		sandbox.Deps.Stringsdeps.TrimSpace(props.Target) == "" &&
		!props.Required && len(props.Clear) == 0
}

// DatabaseFieldEdited rebuilds one declared field with the changes applied,
// holding the result to every rule NewDatabaseField holds a new one to. An
// edited field and a declared one are the same bytes, because the two go
// through one constructor.
func DatabaseFieldEdited(sandbox *api.Sandbox, current databaseconf.Field,
	conf *databaseconf.DatabaseConf, props api.DatabaseFieldEditProps) (databaseconf.Field, error) {

	cleared, err := RouteClearSet(sandbox, props.Clear, DatabaseFieldClearKeys)
	if err != nil {
		return databaseconf.Field{}, err
	}

	built := api.DatabaseFieldProps{
		Parent:   props.Parent,
		Name:     current.Name,
		Type:     current.Type,
		Required: current.Required,
		Target:   current.Target,
	}

	if cleared["required"] {
		built.Required = false
	}
	if cleared["target"] {
		built.Target = ""
	}

	if rename := DatabaseName(sandbox, props.Rename); rename != "" {
		built.Name = rename
	}
	if kind := sandbox.Deps.Stringsdeps.TrimSpace(props.Type); kind != "" {
		built.Type = kind
	}
	if target := DatabaseName(sandbox, props.Target); target != "" {
		built.Target = target
	}
	if props.Required {
		built.Required = true
	}

	field, err := NewDatabaseField(sandbox, conf, built)
	if err != nil {
		return field, err
	}

	// A field that was already a nested collection keeps what it holds: the
	// editors of those fields are add-table-field --parent and its inverse,
	// not this one.
	if field.Type == databaseconf.FieldDatabase && current.Type == databaseconf.FieldDatabase {
		field.Fields = current.Fields
	}

	return field, nil
}
