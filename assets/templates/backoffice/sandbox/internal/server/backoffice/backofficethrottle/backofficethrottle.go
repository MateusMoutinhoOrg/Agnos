package backofficethrottle

import (
	"{{.Module}}/sandbox/api"
)

// Sign-ins and API tokens are counted per client ip, and sign-ins per account
// as well, within a window of WindowSeconds. An attempt is counted *before*
// its password or token is checked — Hit answers the count with it included,
// atomically — so requests sent at the same time can never all pass the
// check before any of them is counted; one that turns out to succeed is taken
// back with Undo. Once a count passes its limit, the attempt is refused with a
// 429, unchecked, until the window closes.
//
// The account is the user the login names — by username or by email, the
// same counter — or, for a login that names nobody, a hash of the login: a key
// never holds what the client typed, so its size is fixed. An ip that signed
// in to an account within TrustedSeconds is counted on its own for that
// account, apart from everyone else: failures sent from other ips never lock
// out the people who use it. The counts live in Deps.RatelimitDeps, in the
// memory of this one process.

// WindowSeconds is how long an attempt counts against its ip or account.
const WindowSeconds = 15 * 60

// TrustedSeconds is how long an ip stays trusted for an account after a
// sign-in to it succeeded from there: thirty days.
const TrustedSeconds = 30 * 24 * 60 * 60

// MaxLoginFailuresPerIp is how many failed sign-ins one client ip may make in
// a window.
const MaxLoginFailuresPerIp = 20

// MaxLoginFailuresPerLogin is how many failed sign-ins one account may take in
// a window, from every ip it is not trusted from together; an ip it is
// trusted from has as many of its own.
const MaxLoginFailuresPerLogin = 10

// MaxTokenFailuresPerIp is how many requests with an invalid API token one
// client ip may make in a window.
const MaxTokenFailuresPerIp = 20

// MaxLoginLength is the longest login, in bytes, a sign-in is checked for:
// the longest an email address may be. A longer one is refused unchecked.
const MaxLoginLength = 254

// MaxPasswordLength is the longest password, in bytes, a sign-in is checked
// for. A longer one is refused unchecked.
const MaxPasswordLength = 1024

// Reservation is one sign-in attempt counted before its password is checked:
// the keys it was counted under, and whether it may be checked at all.
type Reservation struct {
	// Allowed is false when the ip or the account passed its limit: the
	// password is not checked.
	Allowed bool
	ipKey   string
	key     string
	trusted string
}

// loginIpKey names the sign-ins of ip.
func loginIpKey(sandbox *api.Sandbox, ip string) string {
	return "login-ip:" + ip
}

// AccountOf is what the sign-ins of login are counted under: the user it
// names when id is not 0, the SHA-256 of login — trimmed, lower-cased —
// otherwise, so a username and an email of one user share one counter and
// no key grows with what the client typed.
func AccountOf(sandbox *api.Sandbox, login string, id int64) string {
	if id != 0 {
		return "user:" + sandbox.Deps.StringsDeps.FormatInt(id, 10)
	}
	strings := sandbox.Deps.StringsDeps
	return "login:" + sandbox.Deps.HashDeps.Sha256Hex([]byte(strings.ToLower(strings.TrimSpace(login))))
}

// trustedKey names the trust of ip for account.
func trustedKey(sandbox *api.Sandbox, account string, ip string) string {
	return "trusted:" + account + "@" + ip
}

// pairKey names the sign-ins of account from ip, an ip it is trusted from.
func pairKey(sandbox *api.Sandbox, account string, ip string) string {
	return "login-pair:" + account + "@" + ip
}

// tokenIpKey names the API tokens of ip.
func tokenIpKey(sandbox *api.Sandbox, ip string) string {
	return "token-ip:" + ip
}

// RetryAfter is the Retry-After header a refused attempt carries: the length
// of a window, in seconds.
func RetryAfter(sandbox *api.Sandbox) string {
	return sandbox.Deps.StringsDeps.FormatInt(WindowSeconds, 10)
}

// ReserveLogin counts one sign-in to account from ip and tells, on the
// Reservation, whether its password may be checked: false once ip, or the
// account from where ip stands, passed its limit.
func ReserveLogin(sandbox *api.Sandbox, ip string, account string) Reservation {
	limiter := sandbox.Deps.RatelimitDeps
	reservation := Reservation{ipKey: loginIpKey(sandbox, ip), trusted: trustedKey(sandbox, account, ip)}
	if limiter.Hit(reservation.ipKey, WindowSeconds) > MaxLoginFailuresPerIp {
		return reservation
	}
	reservation.key = "login:" + account
	if limiter.Count(reservation.trusted, TrustedSeconds) > 0 {
		reservation.key = pairKey(sandbox, account, ip)
	}
	reservation.Allowed = limiter.Hit(reservation.key, WindowSeconds) <= MaxLoginFailuresPerLogin
	return reservation
}

// LoginSucceeded takes back the attempt reservation counted — a sign-in that
// succeeded is no failure — and trusts its ip for the account for
// TrustedSeconds. The failures counted before it stay: signing in to an
// account of one's own between guesses at another's buys no more guesses.
func LoginSucceeded(sandbox *api.Sandbox, reservation Reservation) {
	limiter := sandbox.Deps.RatelimitDeps
	limiter.Undo(reservation.ipKey)
	limiter.Undo(reservation.key)
	if limiter.Count(reservation.trusted, TrustedSeconds) == 0 {
		limiter.Hit(reservation.trusted, TrustedSeconds)
	}
}

// ReserveToken counts one API token sent from ip and tells whether it may be
// checked: false once ip passed its limit of invalid tokens.
func ReserveToken(sandbox *api.Sandbox, ip string) bool {
	return sandbox.Deps.RatelimitDeps.Hit(tokenIpKey(sandbox, ip), WindowSeconds) <= MaxTokenFailuresPerIp
}

// TokenSucceeded takes back the attempt ReserveToken counted for ip: a valid
// token is no failure.
func TokenSucceeded(sandbox *api.Sandbox, ip string) {
	sandbox.Deps.RatelimitDeps.Undo(tokenIpKey(sandbox, ip))
}
