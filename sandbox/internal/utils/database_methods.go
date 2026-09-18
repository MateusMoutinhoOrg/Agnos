package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/databaseconf"
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
	// Table is the collection it reaches, as GetSchema names it.
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
	case databaseconf.FieldInt, databaseconf.FieldLink:
		return "int64"
	case databaseconf.FieldFloat:
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
	kind := ExportedName(sandbox, table.Name)
	item := kind + "Item"

	methods := []DatabaseMethod{
		databaseMethod(sandbox, "add", table, "Add"+kind, "props "+kind+"New", "("+item+", error)", "props",
			"inserts one "+table.Name+" record"),
		databaseMethod(sandbox, "findById", table, "Find"+kind+"ById", "id int64", "("+item+", bool)", "id",
			"reads one "+table.Name+" record by its permanent id"),
	}

	for _, field := range table.Fields {
		if field.Type != databaseconf.FieldKey {
			continue
		}
		name := "Find" + kind + "By" + ExportedName(sandbox, field.Name)
		method := databaseMethod(sandbox, "findByKey", table, name, "value string", "("+item+", bool)", "value",
			"reads one "+table.Name+" record by its indexed "+field.Name)
		method.Field = field.Name
		methods = append(methods, method)
	}

	methods = append(methods,
		databaseMethod(sandbox, "list", table, "List"+kind, "filtrage "+kind+"Filtrage", "([]"+item+", error)", "filtrage",
			"reads every "+table.Name+" record the filtrage keeps"),
		databaseMethod(sandbox, "page", table, "Page"+kind, "position int, chunk int", "([]"+item+", error)", "position, chunk",
			"reads one page of "+table.Name+" records, counted from 1"),
		databaseMethod(sandbox, "count", table, "Count"+kind, "", "(int, error)", "",
			"is how many "+table.Name+" records are live"),
	)

	for _, field := range table.Fields {
		if field.Type == databaseconf.FieldDatabase {
			continue
		}
		name := "Update" + kind + ExportedName(sandbox, field.Name)
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
		target := ExportedName(sandbox, field.Target)
		name := "Get" + kind + ExportedName(sandbox, field.Name)
		method := databaseMethod(sandbox, "getLink", table, name, "id int64", "("+target+"Item, bool)", "id",
			"resolves the "+field.Name+" link of one "+table.Name+" record to the "+field.Target+" it points at")
		method.Field = field.Name
		method.Record = target
		methods = append(methods, method)
	}

	for _, field := range table.Fields {
		if field.Type != databaseconf.FieldDatabase {
			continue
		}
		sub := ExportedName(sandbox, field.Name)

		add := databaseMethod(sandbox, "addSub", table, "Add"+kind+sub,
			"parent_id int64, props "+sub+"New", "("+sub+"Item, error)", "parent_id, props",
			"inserts one "+field.Name+" record under one "+table.Name+" record")
		add.Field = field.Name
		add.Record = sub

		list := databaseMethod(sandbox, "listSub", table, "List"+kind+sub,
			"parent_id int64", "([]"+sub+"Item, error)", "parent_id",
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

	record := ExportedName(sandbox, table.Name)
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
