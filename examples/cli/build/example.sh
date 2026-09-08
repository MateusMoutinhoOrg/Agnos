# The build example: regenerate every generated file of a project
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos build --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/docs
cp -R TestDir/docs/. AssertDir/docs/
mkdir -p AssertDir/AgnosConfig
cp -R TestDir/AgnosConfig/. AssertDir/AgnosConfig/
