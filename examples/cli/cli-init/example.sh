# The cli-init example: add the cli layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos cli-init --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/cmd
cp -R test-dir/cmd/. assert-dir/cmd/
mkdir -p assert-dir/sandbox/internal/cli
cp -R test-dir/sandbox/internal/cli/. assert-dir/sandbox/internal/cli/
mkdir -p assert-dir/sandbox/internal/commands
cp -R test-dir/sandbox/internal/commands/. assert-dir/sandbox/internal/commands/

# The constructor the layer brought with it, and the generated.new.go that calls it:
# sandbox/generated.new.go is one call per directory of sandbox/constructors/, so this is
# the whole of how Sandbox.Cli comes to be filled.
mkdir -p assert-dir/sandbox/constructors
cp -R test-dir/sandbox/constructors/. assert-dir/sandbox/constructors/
cp test-dir/sandbox/generated.new.go assert-dir/sandbox/generated.new.go

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
