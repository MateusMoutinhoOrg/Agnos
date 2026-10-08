# The set-adapter example: point one binding at another adapter
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos deps-init --path test-dir -q
agnos add-dep sortdeps --path test-dir -q

agnos add-adapter reflectsort --path test-dir -q
agnos add-binding lambda --path test-dir -q

agnos set-adapter sortdeps reflectsort --path test-dir --binding lambda

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/adapters
cp -R test-dir/adapters/. assert-dir/adapters/
