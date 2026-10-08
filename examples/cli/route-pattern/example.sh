# The route-pattern example: --pattern compiles a url shape into the paths of a
# route — literal runs, typed captures, a {*rest} tail — and fixes the segment
# count when there is no tail
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

# One literal run and one typed capture: segments 2, Input.Article an int.
agnos add-route get-article --pattern '/get-article/{article:integer}' --path test-dir -q

# Two captures, one a uuid, between literals.
agnos add-route get-comment --pattern '/users/{user:uuid}/comments/{comment}' --path test-dir -q

# A tail: the rest of the path, any count of segments.
agnos add-route files --pattern '/files/{*file}' --path test-dir -q

agnos show-route get-article --path test-dir
agnos show-route get-comment --path test-dir
agnos show-route files --path test-dir

# What result.yaml records: each declaration, and the Input each generates.
for route in get_article get_comment files; do
  mkdir -p assert-dir/sandbox/internal/routes/$route
  cp test-dir/sandbox/internal/routes/$route/route.yaml assert-dir/sandbox/internal/routes/$route/route.yaml
  cp test-dir/sandbox/internal/routes/$route/input.go assert-dir/sandbox/internal/routes/$route/input.go
done
