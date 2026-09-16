# The enable-extension example: turn one generation mechanic back on
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

# `readme` is on in a fresh project, so turn it off first: what the mechanic
# wrote stays on disk, and the next build leaves it alone.
agnos disable-extension readme --path TestDir -q
rm TestDir/README.md

agnos enable-extension readme --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
cp TestDir/README.md AssertDir/README.md
