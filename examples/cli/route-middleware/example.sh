# The route-middleware example: a guard in front of the routes it protects, an
# access log in front of every request, and the chain they make
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

# A route on the default rung, 100. A prefix trigger holds on a segment
# boundary, so /admin/users is this route and /administrator is not.
agnos add-route admin --trigger /admin --trigger-type starts-with --path test-dir -q

# A middleware: ANY method, rung 10 by default — here one rung below admin —
# and a stub handler that answers nothing, so the chain goes on.
agnos add-route admin-guard --middleware --trigger /admin --before admin --path test-dir -q

# A middleware on rung 0, in front of everything: it sees every request and
# answers none of them.
agnos add-route access-log --middleware --priority 0 --path test-dir -q

# The same chain, laid down again ten rungs apart.
agnos list-routes --path test-dir
agnos rebalance-routes --path test-dir -q
agnos list-routes --path test-dir

# A rename moves the hand-written files and rewrites their package clause.
agnos rename-route admin-guard admin-auth --path test-dir -q
agnos show-route admin-auth --path test-dir

# What result.yaml records: the declarations, the two middleware stubs, and the
# generated server/generated.new.go — the run order, and the eight handlers beside it.
for route in admin admin_auth access_log; do
  mkdir -p assert-dir/sandbox/internal/routes/$route
  cp test-dir/sandbox/internal/routes/$route/route.yaml assert-dir/sandbox/internal/routes/$route/route.yaml
done
cp test-dir/sandbox/internal/routes/admin_auth/handler.go assert-dir/sandbox/internal/routes/admin_auth/handler.go
cp test-dir/sandbox/internal/routes/access_log/handler.go assert-dir/sandbox/internal/routes/access_log/handler.go
mkdir -p assert-dir/sandbox/internal/server
cp test-dir/sandbox/internal/server/generated.new.go assert-dir/sandbox/internal/server/generated.new.go
