# The route-chain example: two routes matching one request, and the six files
# that answer what no route does
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

# A middleware: it sits on the lowest rung and its trigger is a prefix, so every
# request reaches it. Its handler writes no status, which is the whole of what
# makes it a middleware rather than an answer.
agnos add-route logger --trigger / --starts-with --priority 0 --help "Logs every request" --category Server --path TestDir -q

# The route it runs in front of, on a higher rung.
agnos add-route hello --trigger /hello --priority 5 --help "Says hello" --category Server --path TestDir -q

# A route reached only by the requests that carry the header it names: a value
# condition is part of what the route matches on, so a request without it is
# not a bad request — it is a request for some other route.
agnos add-route admin --trigger /admin --priority 5 --help "The admin area" --category Server --path TestDir -q
agnos add-header authorization --route admin --starts-with Bearer --description "the bearer token" --path TestDir -q

agnos show-route logger --path TestDir
agnos show-route admin --path TestDir

# What result.yaml records: the three declarations, the api.Route build
# generated from the prefix one, the order the collector put them in, and the
# six handlers that answer every failure.
mkdir -p AssertDir/sandbox/internal/routes/logger
cp TestDir/sandbox/internal/routes/logger/route.yaml AssertDir/sandbox/internal/routes/logger/route.yaml
cp TestDir/sandbox/internal/routes/logger/new.go AssertDir/sandbox/internal/routes/logger/new.go
mkdir -p AssertDir/sandbox/internal/routes/hello
cp TestDir/sandbox/internal/routes/hello/route.yaml AssertDir/sandbox/internal/routes/hello/route.yaml
mkdir -p AssertDir/sandbox/internal/routes/admin
cp TestDir/sandbox/internal/routes/admin/route.yaml AssertDir/sandbox/internal/routes/admin/route.yaml
cp TestDir/sandbox/internal/routes/admin/new.go AssertDir/sandbox/internal/routes/admin/new.go

# server/new.go is where the run order shows: the routes are laid down lowest
# priority first, and beside them the switch that reaches the six handlers.
mkdir -p AssertDir/sandbox/internal/server
cp TestDir/sandbox/internal/server/new.go AssertDir/sandbox/internal/server/new.go
cp TestDir/sandbox/internal/server/handle_not_found.go AssertDir/sandbox/internal/server/handle_not_found.go
cp TestDir/sandbox/internal/server/handle_method_not_allowed.go AssertDir/sandbox/internal/server/handle_method_not_allowed.go
