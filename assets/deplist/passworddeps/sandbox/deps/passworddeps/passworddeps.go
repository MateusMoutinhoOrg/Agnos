package passworddeps

// This package is the sandbox's *copy* of the api a password hashing library
// exposes — the same mechanic as hashdeps and randdeps, for the same reason:
// the sandbox may import nothing but the sandbox, so `crypto/pbkdf2` may not
// appear inside it. The contract is restated here, and the adapter — which
// lives outside the sandbox — is what fills it.
//
// A password is never stored as a plain digest: a fast hash lets a leaked
// database be guessed at billions of tries a second, and one salt shared by
// every user gives equal passwords equal hashes. A hash this contract hands
// out is slow on purpose, carries a random salt of its own and names the
// parameters it was derived with, so Verify reads them back from it.

// Sandbox is the password hashing library injected whole as the
// Deps.Passworddeps field.
type Sandbox struct {
	// Hash derives a hash of password to store, over a fresh random salt,
	// spelled with everything Verify needs to check a password against it.
	// The error reports a random source that could not be read.
	Hash func(password string) (string, error)

	// Verify tells whether password is the one hash was derived from,
	// comparing in constant time. The error reports a hash it cannot read.
	Verify func(hash string, password string) (bool, error)
}
