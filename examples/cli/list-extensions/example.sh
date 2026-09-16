# The list-extensions example: what agnos generates for this project
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos list-extensions --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
