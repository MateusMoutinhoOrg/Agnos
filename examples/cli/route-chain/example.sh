# The route-chain example: two routes matching one request, and the six files
# that answer what no route does
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

# A middleware: it sits on the lowest rung and its trigger is a prefix, so every
# request reaches it. Its handler writes no status, which is the whole of what
# makes it a middleware rather than an answer.
agnos add-route logger --trigger / --trigger-type prefix --priority 0 --summary "Logs every request" --category Server --path test-dir -q

# The route it runs in front of, on a higher rung.
agnos add-route hello --trigger /hello --priority 5 --summary "Says hello" --category Server --path test-dir -q

# A route reached only by the requests that carry the header it names: a
# parameter trigger is part of what the route matches on, so a request without it is
# not a bad request — it is a request for some other route.
agnos add-route admin --trigger /admin --trigger-type prefix --priority 5 --summary "The admin area" --category Server --path test-dir -q
agnos add-parameter authorization --route admin --source header --trigger Bearer --trigger-type prefix --description "the bearer token" --path test-dir -q

agnos show-route logger --path test-dir
agnos show-route admin --path test-dir

# What result.yaml records: the three declarations, the api.Route build
# generated from the prefix one, the order the collector put them in, and the
# six handlers that answer every failure.
mkdir -p assert-dir/sandbox/internal/routes/logger
cp test-dir/sandbox/internal/routes/logger/route.yaml assert-dir/sandbox/internal/routes/logger/route.yaml
cp test-dir/sandbox/internal/routes/logger/generated.new.go assert-dir/sandbox/internal/routes/logger/generated.new.go
mkdir -p assert-dir/sandbox/internal/routes/hello
cp test-dir/sandbox/internal/routes/hello/route.yaml assert-dir/sandbox/internal/routes/hello/route.yaml
mkdir -p assert-dir/sandbox/internal/routes/admin
cp test-dir/sandbox/internal/routes/admin/route.yaml assert-dir/sandbox/internal/routes/admin/route.yaml
cp test-dir/sandbox/internal/routes/admin/generated.new.go assert-dir/sandbox/internal/routes/admin/generated.new.go

# server/generated.new.go is where the run order shows: the routes are laid down lowest
# priority first, and beside them the switch that reaches the six handlers.
mkdir -p assert-dir/sandbox/internal/server assert-dir/sandbox/internal/server/errors
cp test-dir/sandbox/internal/server/generated.new.go assert-dir/sandbox/internal/server/generated.new.go
cp test-dir/sandbox/internal/server/errors/handle_not_found.go assert-dir/sandbox/internal/server/errors/handle_not_found.go
cp test-dir/sandbox/internal/server/errors/handle_method_not_allowed.go assert-dir/sandbox/internal/server/errors/handle_method_not_allowed.go
