package migrations

import (
	api "{{.Module}}/sandbox/api"
	databasedeps "{{.Module}}/sandbox/deps/databasedeps"
)

// {{.Function}} is the hand-written step of version {{.Version}} of {{.DatabaseName}}: what
// brings the data of version {{.Previous}} to version {{.Version}} that no operation of
// migrations/v{{.Version}}.ops.yaml says — a seed, a field computed from others, a
// value split in two. migrate-database runs it on a copy of the database's
// folder, after those operations: from reads the copy with the schema of
// version {{.Previous}}, to with the schema of version {{.Version}}, record by record:
//
//	table, _ := to.Collection("<table>")
//	records, _ := table.ListAll()
//	for _, record := range records {
//		record.Update("<field>", "<value>")
//	}
//
// An error stops the migration, and the data in use stays exactly as it was.
//
// Written once by `{{.GeneratorName}}` when version {{.Version}} was opened, then yours: no
// build rewrites it.
func {{.Function}}(sandbox *api.Sandbox, from databasedeps.Database, to databasedeps.Database) error {
	return nil
}
