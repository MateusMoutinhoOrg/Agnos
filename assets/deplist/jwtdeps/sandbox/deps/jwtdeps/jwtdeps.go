package jwtdeps

// This package is the sandbox's *copy* of the api a JSON Web Token library
// exposes — the same mechanic as hashdeps, embeddeps and serverdeps, for the
// same reason: the sandbox may import nothing but the sandbox, so neither a
// jwt library nor `crypto/hmac` may appear inside it. The contract is restated
// here, and the adapter — which lives outside the sandbox — is what fills it.
//
// Only builtin types cross this boundary: an instant is a count of seconds
// since the Unix epoch, never a `time.Time`, and the claims are the one flat
// struct below, never the library's own claims type.

// Sandbox is the JSON Web Token library injected whole as the Deps.Jwtdeps
// field. Every token it signs and every token it accepts is HS256, keyed by
// the secret handed to the call.
type Sandbox struct {
	// Sign returns the compact HS256 token carrying claims, signed with
	// secret. The error reports an empty secret or a claims set the library
	// could not encode.
	Sign func(claims Claims, secret string) (string, error)

	// Parse checks token against secret and returns the claims it carries.
	// The error reports a token that is malformed, signed by another key or
	// another algorithm than HS256, or whose ExpiresAt has passed — so a
	// nil error is a token that is valid right now.
	Parse func(token string, secret string) (Claims, error)
}

// Claims is the set of registered claims a token carries, plus the private
// `ip` claim.
type Claims struct {
	// Id is the `jti` claim: the one token among every token issued.
	Id string
	// Subject is the `sub` claim: who the token was issued for.
	Subject string
	// IssuedAt is the `iat` claim, in seconds since the Unix epoch.
	IssuedAt int64
	// ExpiresAt is the `exp` claim, in seconds since the Unix epoch. A token
	// is refused by Parse from that instant on.
	ExpiresAt int64
	// Ip is the private `ip` claim: the client ip the token was issued to.
	Ip string
}
