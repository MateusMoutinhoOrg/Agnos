# The set-route-field example: add the key that was forgotten to a
# declaration that already exists, instead of removing it and declaring it again
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route create-user --trigger /users/ --trigger-type prefix --method POST --help "Create a user under a tenant" --category Users --path TestDir -q
agnos add-path tenant --route create-user --start 1 --end 1 --path TestDir -q
agnos add-parameter authorization --route create-user --font header --path TestDir -q
agnos add-parameter page --route create-user --type number --default 1 --path TestDir -q
agnos add-body-field age --route create-user --path TestDir -q

# Each add- has a set- beside it. The keys given are written over the ones
# already declared, --clear takes one off, and --rename changes the name the
# entry answers to.
agnos set-path tenant --route create-user --description "the tenant the user belongs to" --path TestDir -q
agnos set-parameter authorization --route create-user --required --description "the bearer token" --path TestDir -q
agnos set-parameter page --route create-user --font query --font header --description "the page to read" --path TestDir -q
agnos set-body-field age --route create-user --type int --min 0 --max 130 --required --path TestDir -q
agnos set-parameter page --route create-user --clear default --rename cursor --path TestDir

# What result.yaml records: the declaration the five edits left behind, and the
# api.Route and Entries build generated from it. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/routeslist/create_user
cp TestDir/sandbox/internal/routeslist/create_user/route.yaml AssertDir/sandbox/internal/routeslist/create_user/route.yaml
cp TestDir/sandbox/internal/routeslist/create_user/new.go AssertDir/sandbox/internal/routeslist/create_user/new.go
cp TestDir/sandbox/internal/routeslist/create_user/entries.go AssertDir/sandbox/internal/routeslist/create_user/entries.go
