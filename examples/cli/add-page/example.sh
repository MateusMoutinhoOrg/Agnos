# The add-page example: declare an html page on a project with the front layer
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos front-init --path TestDir -q

agnos add-page about --title "About" --path TestDir
agnos add-page blog/post --path TestDir -q

# What result.yaml records: the whole of what a page is — one html file under
# assets/frontend, beside the index.html front-init wrote. No route is
# declared: the frontend route serves every file of that tree.
mkdir -p AssertDir/assets/frontend
cp -R TestDir/assets/frontend/. AssertDir/assets/frontend/
