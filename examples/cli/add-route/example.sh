# The add-route example: declare a route and one entry of every place its
# declaration holds something
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route create-user --trigger /users/ --trigger-type prefix --method POST --help "Create a user under a tenant" --category Users --path TestDir -q
agnos add-path tenant --route create-user --start 1 --end 1 --path TestDir -q
agnos add-parameter authorization --route create-user --font header --required --path TestDir -q
agnos add-parameter page --route create-user --type number --default 1 --path TestDir -q
agnos set-body create-user --type json --required --path TestDir -q
agnos add-body-field email --route create-user --required --format email --path TestDir

# What result.yaml records: the declaration this example wrote and the
# api.Route build generated from it. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/routeslist/create_user
cp TestDir/sandbox/internal/routeslist/create_user/route.yaml AssertDir/sandbox/internal/routeslist/create_user/route.yaml
cp TestDir/sandbox/internal/routeslist/create_user/new.go AssertDir/sandbox/internal/routeslist/create_user/new.go
cp TestDir/sandbox/internal/routeslist/create_user/entries.go AssertDir/sandbox/internal/routeslist/create_user/entries.go
