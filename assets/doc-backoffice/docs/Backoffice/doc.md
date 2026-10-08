# Backoffice

An admin area over http: a login, backoffice users with a `root` or `viewer` role, sessions,
API tokens, and `/api/admin`, the JSON twin of the pages. `{{.GeneratorName}} backoffice-init`
wrote every file of it once; each one is the project's from then on.

## Run it

```bash
export {{.SecretEnv}}=$(openssl rand -hex 32)          # optional: signs the sessions; >= 32 chars, never a flag or a file
{{.ProjectName}} add-backoffice-user --username admin --email admin@example.com --role root   # prints the password once
{{.ProjectName}} start-server --insecure-http                 # local plain http; drop the flag behind https
```

Open `/admin/login`. Without `{{.SecretEnv}}`, `start-server` generates a secret for the run and
warns: it lives only in memory, so every session ends at a restart and no other instance accepts
them. Shorter than 32 characters, the server does not start. The name is the project's name
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
| `sandbox/internal/server/backoffice/backofficethrottle/` | login and token rate limits |
| `sandbox/internal/server/backoffice/backofficehttp/` | client ip, security headers, same-origin |
| `sandbox/internal/server/backoffice/backofficerender/` | the pages, from `assets/backoffice/*.html` |
| `sandbox/internal/server/backoffice/backofficeapi/` | the JSON documents of `/api/admin` |
| `sandbox/internal/databases/backoffice_db/` | `backoffice-user` (+ nested `session`) and `api-token`; the store is `./data/backofficedb`, gitignored |
| `sandbox/internal/routeprops/backoffice.go` | `props.ClientIp`, `props.User`, `props.Session`, `props.ApiToken` |
| `sandbox/internal/commands/backoffice/`, `sandbox/internal/commands/middleware/` | `add-backoffice-user`; `backoffice-start-server`, the middleware in front of `start-server` |
| `assets/backoffice/*.html`, `assets/front/backoffice/backoffice.js` | page templates (`text/template`, values escaped with `html`) and their script |

## Routes

Middlewares run lowest priority first; a page answers or the chain goes on.

| Route | Priority | Matches | Does |
|---|---|---|---|
| `backoffice-client-ip` | 7 | every path | sets `props.ClientIp` |
| `backoffice-security-headers` | 8 | `/admin`, `/api/admin` | CSP (no inline script), XFO, nosniff, Referrer-Policy, no-store, HSTS |
| `backoffice-same-origin` | 9 | `/admin` | `403` to an `Origin` that is not the `Host` |
| `backoffice-session-auth` | 10 | `/admin` but `/admin/login` | session cookie → `props.User`, or the login page |
| `backoffice-root-guard` | 11 | `/admin/root` | `403` to a non-root |
| `backoffice-api-token-auth` | 10 | `/api/admin` | `Authorization: Bearer <token>` → `props.User`, `props.ApiToken` |
| `backoffice-api-root-guard` | 11 | `/api/admin/root` | `403` to a non-root |

Pages (`/admin/...`): `login` (POST), `logout` (POST), `home`, `list-backoffice-users`,
`list-backoffice-api-tokens`, `add-backoffice-api-token` (GET form, POST), `revoke-backoffice-api-token/{id}` (POST);
root only (`/admin/root/...`): `add-backoffice-user` (GET form, POST), `set-backoffice-user/{id}` (GET form, POST),
`remove-backoffice-user/{id}` (POST). Each is one route per method, named after the surface it
answers: `backoffice-home`, `backoffice-login`, `list-backoffice-users-page` (GET html),
`add-backoffice-user-form` (POST html). An action answers `303` to its list with a `?notice=`.

JSON (`/api/admin/...`, body `application/json`, every action a `POST`): `me` (GET),
`list-backoffice-users` `{search?, role?, page?, limit?}`, `get-backoffice-user` `{id}`; root only:
`root/add-backoffice-user` (`201`), `root/set-backoffice-user` `{id, username, email, role, password?}`,
`root/remove-backoffice-user` `{id}`. `role` travels as `"root"`/`"viewer"`. Failures are the
`{"error","field"}` of `sandbox/internal/server/errors/handle_*.go`: `401` bad token, `429` too
many, `403` not root, `404` unknown id, `400` invalid. `{{.GeneratorName}} explain-route GET /admin/home`
shows which route answers.

## Rules

- Passwords: PBKDF2-HMAC-SHA256, 600k iterations, a salt each (`Deps.PasswordDeps`); at least 8 characters.
  `add-backoffice-user` generates one and prints it once, so none travels in argv.
- Username and email are unique across both columns, ignoring case. A root cannot remove itself;
  the last root cannot be demoted. A new password ends every session of the user and revokes
  their API tokens.
- Session: an HS256 JWT (`Deps.JwtDeps`) in the `backoffice_session` cookie (HttpOnly, SameSite=Strict,
  Secure, 30 min), its `jti` a `session` record under the user and its `ip` the client's.
  Logout or removing the user deletes it.
- API token: `bo_` + 64 hex (`Deps.RandDeps`), shown once; only its SHA-256 is stored
  (`FindApiTokenByTokenSha256`). Expires in 7/30/60/90/365 days, on a date, or never; may be
  limited to ips. It acts as its owner, role read fresh per request.
- Rate limits (`Deps.RatelimitDeps`, in memory, 15 min): login 20 failures per ip and 10 per
  login; tokens 20 invalid per ip. Both answer `429` with `Retry-After`.
- Page script lives in `assets/front/backoffice/backoffice.js`: the CSP runs no inline script, so
  confirm a form with `data-confirm`, never an `on*` attribute.
- Everything is named `backoffice*`, so application users can take the plain names later.

`{{.GeneratorName}} backoffice-purge` removes every file above and keeps the layers, the deps and
`./data/backofficedb`.
