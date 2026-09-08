# The verify example: check a project against the schema, writing nothing
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos verify --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
# verify writes nothing, so what this asserts is the config it read, untouched.
mkdir -p AssertDir/AgnosConfig
cp -R TestDir/AgnosConfig/. AssertDir/AgnosConfig/
