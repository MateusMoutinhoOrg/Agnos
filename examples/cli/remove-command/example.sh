# The remove-command example: delete a command and unwire its dispatch
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos cli-init --path TestDir -q
agnos add-command greet --help "Greet someone" --category "Core" --path TestDir -q

agnos remove-command greet --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/commands
cp -R TestDir/sandbox/internal/commands/. AssertDir/sandbox/internal/commands/
mkdir -p AssertDir/docs/Commands
cp -R TestDir/docs/Commands/. AssertDir/docs/Commands/
