# The server-init example: add the http server layer to a project that has none
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q

agnos server-init --path TestDir

# What result.yaml records: the paths this example asserts, copied out of
# TestDir. sandbox/internal/commands/start_server proves the implicit cli-init
# ran: a server needs a command that starts it.
mkdir -p AssertDir/sandbox/api
cp TestDir/sandbox/api/server.go AssertDir/sandbox/api/server.go
mkdir -p AssertDir/sandbox/internal/server
cp -R TestDir/sandbox/internal/server/. AssertDir/sandbox/internal/server/
mkdir -p AssertDir/sandbox/internal/routes
cp -R TestDir/sandbox/internal/routes/. AssertDir/sandbox/internal/routes/
mkdir -p AssertDir/sandbox/internal/commands/start_server
cp -R TestDir/sandbox/internal/commands/start_server/. AssertDir/sandbox/internal/commands/start_server/
