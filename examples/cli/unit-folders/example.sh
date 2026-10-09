# The unit-folders example: routes and commands grouped in folders. A
# directory is a route by holding a route.yaml, and a command by holding a
# command.yaml, at any depth; every other directory is a folder grouping them.
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

# --dir is the folder the unit lands in: routes/admin/login.
agnos add-route login --dir admin --method POST --summary "Log in" --path test-dir -q
agnos add-route users --dir admin/v1 --summary "List the users" --path test-dir -q
agnos add-command deploy --dir ops/cloud --summary "Deploy the project" --path test-dir -q

# The editors find a unit by its name, whatever folder holds it.
agnos add-parameter token --route login --source header --path test-dir -q
agnos add-flag region -c deploy --path test-dir -q

# rename with --dir moves it; the folder it leaves empty goes with it.
agnos rename-route users users --dir api --path test-dir -q

# remove drops the folder it leaves empty too: admin/ goes with login.
agnos remove-route login --path test-dir

# What result.yaml records: the declarations where they now sit, and the
# generated dispatch importing each from its own folder.
mkdir -p assert-dir/sandbox/internal/routes/api/users assert-dir/sandbox/internal/commands/ops/cloud/deploy
mkdir -p assert-dir/sandbox/internal/server assert-dir/sandbox/internal/cli
cp test-dir/sandbox/internal/routes/api/users/route.yaml assert-dir/sandbox/internal/routes/api/users/route.yaml
cp test-dir/sandbox/internal/commands/ops/cloud/deploy/command.yaml assert-dir/sandbox/internal/commands/ops/cloud/deploy/command.yaml
cp test-dir/sandbox/internal/server/generated.new.go assert-dir/sandbox/internal/server/generated.new.go
cp test-dir/sandbox/internal/cli/generated.new.go assert-dir/sandbox/internal/cli/generated.new.go
ls test-dir/sandbox/internal/routes
