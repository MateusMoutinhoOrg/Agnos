# The add-flag example: declare one flag on a command
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos cli-init --path test-dir -q
agnos add-command greet --summary "Greet someone" --category "Core" --path test-dir -q

agnos add-flag name --command greet --type string --default world --description "who to greet" --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. The lib side copies the same set.
mkdir -p assert-dir/sandbox/internal/commands/core/greet
cp -R test-dir/sandbox/internal/commands/core/greet/. assert-dir/sandbox/internal/commands/core/greet/
