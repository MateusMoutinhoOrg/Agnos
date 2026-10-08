# The remove-doc example: delete a doc directory
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos add-doc Report --theme reference --description "How a report is written" --path test-dir -q

agnos remove-doc Report --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/docs
cp -R test-dir/docs/. assert-dir/docs/
