# The cli-purge example: remove the cli layer and every command in it
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos cli-init --path TestDir -q

agnos cli-purge --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal
cp -R TestDir/sandbox/internal/. AssertDir/sandbox/internal/
mkdir -p AssertDir/docs
cp -R TestDir/docs/. AssertDir/docs/
