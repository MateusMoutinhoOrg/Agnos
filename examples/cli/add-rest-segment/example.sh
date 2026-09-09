# The add-rest-segment example: a route whose last segment takes the rest of
# the path, so one declaration serves /static/a, /static/a/b.png and anything
# deeper.
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route static --trigger /static --help "Serve a file under /static" --category Files --path TestDir -q
agnos add-segment rest --route static --array --path TestDir

# What result.yaml records: the declaration this example wrote, the []string
# field build generated from it, and the dispatch, which is where the route
# stops asking for an exact segment count. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/routes/static AssertDir/sandbox/internal/server
cp TestDir/sandbox/internal/routes/static/route.yaml AssertDir/sandbox/internal/routes/static/route.yaml
cp TestDir/sandbox/internal/routes/static/entries.go AssertDir/sandbox/internal/routes/static/entries.go
cp TestDir/sandbox/internal/server/servermain.go AssertDir/sandbox/internal/server/servermain.go
