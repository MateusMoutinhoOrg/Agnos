package main

import (
	"os"
	"path/filepath"

	"github.com/MateusMoutinhoOrg/Agnos/adapters/bindings/standard"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// The database example: declare a database and let agnos generate its methods.
//
// It calls the same actions `agnos database-init`, `agnos add-database`,
// `agnos add-table` and `agnos add-table-field` call, and writes only inside
// test-dir. DatabaseInit installs the store as a remote dep and turns the
// mechanic on; it scaffolds no database, because which tables a project wants
// is a declaration — so every table below is declared here and api.go, new.go
// and methods.go are generated from that declaration alone.
func main() {

	deps := standard.New()    // every adapter lib bound
	lib := sandbox.New(&deps) // *api.Sandbox
	module := "Test"

	if err := lib.Actions.Start(api.StartProps{
		Path:        "test-dir",
		ProjectName: "Test",
		Module:      &module,
	}); err != nil {
		panic(err)
	}

	if err := lib.Actions.DatabaseInit(api.DatabaseInitProps{Path: "test-dir"}); err != nil {
		panic(err)
	}

	if err := lib.Actions.AddDatabase(api.AddDatabaseProps{Path: "test-dir", Name: "app-database", KeyPrefix: "app"}); err != nil {
		panic(err)
	}

	// Two tables, so the link below has somewhere to point.
	for _, table := range []string{"user", "url"} {
		if err := lib.Actions.AddTable(api.AddTableProps{Path: "test-dir", Database: "app-database", Name: table}); err != nil {
			panic(err)
		}
	}

	// One field of every declared type. A `key` is the only indexed one, so it
	// is the only one that generates a Find; a `link` generates a Get, a
	// `database` the Add/List pair for the collection nested under each
	// record, and every plain field an Update and a place in the table's
	// filter. A field declared with a Parent lands inside that nested
	// collection instead of on the table.
	fields := []api.AddTableFieldProps{
		{Table: "user", Name: "email", Type: "key", Required: true},
		{Table: "user", Name: "name", Type: "string"},
		{Table: "user", Name: "avatar", Type: "bytes"},
		{Table: "url", Name: "alias", Type: "key", Required: true},
		{Table: "url", Name: "link", Type: "string", Required: true},
		{Table: "url", Name: "redirects", Type: "integer"},
		{Table: "url", Name: "score", Type: "number"},
		{Table: "url", Name: "owner", Type: "link", Target: "user"},
		{Table: "url", Name: "visits", Type: "object"},
		{Table: "url", Parent: "visits", Name: "agent", Type: "string"},
		{Table: "url", Parent: "visits", Name: "at", Type: "integer"},
	}
	for _, field := range fields {
		field.Path = "test-dir"
		field.Database = "app-database"
		if err := lib.Actions.AddTableField(field); err != nil {
			panic(err)
		}
	}

	// What result.yaml records: the same set the cli side copies.
	for _, dir := range []string{
		"sandbox/internal/databases",
		"sandbox/deps/OpinionatedAgnosDatabase",
		"docs/Databases",
	} {
		copyTree("test-dir/"+dir, "assert-dir/"+dir)
	}
	copyFile("test-dir/AgnosConfig/extensions.yaml", "assert-dir/AgnosConfig/extensions.yaml")
}

// copyTree copies every file under source into dest, keeping the place each
// one holds in the tree.
func copyTree(source string, dest string) {
	err := filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}

		return write(filepath.Join(dest, relative), path)
	})
	if err != nil {
		panic(err)
	}
}

// copyFile copies one file, creating the directory it lands in.
func copyFile(source string, dest string) {
	if err := write(dest, source); err != nil {
		panic(err)
	}
}

// write is the one copy both of them go through.
func write(target string, source string) error {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return os.WriteFile(target, content, 0o644)
}
