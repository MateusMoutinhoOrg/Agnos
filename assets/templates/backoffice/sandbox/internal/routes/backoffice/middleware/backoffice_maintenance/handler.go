package backoffice_maintenance

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/snapshots"
)

// RetryAfterSeconds is the Retry-After a request refused here carries.
const RetryAfterSeconds = "5"

// Handle runs in front of every ANY /*, on the lowest rung of the
// chain: it is a middleware.
//
// While a restore — or the roll-back of one the last run did not finish —
// writes the databases, every request is refused with a 503: none reads a
// half-written database, nor writes into one about to be replaced. While a
// snapshot reads them file by file, a request that may write — any method but
// GET, HEAD and OPTIONS — is refused the same way, so every file of the
// snapshot is taken at one instant. Anything else declines, so the route
// after it runs. A request already past this rung when a job starts is not
// stopped, and neither is a write made outside the server, by a command line.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	if snapshots.Restoring(sandbox) {
		response.SetHeader("Retry-After", RetryAfterSeconds)
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusServiceUnavailable, "", "a backup is being restored: try again in a few seconds")
	}
	if snapshots.Creating(sandbox) && !reads(sandbox, input.XRequestMethod) {
		response.SetHeader("Retry-After", RetryAfterSeconds)
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusServiceUnavailable, "", "a backup is being taken: try again in a few seconds")
	}
	return nil
}

// reads tells whether method is one that only reads: GET, HEAD or OPTIONS.
func reads(sandbox *api.Sandbox, method string) bool {
	return method == "GET" || method == "HEAD" || method == "OPTIONS"
}
