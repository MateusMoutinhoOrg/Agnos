# The dep-install example: install one dep of the built-in list
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos deps-init --path TestDir -q

agnos dep-install iodeps --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/sandbox/deps
cp -R TestDir/sandbox/deps/. AssertDir/sandbox/deps/
mkdir -p AssertDir/adapters
cp -R TestDir/adapters/. AssertDir/adapters/
