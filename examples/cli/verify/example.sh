# The verify example: check a project against the schema, writing nothing
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos verify --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
# verify writes nothing, so what this asserts is the config it read, untouched.
mkdir -p assert-dir/AgnosConfig
cp -R test-dir/AgnosConfig/. assert-dir/AgnosConfig/
