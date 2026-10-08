package routeprops

import (
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
)

// Backoffice is the part of RouteProps the backoffice middlewares hand on to
// the routes after them, embedded so each field is read as props.<Field>.
//
// Written once by `{{.GeneratorName}} backoffice-init` and the project's from
// then on; `{{.GeneratorName}} backoffice-purge` removes it.
type Backoffice struct {
	// ClientIp is the ip the request came from, as the backoffice-client-ip middleware
	// worked it out: the connection's own, or the one the reverse proxy in
	// front appended to X-Forwarded-For when start-server trusts it.
	ClientIp string
	// User is the backoffice user the backoffice-session-auth middleware
	// authenticated from the session cookie, or the one the
	// backoffice-api-token-auth middleware authenticated from an API
	// token; nil when none.
	User *backoffice_db.BackofficeUserRecord
	// Session is the session of User the session cookie names, nil when
	// none — and always nil on /api/admin, which reads no cookie.
	Session *backoffice_db.BackofficeUserSessionRecord
	// ApiToken is the API token of User the Authorization header carried
	// on /api/admin, nil when none — and always nil on /admin.
	ApiToken *backoffice_db.ApiTokenRecord
}
