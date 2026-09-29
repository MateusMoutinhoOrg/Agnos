# The form-body example: declare the form-schema of a route an html
# <form method="POST"> posts to, and carry it from a json body and back
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q
agnos add-route login --trigger /login --method POST --help "Sign in from a form" --category Auth --path TestDir -q

# A json body turned form keeps its flat schema: Entries.Body keeps its fields.
agnos set-body login --type json --required --path TestDir -q
agnos add-body-field username --route login --required --max 254 --path TestDir -q
agnos set-body login --type form --path TestDir -q

agnos add-body-field password --route login --required --min 8 --path TestDir -q
agnos add-body-field remember --route login --type boolean --path TestDir -q
agnos add-body-field age --route login --type int --min 0 --path TestDir -q
agnos add-body-field tag --route login --array --max-items 3 --path TestDir -q

# A form is one flat list of key=value pairs: nothing nests.
agnos add-body-field address.city --route login --path TestDir -q || echo "refused: address.city nests"
agnos add-body-field nickname --route login --nullable --path TestDir -q || echo "refused: nickname is nullable"

agnos show-route login --path TestDir

# What result.yaml records: the form-schema declared, the Body struct and the
# ReadBody build generated from it, and the page the route's callers read.
mkdir -p AssertDir/sandbox/internal/routeslist/login AssertDir/docs/Routes
cp TestDir/sandbox/internal/routeslist/login/route.yaml AssertDir/sandbox/internal/routeslist/login/route.yaml
cp TestDir/sandbox/internal/routeslist/login/entries.go AssertDir/sandbox/internal/routeslist/login/entries.go
cp TestDir/docs/Routes/login.md AssertDir/docs/Routes/login.md
