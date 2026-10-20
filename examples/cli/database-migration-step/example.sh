# The database-migration-step example: a version no editing command opens,
# carried by a step written by hand, and a failing step that leaves the data
# in use exactly as it was
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.
#
# add-database-migration opens the next version with no operation of its own
# and writes its step, migrations/v2.go, once: V2 is handed the copy a
# migration works on, read with the schema of version 1 (from) and of version
# 2 (to). It reaches the data through databasedeps alone, so it compiles
# whatever the schema becomes.

agnos start --path test-dir --project-name Test --module Test -q
agnos cli-init --path test-dir -q
agnos database-init --path test-dir -q
agnos add-database people --path test-dir -q
agnos add-table person --database people --path test-dir -q
agnos add-table-field name --database people --table person --type key --required --path test-dir -q

# Two commands of the project's own: one stores a person, one prints them all.
agnos add-command add-person --summary "Store one person" --category People --path test-dir -q
agnos add-arg name --command add-person --start 1 --end 1 --required --path test-dir -q
agnos add-command list-people --summary "Print every person" --category People --path test-dir -q
cat > test-dir/sandbox/internal/commands/people/add_person/handler.go <<'GO'
package add_person

import (
	"Test/sandbox/api"
	"Test/sandbox/internal/commandprops"
	"Test/sandbox/internal/databases/people"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	if _, err := people.New(sandbox).AddPerson(people.PersonInput{Name: input.Name}); err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	response.SetStatus(api.ExitOk)
	return nil
}
GO
cat > test-dir/sandbox/internal/commands/people/list_people/handler.go <<'GO'
package list_people

import (
	"Test/sandbox/api"
	"Test/sandbox/internal/commandprops"
	"Test/sandbox/internal/databases/people"
)

func Handle(sandbox *api.Sandbox, props *commandprops.CommandProps, input *Input, response *api.CommandResponse) error {
	records, err := people.New(sandbox).ListPersons(people.PersonFilter{})
	if err != nil {
		return sandbox.Deps.OpinionatedAgnosCli.Fail(api.ExitFailure, "", err.Error())
	}
	for _, record := range records {
		response.Printf("%+v\n", record)
	}
	response.SetStatus(api.ExitOk)
	return nil
}
GO

# A required field changes what every stored record has to hold, so declaring
# one opened version 2. The first command line finds no data and stamps it
# there. The binary runs inside test-dir, where its data folder is.
(cd test-dir && go build -o app.bin ./cmd/main && ./app.bin add-person Ana)

# Version 3, a seed. The step first refuses: the migration stops, and the
# folder in use is byte for byte what it was.
agnos add-database-migration people --path test-dir -q
cat > test-dir/sandbox/internal/databases/people/migrations/v3.go <<'GO'
package migrations

import (
	api "Test/sandbox/api"
	databasedeps "Test/sandbox/deps/databasedeps"
)

func V3(sandbox *api.Sandbox, from databasedeps.Database, to databasedeps.Database) error {
	return sandbox.Deps.StdDeps.Errorf("the seed is not written yet")
}
GO
agnos build --path test-dir -q
(cd test-dir && go build -o app.bin ./cmd/main)
cp -R test-dir/data/people test-dir/people-before
echo "== a step that fails"
(cd test-dir && ./app.bin migrate-database 2>&1); echo "exit $?"
diff -r test-dir/people-before test-dir/data/people && echo "data/people is byte for byte what it was"
ls test-dir/data

# The step written: it seeds the people a fresh install starts with, through
# the schema of the version it brings the data to.
cat > test-dir/sandbox/internal/databases/people/migrations/v3.go <<'GO'
package migrations

import (
	api "Test/sandbox/api"
	databasedeps "Test/sandbox/deps/databasedeps"
)

func V3(sandbox *api.Sandbox, from databasedeps.Database, to databasedeps.Database) error {
	people, _ := to.Collection("person")
	for _, name := range []string{"Root", "Guest"} {
		if _, failure := people.Insert(map[string]any{"name": name}); failure != nil {
			return sandbox.Deps.StdDeps.Errorf("seeding %s: %s", name, failure.Message)
		}
	}
	return nil
}
GO
agnos build --path test-dir -q
(cd test-dir && go build -o app.bin ./cmd/main)
echo "== the step written"
(cd test-dir && ./app.bin migrate-database 2>&1); echo "exit $?"
(cd test-dir && ./app.bin list-people)

# What result.yaml records: the declaration, its migrations folder — the step
# among them, the project's from the moment it was written — and the
# generated.migrations.go that runs it.
mkdir -p assert-dir/sandbox/internal/databases/people
cp -R test-dir/sandbox/internal/databases/people/. assert-dir/sandbox/internal/databases/people/
