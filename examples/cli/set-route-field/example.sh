# The set-route-field example: add the bound that was forgotten to a
# declaration that already exists, instead of removing it and declaring it again
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route create-user --trigger /users --method POST --help "Create a user under a tenant" --category Users --path TestDir -q
agnos add-segment tenant --route create-user --path TestDir -q
agnos add-header authorization --route create-user --path TestDir -q
agnos add-param page --route create-user --type int --default 1 --path TestDir -q
agnos add-body-field age --route create-user --path TestDir -q

# Each add- has a set- beside it. The keys given are written over the ones
# already declared, --clear takes one off, and --rename changes the name the
# field answers to.
agnos set-segment tenant --route create-user --description "the tenant the user belongs to" --path TestDir -q
agnos set-header authorization --route create-user --required --description "the bearer token" --path TestDir -q
agnos set-param page --route create-user --max 50 --description "the page to read" --path TestDir -q
agnos set-body-field age --route create-user --type int --min 0 --max 130 --required --path TestDir -q
agnos set-param page --route create-user --clear default --rename cursor --path TestDir

# What result.yaml records: the declaration the five edits left behind, and the
# api.Route build generated from it. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/routes/create_user
cp TestDir/sandbox/internal/routes/create_user/route.yaml AssertDir/sandbox/internal/routes/create_user/route.yaml
cp TestDir/sandbox/internal/routes/create_user/new.go AssertDir/sandbox/internal/routes/create_user/new.go
