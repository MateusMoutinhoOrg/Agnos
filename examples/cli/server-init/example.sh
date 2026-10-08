# The server-init example: add the http server layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q

agnos server-init --path test-dir

# What result.yaml records: the paths this example asserts, copied out of
# test-dir. sandbox/internal/commands/server/start_server proves the implicit cli-init
# ran: a server needs a command that starts it.
mkdir -p assert-dir/sandbox/api
cp test-dir/sandbox/api/server.go assert-dir/sandbox/api/server.go
cp test-dir/sandbox/api/route.go assert-dir/sandbox/api/route.go
mkdir -p assert-dir/sandbox/internal/server
cp -R test-dir/sandbox/internal/server/. assert-dir/sandbox/internal/server/
mkdir -p assert-dir/sandbox/internal/routes
cp -R test-dir/sandbox/internal/routes/. assert-dir/sandbox/internal/routes/
mkdir -p assert-dir/sandbox/internal/commands/server/start_server
cp -R test-dir/sandbox/internal/commands/server/start_server/. assert-dir/sandbox/internal/commands/server/start_server/

# The declaration the pair wrote: this is the whole of what tells the build the
# mechanic is on or off from here.
mkdir -p assert-dir/AgnosConfig
cp test-dir/AgnosConfig/extensions.yaml assert-dir/AgnosConfig/extensions.yaml
