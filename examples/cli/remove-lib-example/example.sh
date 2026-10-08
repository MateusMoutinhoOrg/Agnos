# The remove-lib-example example: delete an example of examples/lib/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos add-lib-example greet --path test-dir -q

agnos remove-lib-example greet --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/examples
cp -R test-dir/examples/. assert-dir/examples/
mkdir -p assert-dir/docs/LibExamples
cp -R test-dir/docs/LibExamples/. assert-dir/docs/LibExamples/
