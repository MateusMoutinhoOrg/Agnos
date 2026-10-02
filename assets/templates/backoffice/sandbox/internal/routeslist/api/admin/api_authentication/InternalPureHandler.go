package api_authentication

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/routeio"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficethrottle"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficetokens"
)

// InternalPureHandler runs in front of every ANY /api/admin/{*Rest}, on a
// lower rung of the chain: it is a middleware.
//
// An `Authorization: Bearer <token>` header holding an API token — created on
// the backoffice page, not revoked, not expired, and allowed from the client
// ip this request came from — puts the user it belongs to on props.User and
// the token on props.ApiToken and declines, so the route after it runs.
// Anything else is refused with a 401 in JSON. Once the client ip sent
// backofficethrottle.MaxTokenFailuresPerIp invalid tokens, every token it
// sends is refused with a 429, unchecked, until the window closes. The
// session cookie is never
// read here, so a browser carrying one cannot be made to call these routes by
// another site, and the api issues no token of its own.
func InternalPureHandler(sandbox *api.Sandbox, props *routeprops.RouteProps, entries *Entries, response *serverdeps.Response) error {
	token := backofficeauth.BearerToken(sandbox, entries.Authorization)
	if token != "" && !backofficethrottle.TokenAllowed(sandbox, props.ClientIp) {
		response.SetHeader("Retry-After", backofficethrottle.RetryAfter(sandbox))
		return routeio.Fail(sandbox, api.StatusTooManyRequests, "authorization", "too many invalid tokens from this ip, try again later")
	}
	user, apiToken, ok, err := backofficetokens.Resolve(sandbox, token, props.ClientIp)
	if err != nil {
		return err
	}
	if ok {
		props.User = &user
		props.ApiToken = &apiToken
		return nil
	}
	if token != "" {
		backofficethrottle.TokenFailed(sandbox, props.ClientIp)
	}

	response.SetHeader("WWW-Authenticate", "Bearer")
	if token == "" {
		return routeio.Fail(sandbox, api.StatusUnauthorized, "authorization", "send an API token, created on "+backofficetokens.ListPath+", in the Authorization header, after Bearer")
	}
	return routeio.Fail(sandbox, api.StatusUnauthorized, "authorization", "the token is invalid, expired, revoked or not allowed from this ip")
}
