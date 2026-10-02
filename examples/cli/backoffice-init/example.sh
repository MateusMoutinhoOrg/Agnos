# The backoffice-init example: add the admin backoffice to a project that has
# none of the layers it stands on
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos backoffice-init --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The backoffice's parts of RouteProps and of api.Config are files of
# their own, embedded by the generated aggregates beside them; the database
# is the backoffice's own; docs/Backoffice names the secret, TEST_SECRET for
# a project named Test; .gitignore keeps the store out.
for dir in \
	sandbox/internal/routeprops \
	sandbox/internal/commands/backoffice \
	sandbox/internal/databases/backofficedb \
	sandbox/internal/server/backoffice/backofficeauth \
	docs/Backoffice; do
	mkdir -p AssertDir/$dir
	cp -R TestDir/$dir/. AssertDir/$dir/
done

for file in \
	.gitignore \
	AgnosConfig/extensions.yaml \
	sandbox/api/config.go \
	sandbox/api/userconfig_backoffice.go; do
	mkdir -p AssertDir/$(dirname $file)
	cp TestDir/$file AssertDir/$file
done
