# The remove-cli-example example: delete an example of examples/cli/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos cli-init --path TestDir -q
agnos add-cli-example greet --path TestDir -q

agnos remove-cli-example greet --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/examples
cp -R TestDir/examples/. AssertDir/examples/
mkdir -p AssertDir/docs/CliExamples
cp -R TestDir/docs/CliExamples/. AssertDir/docs/CliExamples/
