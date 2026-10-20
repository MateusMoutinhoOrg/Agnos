# The database-migration example: rename a field whose records are already
# stored, and let the project's migrate-database carry them across
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.
#
# An edit that changes what stored data holds opens a version of its own: the
# declaration as it was is frozen into migrations/v1.yaml, the operation that
# carries the data across goes to migrations/v2.ops.yaml, and database.yaml
# says version 2. No migration is written by hand: the command said it is a
# rename, so the values move to the new name, ids and all.

agnos start --path test-dir --project-name Test --module Test -q
agnos cli-init --path test-dir -q
agnos database-init --path test-dir -q
agnos add-database people --path test-dir -q
agnos add-table person --database people --path test-dir -q
agnos add-table-field nome --database people --table person --path test-dir -q

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
	if _, err := people.New(sandbox).AddPerson(people.PersonInput{Nome: input.Name}); err != nil {
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

# Version 1 of the binary stores two people. Its first command line finds no
# data at all, so the database-version middleware stamps it at version 1. The
# binary runs inside test-dir, where its data folder is.
(cd test-dir && go build -o app-v1.bin ./cmd/main)
(cd test-dir && ./app-v1.bin add-person Ana && ./app-v1.bin add-person Bia && ./app-v1.bin list-people)

# The rename. The handler that writes the field is the project's, so it is
# moved to the new name with it.
handler=test-dir/sandbox/internal/commands/people/add_person/handler.go
sed 's/Nome:/Name:/' $handler > $handler.new && mv $handler.new $handler
agnos set-table-field nome --rename name --database people --table person --path test-dir -q

# Version 2 of the binary refuses the data until it is migrated; then it reads
# every record under the new name. The data it replaced is kept aside.
(cd test-dir && go build -o app.bin ./cmd/main)
echo "== version 2, data at 1"
(cd test-dir && ./app.bin list-people 2>&1); echo "exit $?"
echo "== migrate-database"
(cd test-dir && ./app.bin migrate-database 2>&1); echo "exit $?"
(cd test-dir && ./app.bin list-people)

# The older binary refuses the newer data instead of writing the old shape
# over it.
echo "== version 1, data at 2"
(cd test-dir && ./app-v1.bin list-people 2>&1); echo "exit $?"

# Undoing a migration is moving the kept folder back. Data with no version
# record at all — written before versions existed — is read as version 1.
echo "== the kept data back, its version record gone"
rm -rf test-dir/data/people
mv test-dir/data/agnos-pre-migration/people-v1 test-dir/data/people
rm -rf test-dir/data/people/agnos-version
(cd test-dir && ./app.bin list-people 2>&1); echo "exit $?"
(cd test-dir && ./app.bin migrate-database 2>&1); echo "exit $?"
(cd test-dir && ./app.bin list-people)

# What result.yaml records: the declaration with its migrations folder and
# the four generated files, the registry of every database, and the two
# commands the database layer writes over the cli.
mkdir -p assert-dir/sandbox/internal/databases
cp -R test-dir/sandbox/internal/databases/. assert-dir/sandbox/internal/databases/
mkdir -p assert-dir/sandbox/internal/commands/middleware/database_version
cp -R test-dir/sandbox/internal/commands/middleware/database_version/. assert-dir/sandbox/internal/commands/middleware/database_version/
mkdir -p assert-dir/sandbox/internal/commands/database/migrate_database
cp -R test-dir/sandbox/internal/commands/database/migrate_database/. assert-dir/sandbox/internal/commands/database/migrate_database/
