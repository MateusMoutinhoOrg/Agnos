# The add-page example: declare an html page on a project with the front layer
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos front-init --path TestDir -q

agnos add-page home --trigger / --title "Home" --path TestDir

# What result.yaml records: the two halves of a page — the route that answers it
# and the html it renders — plus the entries.go the follow-up build generated
# from the declaration.
mkdir -p AssertDir/sandbox/internal/routes/home
cp -R TestDir/sandbox/internal/routes/home/. AssertDir/sandbox/internal/routes/home/
mkdir -p AssertDir/assets/frontend/pages
cp TestDir/assets/frontend/pages/home.html AssertDir/assets/frontend/pages/home.html
