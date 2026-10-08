# The add-route example: declare a route and one entry of every place its
# declaration holds something
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

agnos add-route create-user --trigger /users/ --trigger-type prefix --method POST --summary "Create a user under a tenant" --category Users --path test-dir -q
agnos add-path tenant --route create-user --start 1 --end 1 --path test-dir -q
agnos add-parameter authorization --route create-user --source header --required --path test-dir -q
agnos add-parameter page --route create-user --type number --default 1 --path test-dir -q
agnos set-body create-user --type json --required --path test-dir -q
agnos add-body-field email --route create-user --required --format email --path test-dir

# What result.yaml records: the declaration this example wrote and the
# api.Route build generated from it. The lib side copies the same set.
mkdir -p assert-dir/sandbox/internal/routes/create_user
cp test-dir/sandbox/internal/routes/create_user/route.yaml assert-dir/sandbox/internal/routes/create_user/route.yaml
cp test-dir/sandbox/internal/routes/create_user/new.go assert-dir/sandbox/internal/routes/create_user/new.go
cp test-dir/sandbox/internal/routes/create_user/input.go assert-dir/sandbox/internal/routes/create_user/input.go
