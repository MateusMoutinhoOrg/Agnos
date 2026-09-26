# The explain-route example: which routes one request reaches, and why every
# other one is skipped, read off the declarations without a server
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos exec-test`.
# The example writes only inside TestDir.

agnos start --path TestDir --project-name Test --module Test -q
agnos server-init --path TestDir -q

agnos add-route admin --trigger /admin --trigger-type prefix --path TestDir -q
agnos add-route admin-guard --middleware --trigger /admin --path TestDir -q
agnos add-parameter authorization --route admin-guard --font header --trigger bearer --trigger-type starts-with --trigger-ignore-case --path TestDir -q
agnos add-route get-article --pattern '/articles/{article:integer}' --path TestDir -q
agnos add-parameter page --route get-article --type integer --default 1 --path TestDir -q
agnos add-route not-api --trigger /api --trigger-type prefix --trigger-negate --priority 500 --path TestDir -q

# The guard runs only with a bearer token; a prefix never reaches /administrator.
agnos explain-route GET /admin/users --header 'authorization=Bearer abc' --path TestDir
agnos explain-route GET /administrator --path TestDir

# A path that does not convert is a non-match; a parameter that does not is a 400.
agnos explain-route GET /articles/abc --path TestDir
agnos explain-route GET '/articles/42?page=x' --path TestDir

# A method no route with explicit methods takes: 405, even with the guard running.
agnos explain-route POST /admin --header 'authorization=Bearer abc' --path TestDir

# A HEAD nothing declares runs again as a GET.
agnos explain-route HEAD /api/health --path TestDir

# What result.yaml records: the two declarations carrying the trigger switches,
# written with the canonical type the alias stood for.
for route in admin_guard not_api; do
  mkdir -p AssertDir/sandbox/internal/routeslist/$route
  cp TestDir/sandbox/internal/routeslist/$route/route.yaml AssertDir/sandbox/internal/routeslist/$route/route.yaml
done
