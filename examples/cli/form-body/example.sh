# The form-body example: declare the form-schema of a route an html
# <form method="POST"> posts to, and carry it from a json body and back
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q
agnos add-route login --trigger /login --method POST --summary "Sign in from a form" --category Auth --path test-dir -q

# A json body turned form keeps its flat schema: Input.Body keeps its fields.
agnos set-body login --type json --required --path test-dir -q
agnos add-body-field username --route login --required --max 254 --path test-dir -q
agnos set-body login --type form --path test-dir -q

agnos add-body-field password --route login --required --min 8 --path test-dir -q
agnos add-body-field remember --route login --type boolean --path test-dir -q
agnos add-body-field age --route login --type integer --min 0 --path test-dir -q
agnos add-body-field tag --route login --array --max-items 3 --path test-dir -q

# A form is one flat list of key=value pairs: nothing nests.
agnos add-body-field address.city --route login --path test-dir -q || echo "refused: address.city nests"
agnos add-body-field nickname --route login --nullable --path test-dir -q || echo "refused: nickname is nullable"

agnos show-route login --path test-dir

# What result.yaml records: the form-schema declared, the Body struct and the
# ReadBody build generated from it, and the page the route's callers read.
mkdir -p assert-dir/sandbox/internal/routes/login assert-dir/docs/Routes
cp test-dir/sandbox/internal/routes/login/route.yaml assert-dir/sandbox/internal/routes/login/route.yaml
cp test-dir/sandbox/internal/routes/login/input.go assert-dir/sandbox/internal/routes/login/input.go
cp test-dir/docs/Routes/login.md assert-dir/docs/Routes/login.md
