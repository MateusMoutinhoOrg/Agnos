# The add-binding example: declare a second binding
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos deps-init --path test-dir -q
agnos add-dep sortdeps --path test-dir -q

agnos add-binding lambda --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/adapters
cp -R test-dir/adapters/. assert-dir/adapters/
