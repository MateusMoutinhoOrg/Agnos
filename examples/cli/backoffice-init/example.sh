# The backoffice-init example: add the admin backoffice to a project that has
# none of the layers it stands on
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos backoffice-init --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The backoffice's parts of RouteProps and of api.Config are files of
# their own, embedded by the generated aggregates beside them; the database
# is the backoffice's own; docs/Backoffice names the secret, TEST_BACKOFFICE_SECRET for
# a project named Test; .gitignore keeps the store out.
for dir in \
	sandbox/internal/routeprops \
	sandbox/internal/commands/backoffice \
	sandbox/internal/commands/middleware/backoffice_start_server \
	sandbox/internal/databases/backoffice_db \
	sandbox/internal/server/backoffice/backofficeauth \
	docs/Backoffice; do
	mkdir -p assert-dir/$dir
	cp -R test-dir/$dir/. assert-dir/$dir/
done

for file in \
	.gitignore \
	AgnosConfig/extensions.yaml \
	sandbox/api/config.go \
	sandbox/api/backofficeconfig.go; do
	mkdir -p assert-dir/$(dirname $file)
	cp test-dir/$file assert-dir/$file
done
