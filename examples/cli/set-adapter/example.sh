# The set-adapter example: point one available at another adapter
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos deps-init --path TestDir -q
agnos add-dep sortdeps --path TestDir -q

agnos add-adapter reflectsort --path TestDir -q
agnos add-available lambda --path TestDir -q

agnos set-adapter sortdeps reflectsort --path TestDir --available lambda

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/adapters
cp -R TestDir/adapters/. AssertDir/adapters/
