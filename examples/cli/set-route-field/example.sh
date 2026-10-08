# The set-route-field example: add the key that was forgotten to a
# declaration that already exists, instead of removing it and declaring it again
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

agnos add-route create-user --trigger /users/ --trigger-type prefix --method POST --summary "Create a user under a tenant" --category Users --path test-dir -q
agnos add-path tenant --route create-user --start 1 --end 1 --path test-dir -q
agnos add-parameter authorization --route create-user --source header --path test-dir -q
agnos add-parameter page --route create-user --type number --default 1 --path test-dir -q
agnos add-body-field age --route create-user --path test-dir -q

# Each add- has a set- beside it. The keys given are written over the ones
# already declared, --clear takes one off, and --rename changes the name the
# entry answers to.
agnos set-path tenant --route create-user --description "the tenant the user belongs to" --path test-dir -q
agnos set-parameter authorization --route create-user --required --description "the bearer token" --path test-dir -q
agnos set-parameter page --route create-user --source query --source header --description "the page to read" --path test-dir -q
agnos set-body-field age --route create-user --type integer --min 0 --max 130 --required --path test-dir -q
agnos set-parameter page --route create-user --clear default --rename cursor --path test-dir

# What result.yaml records: the declaration the five edits left behind, and the
# api.Route and Input build generated from it. The lib side copies the same set.
mkdir -p assert-dir/sandbox/internal/routes/create_user
cp test-dir/sandbox/internal/routes/create_user/route.yaml assert-dir/sandbox/internal/routes/create_user/route.yaml
cp test-dir/sandbox/internal/routes/create_user/new.go assert-dir/sandbox/internal/routes/create_user/new.go
cp test-dir/sandbox/internal/routes/create_user/input.go assert-dir/sandbox/internal/routes/create_user/input.go
