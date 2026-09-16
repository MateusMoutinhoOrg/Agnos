# The cli-init example: add the cli layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos cli-init --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/cmd
cp -R TestDir/cmd/. AssertDir/cmd/
mkdir -p AssertDir/sandbox/internal/cli
cp -R TestDir/sandbox/internal/cli/. AssertDir/sandbox/internal/cli/
mkdir -p AssertDir/sandbox/internal/commands
cp -R TestDir/sandbox/internal/commands/. AssertDir/sandbox/internal/commands/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
