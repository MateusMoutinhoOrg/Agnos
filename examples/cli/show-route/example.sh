# The show-route example: read a whole route declaration back as a tree —
# its path, its headers, its query parameters and its body schema
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route create-user --trigger /users/ --trigger-type prefix --method POST --help "Create a user under a tenant" --category Users --path TestDir -q
agnos add-path tenant --route create-user --start 1 --end 1 --path TestDir -q
agnos add-parameter authorization --route create-user --font header --required --path TestDir -q
agnos add-parameter page --route create-user --type number --default 1 --path TestDir -q
agnos add-body-field email --route create-user --required --format email --max 254 --path TestDir -q
agnos add-body-field address.city --route create-user --required --path TestDir -q
agnos add-body-field tags --route create-user --array --unique-items --max-items 10 --path TestDir -q

agnos show-route create-user --path TestDir

# What result.yaml records: the output above — the tree is the whole of what
# this command answers — plus the declaration it was read from. The lib side
# copies the same set.
mkdir -p AssertDir/sandbox/internal/routeslist/create_user
cp TestDir/sandbox/internal/routeslist/create_user/route.yaml AssertDir/sandbox/internal/routeslist/create_user/route.yaml
