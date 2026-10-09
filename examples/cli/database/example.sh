# The database example: declare a database and let agnos generate its methods
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.
#
# `database-init` installs the store as a remote dep and the
# OpinionatedAgnosDatabase lib over it, and turns the mechanic on. It scaffolds no database: which tables a
# project wants is a declaration, so every table below is declared by hand and
# api.go, new.go and methods.go are generated from that declaration alone.

agnos start --path test-dir --project-name Test --module Test -q

agnos database-init --path test-dir -q

agnos add-database app-database --key-prefix app --path test-dir -q

# Two tables, so the link below has somewhere to point.
agnos add-table user --database app-database --path test-dir -q
agnos add-table url --database app-database --path test-dir -q

# One field of every declared type. A `key` is the only indexed one, so it is
# the only one that generates a Find; a `link` generates a Get, a `database`
# the Add/List pair for the collection nested under each record, and every
# plain field an Update and a place in the table's filter.
agnos add-table-field email --database app-database --table user --type key --required --path test-dir -q
agnos add-table-field name --database app-database --table user --type string --path test-dir -q
agnos add-table-field avatar --database app-database --table user --type bytes --path test-dir -q

agnos add-table-field alias --database app-database --table url --type key --required --path test-dir -q
agnos add-table-field link --database app-database --table url --type string --required --path test-dir -q
agnos add-table-field redirects --database app-database --table url --type integer --path test-dir -q
agnos add-table-field score --database app-database --table url --type number --path test-dir -q
agnos add-table-field owner --database app-database --table url --type link --target user --path test-dir -q
agnos add-table-field visits --database app-database --table url --type object --path test-dir -q

# A field of the nested collection, which is what --parent names.
agnos add-table-field agent --database app-database --table url --parent visits --type string --path test-dir -q
agnos add-table-field at --database app-database --table url --parent visits --type integer --path test-dir -q

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The declaration and the three files generated from it, the
# contract of the OpinionatedAgnosDatabase lib they share, the page the build
# wrote, and the key that says the mechanic is on. The lib side copies the
# same set.
mkdir -p assert-dir/sandbox/internal/databases
cp -R test-dir/sandbox/internal/databases/. assert-dir/sandbox/internal/databases/
mkdir -p assert-dir/sandbox/deps/OpinionatedAgnosDatabase
cp -R test-dir/sandbox/deps/OpinionatedAgnosDatabase/. assert-dir/sandbox/deps/OpinionatedAgnosDatabase/
mkdir -p assert-dir/docs/Databases
cp -R test-dir/docs/Databases/. assert-dir/docs/Databases/
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
