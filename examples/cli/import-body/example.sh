# The import-body example: declare a route's whole body json-schema from one
# example payload, instead of one add-body-field per key
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q
agnos add-route create-user --trigger /users --method POST --summary "Create a user" --category Users --path test-dir -q

cat > test-dir/payload.json <<'JSON'
{
  "email": "someone@example.com",
  "age": 30,
  "score": 1.5,
  "active": true,
  "nickname": null,
  "created_at": "2024-01-02T03:04:05Z",
  "id": "3f2504e0-4f89-11d3-9a0c-0305e82c3301",
  "site": "https://example.com",
  "tags": ["blue"],
  "address": { "city": "Sao Paulo" }
}
JSON

agnos import-body create-user --file test-dir/payload.json --required --infer-format --path test-dir

# What result.yaml records: the schema the payload was read as, and the Body
# struct build generated from it. The lib side copies the same set.
mkdir -p assert-dir/sandbox/internal/routes/create_user
cp test-dir/sandbox/internal/routes/create_user/route.yaml assert-dir/sandbox/internal/routes/create_user/route.yaml
cp test-dir/sandbox/internal/routes/create_user/generated.input.go assert-dir/sandbox/internal/routes/create_user/generated.input.go
