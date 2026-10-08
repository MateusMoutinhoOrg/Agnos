# The disable-extension example: stop generating one mechanic, keep its files
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos disable-extension readme --path test-dir

# Nothing was removed: turning a mechanic off only stops the generation, so
# README.md is still there and is the project's from here on.
echo "README.md still here: $(test -f test-dir/README.md && echo yes || echo no)"

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
cp test-dir/README.md assert-dir/README.md
