# The unit-folders example: routes and commands grouped in folders. A
# directory is a route by holding a route.yaml, and a command by holding a
# command.yaml, at any depth; every other directory is a folder grouping them.
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

# --dir is the folder the unit lands in: routeslist/admin/login.
agnos add-route login --dir admin --method POST --help "Log in" --path TestDir -q
agnos add-route users --dir admin/v1 --help "List the users" --path TestDir -q
agnos add-command deploy --dir ops/cloud --help "Deploy the project" --path TestDir -q

# The editors find a unit by its name, whatever folder holds it.
agnos add-parameter token --route login --font header --path TestDir -q
agnos add-flag region -c deploy --path TestDir -q

# rename with --dir moves it; the folder it leaves empty goes with it.
agnos rename-route users users --dir api --path TestDir -q

# remove drops the folder it leaves empty too: admin/ goes with login.
agnos remove-route login --path TestDir

# What result.yaml records: the declarations where they now sit, and the
# generated dispatch importing each from its own folder.
mkdir -p AssertDir/sandbox/internal/routeslist/api/users AssertDir/sandbox/internal/commands/ops/cloud/deploy
mkdir -p AssertDir/sandbox/internal/generated/server/server AssertDir/sandbox/internal/generated/cli/cli
cp TestDir/sandbox/internal/routeslist/api/users/route.yaml AssertDir/sandbox/internal/routeslist/api/users/route.yaml
cp TestDir/sandbox/internal/commands/ops/cloud/deploy/command.yaml AssertDir/sandbox/internal/commands/ops/cloud/deploy/command.yaml
cp TestDir/sandbox/internal/generated/server/server/new.go AssertDir/sandbox/internal/generated/server/server/new.go
cp TestDir/sandbox/internal/generated/cli/cli/new.go AssertDir/sandbox/internal/generated/cli/cli/new.go
ls TestDir/sandbox/internal/routeslist
