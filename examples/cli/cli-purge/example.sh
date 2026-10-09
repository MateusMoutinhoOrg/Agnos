# The cli-purge example: remove the cli layer and every command in it
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos cli-init --path test-dir -q

agnos cli-purge --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/sandbox/internal
cp -R test-dir/sandbox/internal/. assert-dir/sandbox/internal/

# The purge takes sandbox/constructors/cli with the layer, so generated.new.go comes out
# of the following build calling nothing at all.
cp test-dir/sandbox/generated.new.go assert-dir/sandbox/generated.new.go
mkdir -p assert-dir/docs
cp -R test-dir/docs/. assert-dir/docs/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
