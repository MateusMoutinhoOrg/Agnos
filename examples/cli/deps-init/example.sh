# The deps-init example: add the dependency layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos deps-init --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/sandbox/deps
cp -R test-dir/sandbox/deps/. assert-dir/sandbox/deps/
mkdir -p assert-dir/adapters
cp -R test-dir/adapters/. assert-dir/adapters/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
