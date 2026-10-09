package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/databaseconf"
)

// DatabaseMethod is one method a table generates, spelled once for the three
// places that need it: api.go declares the function field, new.go closes over
// the function, methods.go writes its body and docs/Databases prints its
// signature. Params, Results and Args are the signature itself, so none of
// those four derives it a second time.
type DatabaseMethod struct {
	// Kind is which body methods.go writes: add, findById, findByKey, list,
	// page, count, update, remove, getLink, addSub or listSub.
	Kind string
	// Name is the exported Go name of the method.
	Name string
	// Table is the collection it reaches, as Database.Collection names it.
	Table string
	// Type is the Go record name of that table.
	Type string
	// Record is the record the method hands back — the table's own, the
	// target of a link, or the nested collection of a sub.
	Record string
	// Field is the declared field the method is born of, "" for the methods
	// every table generates.
	Field string
	// Params is the parameter list after the sandbox and the receiver.
	Params string
	// Results is the result list.
	Results string
	// Args is Params reduced to the names, for the closure of new.go.
	Args string
	// Help is the one line the doc comment and docs/Databases print.
	Help string
}

// DatabaseGoType is the Go type one declared field type is carried as. A link
// is the id of the record it points at, which is an int64 like any other id.
func DatabaseGoType(kind string) string {
	switch kind {
	case databaseconf.FieldInteger, databaseconf.FieldLink:
		return "int64"
	case databaseconf.FieldNumber:
		return "float64"
	}
	return "string"
}

// DatabaseMethodNames is DatabaseMethods reduced to the names, which is what
// `verify` checks methods_custom.go against.
func DatabaseMethodNames(sandbox *api.Sandbox, table databaseconf.Table) []string {
	methods := DatabaseMethods(sandbox, table)
	names := make([]string, 0, len(methods))
	for _, method := range methods {
		names = append(names, method.Name)
	}
	return names
}

// DatabaseMethods is every method one table generates, in the order api.go
// declares them, new.go wires them and docs/Databases prints them. The set is
// derived from the fields alone: a Find is born of a `key` field because that
// is the only one the database indexes, a Get of a `link`, and the pair
// Add/List of a nested collection.
//
// Each entry carries its own signature — Params, Results and the Args the
// closure of new.go calls through with — so the three templates spell one
// method the same way without deriving it three times.
func DatabaseMethods(sandbox *api.Sandbox, table databaseconf.Table) []DatabaseMethod {
	kind := GoIdentifier(sandbox, table.Name)
	item := kind + "Record"
	plural := Plural(kind)

	methods := []DatabaseMethod{
		databaseMethod(sandbox, "add", table, "Add"+kind, "props "+kind+"Input", "("+item+", error)", "props",
			"inserts one "+table.Name+" record"),
		databaseMethod(sandbox, "findById", table, "Find"+kind+"ById", "id int64", "("+item+", bool)", "id",
			"reads one "+table.Name+" record by its permanent id"),
	}

	for _, field := range table.Fields {
		if field.Type != databaseconf.FieldKey {
			continue
		}
		name := "Find" + kind + "By" + GoIdentifier(sandbox, field.Name)
		method := databaseMethod(sandbox, "findByKey", table, name, "value string", "("+item+", bool)", "value",
			"reads one "+table.Name+" record by its indexed "+field.Name)
		method.Field = field.Name
		methods = append(methods, method)
	}

	methods = append(methods,
		databaseMethod(sandbox, "list", table, "List"+plural, "filter "+kind+"Filter", "([]"+item+", error)", "filter",
			"reads every "+table.Name+" record the filter keeps"),
		databaseMethod(sandbox, "page", table, "List"+plural+"Page", "offset int, limit int", "([]"+item+", error)", "offset, limit",
			"reads up to limit "+table.Name+" records after the first offset, every one past them when limit is 0"),
		databaseMethod(sandbox, "count", table, "Count"+kind, "", "(int, error)", "",
			"is how many "+table.Name+" records are live"),
	)

	for _, field := range table.Fields {
		if field.Type == databaseconf.FieldObject {
			continue
		}
		name := "Set" + kind + GoIdentifier(sandbox, field.Name)
		method := databaseMethod(sandbox, "update", table, name,
			"id int64, value "+DatabaseGoType(field.Type), "error", "id, value",
			"writes a new "+field.Name+" on one "+table.Name+" record")
		method.Field = field.Name
		methods = append(methods, method)
	}

	methods = append(methods,
		databaseMethod(sandbox, "remove", table, "Remove"+kind, "id int64", "error", "id",
			"deletes one "+table.Name+" record and everything nested under it"),
	)

	for _, field := range table.Fields {
		if field.Type != databaseconf.FieldLink {
			continue
		}
		target := GoIdentifier(sandbox, field.Target)
		name := "Get" + kind + GoIdentifier(sandbox, field.Name)
		method := databaseMethod(sandbox, "getLink", table, name, "id int64", "("+target+"Record, bool)", "id",
			"resolves the "+field.Name+" link of one "+table.Name+" record to the "+field.Target+" it points at")
		method.Field = field.Name
		method.Record = target
		methods = append(methods, method)
	}

	for _, field := range table.Fields {
		if field.Type != databaseconf.FieldObject {
			continue
		}
		sub := kind + GoIdentifier(sandbox, field.Name)

		add := databaseMethod(sandbox, "addSub", table, "Add"+sub,
			"parentId int64, props "+sub+"Input", "("+sub+"Record, error)", "parentId, props",
			"inserts one "+field.Name+" record under one "+table.Name+" record")
		add.Field = field.Name
		add.Record = sub

		list := databaseMethod(sandbox, "listSub", table, "List"+Plural(sub),
			"parentId int64", "([]"+sub+"Record, error)", "parentId",
			"reads every "+field.Name+" record of one "+table.Name+" record")
		list.Field = field.Name
		list.Record = sub

		methods = append(methods, add, list)
	}

	return methods
}

// databaseMethod is one method entry with the keys every kind carries. Field
// and Record are filled by the caller for the kinds that name one.
func databaseMethod(sandbox *api.Sandbox, kind string, table databaseconf.Table, name string,
	params string, results string, args string, help string) DatabaseMethod {

	record := GoIdentifier(sandbox, table.Name)
	return DatabaseMethod{
		Kind:    kind,
		Name:    name,
		Table:   table.Name,
		Type:    record,
		Record:  record,
		Params:  params,
		Results: results,
		Args:    args,
		Help:    help,
	}
}

// Plural is the plural of one exported Go name, for the methods that hand
// back many records: Product -> Products, Category -> Categories, Address ->
// Addresses.
func Plural(name string) string {
	if name == "" {
		return name
	}
	last := name[len(name)-1]
	if last == 'y' && len(name) > 1 && !isVowel(name[len(name)-2]) {
		return name[:len(name)-1] + "ies"
	}
	for _, suffix := range []string{"s", "x", "z", "ch", "sh"} {
		if len(name) >= len(suffix) && name[len(name)-len(suffix):] == suffix {
			return name + "es"
		}
	}
	return name + "s"
}

// isVowel reports a lower-case ASCII vowel.
func isVowel(letter byte) bool {
	return letter == 'a' || letter == 'e' || letter == 'i' || letter == 'o' || letter == 'u'
}
