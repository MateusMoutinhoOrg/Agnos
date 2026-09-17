# The show-route example: read a whole route declaration back as a tree —
# its path, its headers, its query parameters and its body schema
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route create-user --trigger /users --method POST --help "Create a user under a tenant" --category Users --path TestDir -q
agnos add-segment tenant --route create-user --path TestDir -q
agnos add-header authorization --route create-user --required --path TestDir -q
agnos add-param page --route create-user --type int --default 1 --min 1 --path TestDir -q
agnos add-body-field email --route create-user --required --format email --max 254 --path TestDir -q
agnos add-body-field address.city --route create-user --required --path TestDir -q
agnos add-body-field tags --route create-user --array --unique-items --max-items 10 --path TestDir -q

agnos show-route create-user --path TestDir

# What result.yaml records: the output above — the tree is the whole of what
# this command answers — plus the declaration it was read from. The lib side
# copies the same set.
mkdir -p AssertDir/sandbox/internal/routes/create_user
cp TestDir/sandbox/internal/routes/create_user/route.yaml AssertDir/sandbox/internal/routes/create_user/route.yaml
