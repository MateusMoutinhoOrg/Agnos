# The route-docs example: the pages docs/Routes generates for whoever calls the
# server — plain words, and curl requests that run as they are — and the
# OpenAPI document beside them, openapi.json, that Postman and Swagger import
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

# A route with a value in the address, one in the query string and a json body.
# Its page carries the smallest request and the one sending every value.
agnos add-route create-user --pattern '/users/{tenant}' --method POST --summary "Create a user under a tenant" --category Users --path test-dir -q
agnos set-path tenant --route create-user --description "the tenant the user belongs to" --path test-dir -q
agnos add-parameter page --route create-user --type integer --default 1 --description "the page to read" --path test-dir -q
agnos set-body create-user --type json --required --path test-dir -q
agnos add-body-field email --route create-user --format email --required --path test-dir -q
agnos add-body-field age --route create-user --type integer --min 0 --max 130 --path test-dir -q
agnos add-body-field address.city --route create-user --path test-dir -q
agnos add-body-field role --route create-user --enum admin --enum member --path test-dir -q

# A typed value in the address: a request whose part is no whole number is not
# this route at all.
agnos add-route get-article --pattern '/articles/{article:integer}' --summary "Read one article" --category Articles --path test-dir -q

# A guard in front of /admin: the header it requires lands on the admin page,
# and in its request, because the guard always runs there.
agnos add-route admin --trigger /admin --trigger-type prefix --summary "The admin panel" --category Admin --path test-dir -q
agnos add-route admin-guard --middleware --trigger /admin --before admin --summary "Checks the bearer token" --path test-dir -q
agnos add-parameter authorization --route admin-guard --source header --required --description "the bearer token" --path test-dir -q

# A path that has to end with .png, one it must not start with, a cookie, and a
# query value the route only runs for.
agnos add-route images --trigger /img --trigger-type prefix --method PUT --category Files --path test-dir -q
agnos add-path ext --route images --start 1 --end -1 --trigger .png --trigger-type suffix --path test-dir -q
agnos add-path not-secret --route images --trigger /img/secret --trigger-type prefix --trigger-negate --path test-dir -q
agnos add-parameter session --route images --source cookie --required --path test-dir -q
agnos add-parameter mode --route images --trigger 'fast,slow' --trigger-type one-of --path test-dir -q
agnos set-body images --type text --path test-dir -q

mkdir -p assert-dir/docs/Routes
cp -R test-dir/docs/Routes/. assert-dir/docs/Routes/
