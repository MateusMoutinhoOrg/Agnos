# The remove-page example: drop a page, its route and its html both
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos front-init --path test-dir -q
agnos add-page about --title "About" --path test-dir -q
agnos add-page blog/post --path test-dir -q

agnos remove-page about --path test-dir

# What result.yaml records: the pages left. about.html is gone, and index.html
# and blog/post.html were not touched.
mkdir -p assert-dir/assets/front
cp -R test-dir/assets/front/. assert-dir/assets/front/
