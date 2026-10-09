# The backoffice-purge example: remove the admin backoffice backoffice-init wrote
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos backoffice-init --path test-dir -q

agnos backoffice-purge --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. routeprops is back to the project's own empty part, api.Config
# embeds ProjectConfig alone, routes holds what server-init and front-init
# wrote and nothing of the backoffice, and the key is off. The layers stay.
for dir in \
	sandbox/internal/routeprops \
	sandbox/internal/routes \
	sandbox/internal/commands; do
	mkdir -p assert-dir/$dir
	cp -R test-dir/$dir/. assert-dir/$dir/
done

for file in \
	AgnosConfig/extensions.yaml \
	sandbox/api/generated.config.go; do
	mkdir -p assert-dir/$(dirname $file)
	cp test-dir/$file assert-dir/$file
done
