# The route-pattern example: --pattern compiles a url shape into the paths of a
# route — literal runs, typed captures, a {*rest} tail — and fixes the segment
# count when there is no tail
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

# One literal run and one typed capture: segments 2, Entries.Article an int.
agnos add-route get-article --pattern '/get-article/{article:integer}' --path TestDir -q

# Two captures, one a uuid, between literals.
agnos add-route get-comment --pattern '/users/{user:uuid}/comments/{comment}' --path TestDir -q

# A tail: the rest of the path, any count of segments.
agnos add-route files --pattern '/files/{*file}' --path TestDir -q

agnos show-route get-article --path TestDir
agnos show-route get-comment --path TestDir
agnos show-route files --path TestDir

# What result.yaml records: each declaration, and the Entries each generates.
for route in get_article get_comment files; do
  mkdir -p AssertDir/sandbox/internal/routeslist/$route
  cp TestDir/sandbox/internal/routeslist/$route/route.yaml AssertDir/sandbox/internal/routeslist/$route/route.yaml
  cp TestDir/sandbox/internal/routeslist/$route/entries.go AssertDir/sandbox/internal/routeslist/$route/entries.go
done
