package passworddeps

import (
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strconv"
	"strings"

	passworddeps "{{.Module}}/sandbox/deps/passworddeps"

	"{{.Module}}/sandbox/deps"
)

// scheme names the algorithm at the head of every hash this adapter writes.
const scheme = "pbkdf2-sha256"

// iterations is how many rounds of HMAC-SHA256 a new hash takes, the count
// OWASP recommends for PBKDF2-HMAC-SHA256. A stored hash carries its own
// count, so raising this one never invalidates an older hash.
const iterations = 600_000

// maxIterations is the most rounds Verify accepts from a stored hash, so a
// corrupted one cannot pin a request on the CPU.
const maxIterations = 10_000_000

// saltBytes is how many random bytes salt each hash.
const saltBytes = 16

// keyBytes is how long the derived key is.
const keyBytes = 32

// encoding spells the salt and the key: standard base64 without padding.
var encoding = base64.RawStdEncoding

// hash fills passworddeps.Sandbox.Hash:
// "pbkdf2-sha256$<iterations>$<salt>$<key>", salt and key in base64.
func hash(password string) (string, error) {
	salt := make([]byte, saltBytes)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, iterations, keyBytes)
	if err != nil {
		return "", err
	}
	return strings.Join([]string{
		scheme,
		strconv.Itoa(iterations),
		encoding.EncodeToString(salt),
		encoding.EncodeToString(key),
	}, "$"), nil
}

// verify fills passworddeps.Sandbox.Verify: password derived again with the
// salt and the count hash carries, compared with its key in constant time.
func verify(hash string, password string) (bool, error) {
	parts := strings.Split(hash, "$")
	if len(parts) != 4 || parts[0] != scheme {
		return false, errors.New("not a " + scheme + " password hash")
	}
	rounds, err := strconv.Atoi(parts[1])
	if err != nil || rounds < 1 || rounds > maxIterations {
		return false, errors.New("password hash has an invalid iteration count")
	}
	salt, err := encoding.DecodeString(parts[2])
	if err != nil {
		return false, errors.New("password hash has an invalid salt")
	}
	stored, err := encoding.DecodeString(parts[3])
	if err != nil || len(stored) == 0 {
		return false, errors.New("password hash has an invalid key")
	}
	key, err := pbkdf2.Key(sha256.New, password, salt, rounds, len(stored))
	if err != nil {
		return false, err
	}
	return subtle.ConstantTimeCompare(key, stored) == 1, nil
}

// Bind fills deps.Deps.Passworddeps with PBKDF2-HMAC-SHA256 from the standard
// library's crypto/pbkdf2.
func Bind(deps *deps.Deps) {
	deps.Passworddeps = passworddeps.Sandbox{
		Hash: func(password string) (string, error) {
			return hash(password)
		},
		Verify: func(hash string, password string) (bool, error) {
			return verify(hash, password)
		},
	}
}
