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

# The constructor the layer brought with it, and the new.go that calls it:
# sandbox/new.go is one call per directory of sandbox/constructors/, so this is
# the whole of how Sandbox.Cli comes to be filled.
mkdir -p AssertDir/sandbox/constructors
cp -R TestDir/sandbox/constructors/. AssertDir/sandbox/constructors/
cp TestDir/sandbox/new.go AssertDir/sandbox/new.go

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
