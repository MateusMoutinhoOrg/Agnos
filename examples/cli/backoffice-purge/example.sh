# The backoffice-purge example: remove the admin backoffice backoffice-init wrote
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos backoffice-init --path TestDir -q

agnos backoffice-purge --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. routeprops is back to the project's own empty part, api.Config
# embeds UserConfig alone, routeslist holds what server-init and front-init
# wrote and nothing of the backoffice, and the key is off. The layers stay.
for dir in \
	sandbox/internal/routeprops \
	sandbox/internal/routeslist \
	sandbox/internal/commands; do
	mkdir -p AssertDir/$dir
	cp -R TestDir/$dir/. AssertDir/$dir/
done

for file in \
	AgnosConfig/extensions.yaml \
	sandbox/api/config.go; do
	mkdir -p AssertDir/$(dirname $file)
	cp TestDir/$file AssertDir/$file
done
