# The database example: declare a database and let agnos generate its methods
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.
#
# `database-init` installs the store as a remote dep, renders the databaseio
# package and turns the mechanic on. It scaffolds no database: which tables a
# project wants is a declaration, so every table below is declared by hand and
# api.go, new.go and methods.go are generated from that declaration alone.

agnos start --path TestDir --project-name Test --module Test -q

agnos database-init --path TestDir -q

agnos add-database app-database --prefix app --path TestDir -q

# Two tables, so the link below has somewhere to point.
agnos add-table user --database app-database --path TestDir -q
agnos add-table url --database app-database --path TestDir -q

# One field of every declared type. A `key` is the only indexed one, so it is
# the only one that generates a Find; a `link` generates a Get, a `database`
# the Add/List pair for the collection nested under each record, and every
# plain field an Update and a place in the table's filtrage.
agnos add-table-field email --database app-database --table user --type key --required --path TestDir -q
agnos add-table-field name --database app-database --table user --type string --path TestDir -q

agnos add-table-field alias --database app-database --table url --type key --required --path TestDir -q
agnos add-table-field link --database app-database --table url --type string --required --path TestDir -q
agnos add-table-field redirects --database app-database --table url --type int --path TestDir -q
agnos add-table-field score --database app-database --table url --type float --path TestDir -q
agnos add-table-field owner --database app-database --table url --type link --target user --path TestDir -q
agnos add-table-field visits --database app-database --table url --type database --path TestDir -q

# A field of the nested collection, which is what --parent names.
agnos add-table-field agent --database app-database --table url --parent visits --type string --path TestDir -q
agnos add-table-field at --database app-database --table url --parent visits --type int --path TestDir -q

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The declaration and the three files generated from it, the package
# they share, the page the build wrote, and the key that says the mechanic is
# on. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/databases
cp -R TestDir/sandbox/internal/databases/. AssertDir/sandbox/internal/databases/
mkdir -p AssertDir/sandbox/internal/generated/databaseio
cp -R TestDir/sandbox/internal/generated/databaseio/. AssertDir/sandbox/internal/generated/databaseio/
mkdir -p AssertDir/docs/Databases
cp -R TestDir/docs/Databases/. AssertDir/docs/Databases/
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
