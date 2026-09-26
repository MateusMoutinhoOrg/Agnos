# The route-middleware example: a guard in front of the routes it protects, an
# access log behind every answer, and the chain they make
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

# A route on the default rung, 100. A prefix trigger holds on a segment
# boundary, so /admin/users is this route and /administrator is not.
agnos add-route admin --trigger /admin --trigger-type starts-with --path TestDir -q

# A middleware: ANY method, rung 10 by default — here one rung below admin —
# and a stub handler that answers nothing, so the chain goes on.
agnos add-route admin-guard --middleware --trigger /admin --before admin --path TestDir -q

# The after phase: runs once the request has been answered, whatever answered it.
agnos add-route access-log --middleware --phase after --path TestDir -q

# The same chain, laid down again ten rungs apart.
agnos list-routes --path TestDir
agnos rebalance-routes --path TestDir -q
agnos list-routes --path TestDir

# A rename moves the hand-written files and rewrites their package clause.
agnos rename-route admin-guard admin-auth --path TestDir -q
agnos show-route admin-auth --path TestDir

# What result.yaml records: the declarations, the two middleware stubs, and the
# generated server/new.go — the run order, and the eight handlers beside it.
for route in admin admin_auth access_log; do
  mkdir -p AssertDir/sandbox/internal/routeslist/$route
  cp TestDir/sandbox/internal/routeslist/$route/route.yaml AssertDir/sandbox/internal/routeslist/$route/route.yaml
done
cp TestDir/sandbox/internal/routeslist/admin_auth/InternalPureHandler.go AssertDir/sandbox/internal/routeslist/admin_auth/InternalPureHandler.go
cp TestDir/sandbox/internal/routeslist/access_log/InternalPureHandler.go AssertDir/sandbox/internal/routeslist/access_log/InternalPureHandler.go
mkdir -p AssertDir/sandbox/internal/server/server
cp TestDir/sandbox/internal/server/server/new.go AssertDir/sandbox/internal/server/server/new.go
