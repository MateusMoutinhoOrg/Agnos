# The explain-route example: which routes one request reaches, and why every
# other one is skipped, read off the declarations without a server
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir.

agnos start --path test-dir --project-name Test --module Test -q
agnos server-init --path test-dir -q

agnos add-route admin --trigger /admin --trigger-type prefix --path test-dir -q
agnos add-route admin-guard --middleware --trigger /admin --path test-dir -q
agnos add-parameter authorization --route admin-guard --source header --trigger bearer --trigger-type starts-with --trigger-ignore-case --path test-dir -q
agnos add-route get-article --pattern '/articles/{article:integer}' --path test-dir -q
agnos add-parameter page --route get-article --type integer --default 1 --path test-dir -q
agnos add-route not-api --trigger /api --trigger-type prefix --trigger-negate --priority 500 --path test-dir -q

# The guard runs only with a bearer token; a prefix never reaches /administrator.
agnos explain-route GET /admin/users --header 'authorization=Bearer abc' --path test-dir
agnos explain-route GET /administrator --path test-dir

# A path that does not convert is a non-match; a parameter that does not is a 400.
agnos explain-route GET /articles/abc --path test-dir
agnos explain-route GET '/articles/42?page=x' --path test-dir

# A method no route with explicit methods takes: 405, even with the guard running.
agnos explain-route POST /admin --header 'authorization=Bearer abc' --path test-dir

# A HEAD nothing declares runs again as a GET.
agnos explain-route HEAD /api/health --path test-dir

# What result.yaml records: the two declarations carrying the trigger switches,
# written with the canonical type the alias stood for.
for route in admin_guard not_api; do
  mkdir -p assert-dir/sandbox/internal/routes/$route
  cp test-dir/sandbox/internal/routes/$route/route.yaml assert-dir/sandbox/internal/routes/$route/route.yaml
done
