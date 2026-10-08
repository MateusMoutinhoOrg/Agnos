# The add-page example: declare an html page on a project with the front layer
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos front-init --path test-dir -q

agnos add-page about --title "About" --path test-dir
agnos add-page blog/post --path test-dir -q

# What result.yaml records: the whole of what a page is — one html file under
# assets/front, beside the index.html front-init wrote. No route is
# declared: the front route serves every file of that tree.
mkdir -p assert-dir/assets/front
cp -R test-dir/assets/front/. assert-dir/assets/front/
