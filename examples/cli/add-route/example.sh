# The add-route example: declare a route and one entry of every place its
# declaration holds something
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route create-user --trigger /users --method POST --help "Create a user under a tenant" --category Users --path TestDir -q
agnos add-segment tenant --route create-user --path TestDir -q
agnos add-header authorization --route create-user --required --path TestDir -q
agnos add-param page --route create-user --type int --default 1 --min 1 --path TestDir -q
agnos set-body create-user --type json --required --path TestDir -q
agnos add-body-field email --route create-user --required --format email --path TestDir

# What result.yaml records: the declaration this example wrote and the struct
# build generated from it. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/routes/create_user
cp TestDir/sandbox/internal/routes/create_user/route.yaml AssertDir/sandbox/internal/routes/create_user/route.yaml
cp TestDir/sandbox/internal/routes/create_user/entries.go AssertDir/sandbox/internal/routes/create_user/entries.go
