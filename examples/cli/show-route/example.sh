# The show-route example: read a whole route declaration back as a tree —
# its path, its headers, its query parameters and its body schema
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

agnos add-route create-user --trigger /users/ --trigger-type prefix --method POST --summary "Create a user under a tenant" --category Users --path test-dir -q
agnos add-path tenant --route create-user --start 1 --end 1 --path test-dir -q
agnos add-parameter authorization --route create-user --source header --required --path test-dir -q
agnos add-parameter page --route create-user --type number --default 1 --path test-dir -q
agnos add-body-field email --route create-user --required --format email --max 254 --path test-dir -q
agnos add-body-field address.city --route create-user --required --path test-dir -q
agnos add-body-field tags --route create-user --array --unique-items --max-items 10 --path test-dir -q

agnos show-route create-user --path test-dir

# What result.yaml records: the output above — the tree is the whole of what
# this command answers — plus the declaration it was read from. The lib side
# copies the same set.
mkdir -p assert-dir/sandbox/internal/routes/create_user
cp test-dir/sandbox/internal/routes/create_user/route.yaml assert-dir/sandbox/internal/routes/create_user/route.yaml
