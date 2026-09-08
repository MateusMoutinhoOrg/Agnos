# The add-doc example: create a doc directory under docs/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos add-doc Report --theme reference --description "How a report is written" --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/docs/Report
cp -R TestDir/docs/Report/. AssertDir/docs/Report/
