# The disable-extension example: stop generating one mechanic, keep its files
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos disable-extension readme --path TestDir

# Nothing was removed: turning a mechanic off only stops the generation, so
# README.md is still there and is the project's from here on.
echo "README.md still here: $(test -f TestDir/README.md && echo yes || echo no)"

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/AgnosConfig
cp TestDir/AgnosConfig/extensions.yaml AssertDir/AgnosConfig/extensions.yaml
cp TestDir/README.md AssertDir/README.md
