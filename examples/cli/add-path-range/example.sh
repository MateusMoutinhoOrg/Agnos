# The add-path-range example: a route whose second path reads every segment
# after the mount, so one declaration serves /static/a, /static/a/b.png and
# anything deeper.
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route static --trigger /static --help "Serve a file under /static" --category Files --path TestDir -q
agnos set-path Route --route static --end 0 --path TestDir -q
agnos add-path rest --route static --start 1 --end -1 --path TestDir

# What result.yaml records: the declaration this example wrote and the
# api.Route and Entries build generated from it: the first path fixes segment
# 0, the second reads segment 1 to the last. The lib side copies the same set.
mkdir -p AssertDir/sandbox/internal/routeslist/static
cp TestDir/sandbox/internal/routeslist/static/route.yaml AssertDir/sandbox/internal/routeslist/static/route.yaml
cp TestDir/sandbox/internal/routeslist/static/new.go AssertDir/sandbox/internal/routeslist/static/new.go
cp TestDir/sandbox/internal/routeslist/static/entries.go AssertDir/sandbox/internal/routeslist/static/entries.go
