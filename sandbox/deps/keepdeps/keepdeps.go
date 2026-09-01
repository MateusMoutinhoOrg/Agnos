package keepdeps

// This package is this library's *copy* of the embedded Keep schema
// database library's public api. The sandbox may not import the embedded
// library — that would be a third-party import — so it restates the shape
// it needs here, field for field. The adapter, which lives outside the
// sandbox, is what fills these structs from the real library.
//
// Copying is cheap precisely because the embedded library exposes structs
// of function fields instead of interfaces: an adapter assigns the real
// library's fields straight into the copy. Where a field hands back
// another api struct — Lib.NewDatabase, KeepDatabase.GetSchema,
// SchemaInstance.NewItem — the adapter wraps it in a closure that copies
// the returned struct too, so nothing of the embedded library ever reaches
// the sandbox.

// Field types, reported by Item.Type. The values mirror the embedded
// library's constants, so they can be assigned across without mapping.
const (
	// Key is a unique, indexed string field.
	Key = iota
	// Int is a plain integer field.
	Int
	// Database is a nested collection of records.
	Database
)

// Failure causes, reported by Error.Type. The values mirror the embedded
// library's constants.
const (
	KeyConflict = iota
	NotFound
	MissingField
	InvalidField
	Internal
)

// Item mirrors the embedded library's api.Item — one field of a schema.
// It carries no behavior, so it is plain data on both sides.
type Item struct {
	// Name is the field's name, as used in the fields map.
	Name string
	// Type is one of Key, Int, or Database.
	Type int
	// Required reports whether a record must provide this field.
	Required bool
	// Itens are the nested fields, for a Database field; nil otherwise.
	Itens []Item
}

// Schema mirrors the embedded library's api.Schema — one collection of
// records and its fields.
type Schema struct {
	// Name is the collection's name, as used in KeepDatabase.GetSchema.
	Name string
	// Itens are the fields each record of the collection can hold.
	Itens []Item
}

// Props mirrors the embedded library's api.Props — the declarative
// description a database is created from.
type Props struct {
	// Path is the prefix every key of the database is written under.
	Path string
	// Schemas are the collections the database holds.
	Schemas []Schema
}

// Error mirrors the embedded library's api.Error — one failure reported
// by a database operation. A nil *Error means success.
type Error struct {
	// Type is one of KeyConflict, NotFound, MissingField, InvalidField,
	// or Internal.
	Type int
	// Key is the field the failure involves.
	Key string
	// KeyValue is the value the failure involves, when relevant.
	KeyValue any
	// Message is the human-readable description of the failure.
	Message string
}

// SchemaItem mirrors the embedded library's api.SchemaItem — one record
// of a collection. The embedded library's Deps field is deliberately not
// copied: the sandbox has no business reading the dependencies the
// embedded library was wired with.
type SchemaItem struct {
	// Items are the schema fields this record's collection declares.
	Items []Item
	// Prefix is the collection's key prefix.
	Prefix string
	// Id is the record's permanent, never-reused identifier.
	Id int64
	// Get returns the typed value stored for a field.
	Get func(fieldName string) (any, *Error)
	// Update writes a new value for a field, re-indexing it when the
	// field is a unique key.
	Update func(fieldName string, value any) *Error
	// Remove deletes the record and everything nested under it,
	// returning nil on success.
	Remove func() *Error
	// CheckKeysPresence reports whether every named field has a stored
	// value for this record.
	CheckKeysPresence func(keys []string) bool
	// ListAll returns every record of a nested (Database) field.
	ListAll func(fieldName string) []SchemaItem
	// NewSubItem inserts a record into a nested (Database) field.
	NewSubItem func(fieldName string, fields map[string]any) (SchemaItem, *Error)
	// String renders the record's plain fields.
	String func() string
}

// SchemaInstance mirrors the embedded library's api.SchemaInstance — one
// collection of records.
type SchemaInstance struct {
	// Items are the fields each record of the collection can hold.
	Items []Item
	// Prefix is the collection's key prefix.
	Prefix string
	// NewItem inserts a record, validating the fields against the schema.
	NewItem func(fields map[string]any) (SchemaItem, *Error)
	// FindByKey looks a record up through a unique Key field. ok is false
	// when the field is not an indexed key or no record matches.
	FindByKey func(key string, keyValue any) (SchemaItem, bool)
	// ListAll returns every record of the collection.
	ListAll func() ([]SchemaItem, *Error)
	// List returns up to chunk records starting at position (1-based).
	List func(position int, chunk int) ([]SchemaItem, *Error)
}

// KeepDatabase mirrors the embedded library's api.KeepDatabase — a
// database bound to a Props description.
type KeepDatabase struct {
	// Props is the description the database was created from.
	Props Props
	// GetSchema returns the collection with the given name. ok is false
	// when the database declares no schema under that name.
	GetSchema func(name string) (SchemaInstance, bool)
}

// Lib mirrors the embedded Keep library's api.Lib — a schema database over
// an injected single-key storage backend. It is injected whole as the
// Deps.KeepLib field — the same mechanic as requestdeps.Lib — and the
// adapter, which lives outside the sandbox, fills it over the embedded Keep
// library. A database is rooted at Props.Path, so no per-base-path
// constructor is needed.
type Lib struct {
	// NewDatabase creates a database from a Props description, rooted at
	// Props.Path.
	NewDatabase func(props Props) KeepDatabase
}
