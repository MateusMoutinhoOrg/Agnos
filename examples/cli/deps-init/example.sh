# The deps-init example: add the dependency layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos deps-init --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/sandbox/deps
cp -R TestDir/sandbox/deps/. AssertDir/sandbox/deps/
mkdir -p AssertDir/adapters
cp -R TestDir/adapters/. AssertDir/adapters/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
