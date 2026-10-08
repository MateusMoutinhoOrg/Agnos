# The add-path-range example: a route whose second path reads every segment
# after the mount, so one declaration serves /static/a, /static/a/b.png and
# anything deeper.
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

agnos add-route static --trigger /static --summary "Serve a file under /static" --category Files --path test-dir -q
agnos set-path Static --route static --end 0 --path test-dir -q
agnos add-path rest --route static --start 1 --end -1 --path test-dir

# What result.yaml records: the declaration this example wrote and the
# api.Route and Input build generated from it: the first path fixes segment
# 0, the second reads segment 1 to the last. The lib side copies the same set.
mkdir -p assert-dir/sandbox/internal/routes/static
cp test-dir/sandbox/internal/routes/static/route.yaml assert-dir/sandbox/internal/routes/static/route.yaml
cp test-dir/sandbox/internal/routes/static/new.go assert-dir/sandbox/internal/routes/static/new.go
cp test-dir/sandbox/internal/routes/static/input.go assert-dir/sandbox/internal/routes/static/input.go
