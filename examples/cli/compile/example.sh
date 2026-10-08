# The compile example: cross-compile a project's cmd/main into release/
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos cli-init --path test-dir -q

agnos compile --target linux86 --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
# The binary in release/ is machine-specific and no golden, so what this
# asserts is the source it was built from, untouched.
mkdir -p assert-dir/cmd
cp -R test-dir/cmd/. assert-dir/cmd/
