package backoffice_api_token_auth

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficethrottle"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapitokens"
)

// Handle runs in front of every ANY /api/admin/{*Rest}, on a
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
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	token := backofficeauth.BearerToken(sandbox, input.Authorization)
	if token != "" && !backofficethrottle.TokenAllowed(sandbox, props.ClientIp) {
		response.SetHeader("Retry-After", backofficethrottle.RetryAfter(sandbox))
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusTooManyRequests, "authorization", "too many invalid tokens from this ip, try again later")
	}
	user, apiToken, ok, err := backofficeapitokens.Resolve(sandbox, token, props.ClientIp)
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
		return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "authorization", "send an API token, created on "+backofficeapitokens.ListPath+", in the Authorization header, after Bearer")
	}
	return sandbox.Deps.OpinionatedAgnosServer.Fail(api.StatusUnauthorized, "authorization", "the token is invalid, expired, revoked or not allowed from this ip")
}
