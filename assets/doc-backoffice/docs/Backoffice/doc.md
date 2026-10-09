# Backoffice

An admin area over http: a login, backoffice users with a `root` or `viewer` role, sessions,
API tokens, [Backups](../Backups/doc.md) of every database, and `/api/admin`, the JSON twin of the
pages. `{{.GeneratorName}} backoffice-init`
wrote every file of it once; each one is the project's from then on.

## Run it

```bash
export {{.SecretEnv}}=$(openssl rand -hex 32)          # optional: signs the sessions; >= 32 chars, never a flag or a file
{{.ProjectName}} add-backoffice-user --username admin --email admin@example.com --role root   # prints the password once
{{.ProjectName}} start-server --insecure-http --addr 127.0.0.1:3000   # local plain http; drop the flag behind https
```

Open `/admin/login`. Without `{{.SecretEnv}}`, `start-server` generates a secret for the run and
warns: it lives only in memory, so every session ends at a restart and no other instance accepts
them. Shorter than 32 characters, or holding fewer than 10 different ones, the server does not
start. `--insecure-http` on every interface (no host in `--addr`) is warned about, and so is a
`.gitignore` that lists neither store under the `--database` folder of the run. The name is the project's name
upper-cased, every other character `_`, then `_BACKOFFICE_SECRET`.

| `start-server` flag | Read by | Effect |
|---|---|---|
| `--insecure-http` | `backoffice-start-server` | Session cookie without `Secure`, no HSTS. Local development only |
| `--allow-x-forwarded-for` | `backoffice-start-server` | Client ip = last `X-Forwarded-For` entry. Only behind one proxy, bound to an address only it reaches (`--addr 127.0.0.1:3000`) |

`backoffice-start-server` is a cli middleware (priority `50`) in front of `start-server`: it reads the
secret and both flags onto `sandbox.Config` (`api.BackofficeConfig`,
`sandbox/api/backofficeconfig.go`). `start-server`'s own files are untouched.

## Where it lives

| Path | What |
|---|---|
| `sandbox/internal/server/backoffice/backofficeauth/` | roles, secret, password hashes, login, session JWT and cookie, bearer token |
| `sandbox/internal/server/backoffice/backofficeusers/` | list, filter, page, add, set, remove; validation and uniqueness |
| `sandbox/internal/server/backoffice/backofficeapitokens/` | API tokens: add, list, revoke, resolve |
| `sandbox/internal/server/backoffice/backofficethrottle/` | login and token rate limits, counted before the check |
| `sandbox/internal/server/backoffice/backofficehttp/` | client ip, security headers, same-origin |
| `sandbox/internal/server/backoffice/backofficerender/` | the pages, from `assets/backoffice/*.html` |
| `sandbox/internal/server/backoffice/backofficeapi/` | the JSON documents of `/api/admin` |
| `sandbox/internal/databases/backoffice_db/` | `backoffice-user` (+ nested `session`) and `api-token`; the store is `./data/backofficedb` (inside `--database`), gitignored |
| `sandbox/internal/snapshots/`, `sandbox/internal/databases/backup/`, `sandbox/internal/server/backoffice/backofficesnapshots/` | the [Backups](../Backups/doc.md); the store is `./data/backup`, gitignored |
| `sandbox/internal/routeprops/backoffice.go` | `props.ClientIp`, `props.User`, `props.Session`, `props.ApiToken` |
| `sandbox/internal/commands/backoffice/`, `sandbox/internal/commands/middleware/` | `add-backoffice-user`; `backoffice-start-server`, the middleware in front of `start-server` |
| `assets/backoffice/*.html`, `assets/front/backoffice/backoffice.js` | page templates (`html/template`, through `Deps.EmbedDeps.RenderHTMLTemplate`: every value escaped for its context) and their script |

## Routes

Middlewares run lowest priority first; a page answers or the chain goes on.

| Route | Priority | Matches | Does |
|---|---|---|---|
| `backoffice-maintenance` | 6 | every path | `503` while a backup is restored, and to writes while one is taken ([Backups](../Backups/doc.md)) |
| `backoffice-client-ip` | 7 | every path | sets `props.ClientIp` |
| `backoffice-security-headers` | 8 | every path | `/admin`, `/api/admin`: CSP (no inline script), XFO, nosniff, Referrer-Policy, no-store, HSTS; any other: nosniff, `X-Frame-Options: SAMEORIGIN`, Referrer-Policy, which a route may override |
| `backoffice-same-origin` | 9 | `/admin` | `403` to an `Origin` that is not the `Host` |
| `backoffice-session-auth` | 10 | `/admin` but `/admin/login` | session cookie → `props.User`, or the login page |
| `backoffice-root-guard` | 11 | `/admin/root` | `403` to a non-root |
| `backoffice-api-token-auth` | 10 | `/api/admin` | `Authorization: Bearer <token>` → `props.User`, `props.ApiToken` |
| `backoffice-api-root-guard` | 11 | `/api/admin/root` | `403` to a non-root |

Pages (`/admin/...`): `login` (GET form, POST), `logout` (POST), `home` (`/admin` itself `303`s there), `list-backoffice-users`,
`list-backoffice-api-tokens`, `add-backoffice-api-token` (GET form, POST), `revoke-backoffice-api-token/{id}` (POST);
root only (`/admin/root/...`): `add-backoffice-user` (GET form, POST), `set-backoffice-user/{id}` (GET form, POST),
`remove-backoffice-user/{id}` (POST), and the backups (`list-backups`…, in [Backups](../Backups/doc.md)). Each is one route per method, named after the surface it
answers: `backoffice-home`, `backoffice-login`, `list-backoffice-users-page` (GET html),
`add-backoffice-user-form` (POST html). An action answers `303` to its list with a `?notice=`.

JSON (`/api/admin/...`, body `application/json`, every action a `POST`): `me` (GET),
`list-backoffice-users` `{search?, role?, page?, limit?}`, `get-backoffice-user` `{id}`; root only:
`root/add-backoffice-user` (`201`), `root/set-backoffice-user` `{id, username, email, role, password?}`,
`root/remove-backoffice-user` `{id}`. `role` travels as `"root"`/`"viewer"`. Failures are the
`{"error","field"}` of `sandbox/internal/server/errors/handle_*.go`: `401` bad token, `429` too
many, `403` not root, `404` unknown id, `400` invalid. `{{.GeneratorName}} explain-route GET /admin/home`
shows which route answers. Every route above carries `private: true`, so `/openapi.json` leaves
it out; [Routes](../Routes/doc.md) still lists it.

## Rules

- Passwords: PBKDF2-HMAC-SHA256, 600k iterations, a salt each (`Deps.PasswordDeps`); 12 to 1024
  characters, at least 5 different ones, none of a list of common ones, not the username nor the
  email. At most 4 are checked at once, the rest wait. `add-backoffice-user` generates one and
  prints it once, so none travels in argv.
- Username and email are unique across both columns, ignoring case: checked and written under one
  lock, and an `add-backoffice-user` racing the server loses its insert when it came second. A
  root cannot remove itself; the last root cannot be demoted, nor removed by a root demoted
  meanwhile. A new password ends every session of the user and revokes
  their API tokens.
- Session: an HS256 JWT (`Deps.JwtDeps`) in the `backoffice_session` cookie (HttpOnly, SameSite=Strict,
  Secure, `Path=/admin`, 30 min), its `jti` a `session` record under the user and its `ip` the client's.
  A page of the application never carries it.
  Logout or removing the user deletes it.
- API token: `bo_` + 64 hex (`Deps.RandDeps`), shown once; only its SHA-256 is stored
  (`FindApiTokenByTokenSha256`). Creating one asks the user's password again. Expires in
  7/30/60/90/365 days, on a date, or never — a viewer's within 90 days; may be limited to ips and
  CIDR ranges, compared as parsed addresses. It acts as its owner, role read fresh per request;
  its last use is written at most once a minute, or when its ip changes.
- A viewer sees every other user's email masked (`a***@example.com`), and searches usernames.
- Rate limits (`Deps.RatelimitDeps`, 15 min), each attempt counted before it is checked and taken
  back once it succeeds: sign-in 20 failures per ip, and 10 per account — its username and email
  together — from the ips it did not sign in from in the last 30 days, each such ip keeping 10 of
  its own; tokens 20 invalid per ip. Both answer `429` with `Retry-After`. A sign-in body is at
  most 4 KiB. The counters live in the memory of one process, at most 65536 keys: a restart forgets
  them, and two instances behind a load balancer allow twice the limit.
- Page script lives in `assets/front/backoffice/backoffice.js`: the CSP runs no inline script, so
  confirm a form with `data-confirm`, never an `on*` attribute.
- Everything is named `backoffice*`, so application users can take the plain names later.

`{{.GeneratorName}} backoffice-purge` removes every file above and keeps the layers, the deps,
`./data/backofficedb` and `./data/backup`.
