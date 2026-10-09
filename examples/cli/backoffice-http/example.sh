# The backoffice-http example: run the backoffice a project gets from
# backoffice-init and talk to it over http — sign-in limits, the session
# cookie, the openapi document, a backup restored over a failure — so its
# behavior, not only the files it writes, is checked.
#
# `agnos` here is this repository's own cli, put on the PATH by `agnos run-examples`.
# The example writes only inside test-dir. It needs curl, and a free port
# among 127.0.0.1:39300 to 39399.

agnos start --path test-dir --project-name Test --module Test -q
agnos backoffice-init --path test-dir -q

# A dotfile copied into the front by mistake, compiled in with the rest.
echo "SECRET=1" > test-dir/assets/front/.env

cd test-dir
go build -o app.bin ./cmd/main
password=$(./app.bin add-backoffice-user --username admin --email admin@example.com --role root 2>&1 | sed -n 's/^password: //p')

TEST_BACKOFFICE_SECRET=0123456789abcdef0123456789abcdef ./app.bin start-server --insecure-http --addr 127.0.0.1:39300:39399 > server.log 2>&1 &
server=$!
trap 'kill $server 2>/dev/null' EXIT
tries=0
until grep -q 'server listening on' server.log || [ $tries -ge 100 ]; do sleep 0.1; tries=$((tries + 1)); done
base=http://$(sed -n 's/^server listening on \([^ ]*\).*/\1/p' server.log)
origin="Origin: $base"
out=http.txt

status() { curl -s -o /dev/null -w '%{http_code}' "$@"; }
count() { sort | uniq -c | awk '{print $2 " x" $1}' | tr '\n' ' '; }
api() { curl -s -H "Authorization: Bearer $token" "$@"; }
idle() {
	tries=0
	until api "$base/api/admin/root/list-backups" | grep -q '"busy":false' || [ $tries -ge 300 ]; do sleep 0.1; tries=$((tries + 1)); done
}

echo "GET /admin/login: $(status "$base/admin/login")" > $out

# Sign in: the cookie is sent to /admin alone.
curl -s -D headers.txt -c cookies.txt -o /dev/null -H "$origin" --data-urlencode "username=admin" --data-urlencode "password=$password" "$base/admin/login"
echo "session cookie: $(grep -i '^set-cookie' headers.txt | grep -o 'Path=[^;]*')" >> $out
echo "GET /admin signed in: $(curl -s -o /dev/null -w '%{http_code} %{redirect_url}' -b cookies.txt "$base/admin" | sed "s|$base||")" >> $out

# An API token: the password is asked again.
echo "token without password: $(status -b cookies.txt -H "$origin" -d 'name=t0&expiration=30' "$base/admin/add-backoffice-api-token")" >> $out
token=$(curl -s -b cookies.txt -H "$origin" --data-urlencode "password=$password" -d 'name=ci&expiration=30' "$base/admin/add-backoffice-api-token" | grep -o 'bo_[0-9a-f]*' | head -1)
echo "GET /api/admin/me: $(status -H "Authorization: Bearer $token" "$base/api/admin/me")" >> $out
spare=$(curl -s -b cookies.txt -H "$origin" --data-urlencode "password=$password" -d 'name=spare&expiration=7' "$base/admin/add-backoffice-api-token" | grep -o 'bo_[0-9a-f]*' | head -1)

# Thirty wrong passwords sent at once: each is counted before it is checked,
# so ten are checked (401) and the rest refused unchecked (429).
pids=""
for i in $(seq 30); do
	curl -s -o /dev/null -w '%{http_code}\n' -H "$origin" -d "username=admin&password=wrong-$i" "$base/admin/login" >> burst.txt &
	pids="$pids $!"
done
wait $pids
echo "30 wrong sign-ins at once: $(count < burst.txt)" >> $out
echo "a 64 KiB sign-in: $(head -c 65536 /dev/zero | tr '\0' a | status -H "$origin" --data-urlencode 'username@-' -d password=x "$base/admin/login")" >> $out

echo "a json body without Content-Type: $(api -o /dev/null -w '%{http_code}' -X POST -d '{"name":"x"}' "$base/api/admin/root/create-backup")" >> $out
echo "GET /.env: $(status "$base/.env")" >> $out
echo "/admin in /openapi.json: $(curl -s "$base/openapi.json" | grep -c '/admin')" >> $out
echo "nosniff on /: $(curl -s -D - -o /dev/null "$base/" | grep -ci '^x-content-type-options: nosniff')" >> $out

# A database of the application, and a snapshot of it.
mkdir -p data/app/items
for i in 1 2 3; do echo "item $i" > data/app/items/$i; done
api -o /dev/null -X POST -H 'Content-Type: application/json' -d '{"name":"first"}' "$base/api/admin/root/create-backup"
idle

# A token revoked after the snapshot stays revoked once it is restored: the
# backoffice's own users and tokens are left out unless asked for.
spare_id=$(curl -s -b cookies.txt "$base/admin/list-backoffice-api-tokens" | grep -o 'revoke-backoffice-api-token/[0-9]*' | sort -t/ -k2 -n | tail -1)
curl -s -o /dev/null -b cookies.txt -H "$origin" -X POST "$base/admin/$spare_id"
first=$(api "$base/api/admin/root/list-backups" | grep -o '"id":[0-9]*,"name":"first"' | grep -o '[0-9][0-9]*' | head -1)
echo "item 2" > data/app/items/2.changed
api -o /dev/null -X POST -H 'Content-Type: application/json' -d "{\"id\":$first}" "$base/api/admin/root/restore-backup"
idle
echo "after restore, the revoked token: $(status -H "Authorization: Bearer $spare" "$base/api/admin/me")" >> $out
echo "after restore, data/app: $(ls data/app/items | tr '\n' ' ')" >> $out

# A snapshot built by hand whose files cannot both exist is refused.
broken=$(api -X POST "$base/api/admin/root/create-empty-backup" | grep -o '"id":[0-9]*' | grep -o '[0-9][0-9]*')
api -o /dev/null -H 'Content-Type: application/octet-stream' --data-binary x "$base/api/admin/root/add-backup-file/$broken/app/zzz/a"
api -o /dev/null -H 'Content-Type: application/octet-stream' --data-binary y "$base/api/admin/root/add-backup-file/$broken/app/zzz/a/b"
echo "close of a file under a file: $(api -o /dev/null -w '%{http_code}' -H 'Content-Type: application/json' -d "{\"id\":$broken}" "$base/api/admin/root/close-backup")" >> $out

# A restore that fails while writing puts everything back as it was: a file
# where the snapshot needs a folder makes the write fail.
failing=$(api -X POST "$base/api/admin/root/create-empty-backup" | grep -o '"id":[0-9]*' | grep -o '[0-9][0-9]*')
api -o /dev/null -H 'Content-Type: application/octet-stream' --data-binary x "$base/api/admin/root/add-backup-file/$failing/qqq/a"
api -o /dev/null -H 'Content-Type: application/json' -d "{\"id\":$failing}" "$base/api/admin/root/close-backup"
echo "blocker" > data/qqq
api -o /dev/null -X POST -H 'Content-Type: application/json' -d "{\"id\":$failing,\"include-backoffice\":true}" "$base/api/admin/root/restore-backup"
idle
echo "after a failed restore, GET /api/admin/me: $(status -H "Authorization: Bearer $token" "$base/api/admin/me")" >> $out
echo "after a failed restore, data/app: $(ls data/app/items | tr '\n' ' ')" >> $out
echo "after a failed restore, the log: $(grep -c 'was put back as it was' server.log)" >> $out

cd ..
cat test-dir/http.txt

# What result.yaml records: the transcript of the requests above.
mkdir -p assert-dir
cp test-dir/http.txt assert-dir/http.txt
