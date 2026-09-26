# The remove-page example: drop a page, its route and its html both
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos front-init --path TestDir -q
agnos add-page about --title "About" --path TestDir -q
agnos add-page blog/post --path TestDir -q

agnos remove-page about --path TestDir

# What result.yaml records: the pages left. about.html is gone, and index.html
# and blog/post.html were not touched.
mkdir -p AssertDir/assets/frontend
cp -R TestDir/assets/frontend/. AssertDir/assets/frontend/
