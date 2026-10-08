package backofficethrottle

import (
	"{{.Module}}/sandbox/api"
)

// Failed sign-ins and failed API tokens are counted per client ip, and failed
// sign-ins per login as well, within a window of WindowSeconds. Once a count
// reaches its limit, the next attempt is refused with a 429 before any
// password or token is checked, until the window closes: guessing a password
// from one ip, or one account's password from many, slows to a few tries
// every quarter of an hour. The counts live in Deps.RatelimitDeps, in memory.

// WindowSeconds is how long a failure counts against its ip or login.
const WindowSeconds = 15 * 60

// MaxLoginFailuresPerIp is how many failed sign-ins one client ip may make in
// a window.
const MaxLoginFailuresPerIp = 20

// MaxLoginFailuresPerLogin is how many failed sign-ins one login — a username
// or an email, regardless of case — may take in a window, from any ip.
const MaxLoginFailuresPerLogin = 10

// MaxTokenFailuresPerIp is how many requests with an invalid API token one
// client ip may make in a window.
const MaxTokenFailuresPerIp = 20

// loginIpKey names the failed sign-ins of ip.
func loginIpKey(sandbox *api.Sandbox, ip string) string {
	return "login-ip:" + ip
}

// loginKey names the failed sign-ins of login.
func loginKey(sandbox *api.Sandbox, login string) string {
	strings := sandbox.Deps.StringsDeps
	return "login:" + strings.ToLower(strings.TrimSpace(login))
}

// tokenIpKey names the failed API tokens of ip.
func tokenIpKey(sandbox *api.Sandbox, ip string) string {
	return "token-ip:" + ip
}

// RetryAfter is the Retry-After header a refused attempt carries: the length
// of a window, in seconds.
func RetryAfter(sandbox *api.Sandbox) string {
	return sandbox.Deps.StringsDeps.FormatInt(WindowSeconds, 10)
}

// LoginAllowed tells whether a sign-in as login from ip may be checked, false
// once either reached its limit of failures.
func LoginAllowed(sandbox *api.Sandbox, ip string, login string) bool {
	limiter := sandbox.Deps.RatelimitDeps
	return limiter.Count(loginIpKey(sandbox, ip), WindowSeconds) < MaxLoginFailuresPerIp &&
		limiter.Count(loginKey(sandbox, login), WindowSeconds) < MaxLoginFailuresPerLogin
}

// LoginFailed counts one failed sign-in as login from ip.
func LoginFailed(sandbox *api.Sandbox, ip string, login string) {
	sandbox.Deps.RatelimitDeps.Hit(loginIpKey(sandbox, ip), WindowSeconds)
	sandbox.Deps.RatelimitDeps.Hit(loginKey(sandbox, login), WindowSeconds)
}

// LoginSucceeded forgets the failed sign-ins of login. The ones of the ip
// stay: signing in to an account of one's own between guesses at another's
// buys no more guesses.
func LoginSucceeded(sandbox *api.Sandbox, login string) {
	sandbox.Deps.RatelimitDeps.Reset(loginKey(sandbox, login))
}

// TokenAllowed tells whether an API token sent from ip may be checked, false
// once ip reached its limit of invalid tokens.
func TokenAllowed(sandbox *api.Sandbox, ip string) bool {
	return sandbox.Deps.RatelimitDeps.Count(tokenIpKey(sandbox, ip), WindowSeconds) < MaxTokenFailuresPerIp
}

// TokenFailed counts one invalid API token sent from ip.
func TokenFailed(sandbox *api.Sandbox, ip string) {
	sandbox.Deps.RatelimitDeps.Hit(tokenIpKey(sandbox, ip), WindowSeconds)
}
