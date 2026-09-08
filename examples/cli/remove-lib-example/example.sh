# The remove-lib-example example: delete an example of examples/lib/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos add-lib-example greet --path TestDir -q

agnos remove-lib-example greet --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/examples
cp -R TestDir/examples/. AssertDir/examples/
mkdir -p AssertDir/docs/LibExamples
cp -R TestDir/docs/LibExamples/. AssertDir/docs/LibExamples/
