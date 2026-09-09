# The remove-page example: drop a page, its route and its html both
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos front-init --path TestDir -q
agnos add-page home --trigger / --title "Home" --path TestDir -q
agnos add-page about --title "About" --path TestDir -q

agnos remove-page about --path TestDir

# What result.yaml records: the pages left. `about` is gone from both places a
# page lives, and `home` was not touched.
mkdir -p AssertDir/sandbox/internal/routes
cp -R TestDir/sandbox/internal/routes/. AssertDir/sandbox/internal/routes/
mkdir -p AssertDir/assets/frontend/pages
cp -R TestDir/assets/frontend/pages/. AssertDir/assets/frontend/pages/
