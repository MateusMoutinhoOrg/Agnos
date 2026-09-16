# The deps-purge example: remove the dependency layer and every installed dep
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos deps-init --path TestDir -q

agnos deps-purge --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/sandbox
cp -R TestDir/sandbox/. AssertDir/sandbox/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
