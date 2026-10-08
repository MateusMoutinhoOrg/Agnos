# The enable-extension example: turn one generation mechanic back on
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

# `readme` is on in a fresh project, so turn it off first: what the mechanic
# wrote stays on disk, and the next build leaves it alone.
agnos disable-extension readme --path test-dir -q
rm test-dir/README.md

agnos enable-extension readme --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
cp test-dir/README.md assert-dir/README.md
