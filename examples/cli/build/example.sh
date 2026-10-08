# The build example: regenerate every generated file of a project
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos build --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/docs
cp -R test-dir/docs/. assert-dir/docs/
mkdir -p assert-dir/AgnosConfig
cp -R test-dir/AgnosConfig/. assert-dir/AgnosConfig/
