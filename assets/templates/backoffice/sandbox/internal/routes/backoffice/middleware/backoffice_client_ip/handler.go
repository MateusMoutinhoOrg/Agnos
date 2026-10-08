package client_ip

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeguard"
)

// InternalPureHandler runs in front of every ANY /*, on the lowest rung of the
// chain: it is a middleware. It puts the ip the request came from on
// props.ClientIp — the connection's own, or, when start-server trusts
// X-Forwarded-For, the one the reverse proxy in front appended — and declines,
// so every route after it reads one ip, worked out once.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	props.ClientIp = backofficeguard.ClientIp(sandbox, entries.XClientIp, entries.XForwardedFor)
	return nil
}
