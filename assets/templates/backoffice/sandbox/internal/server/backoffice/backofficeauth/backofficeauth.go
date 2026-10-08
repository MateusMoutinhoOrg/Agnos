package backofficeauth

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/jwtdeps"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
)

// CookieName is the cookie the session token travels in. The backoffice-session-auth
// route declares a cookie parameter under the same key.
const CookieName = "backoffice_session"

// SessionSeconds is how long a session token is valid after login.
const SessionSeconds = 30 * 60

// SecretEnv is the environment variable that may hold the secret that signs
// session tokens: the project's name upper-cased, every character
// but a letter or a digit turned into "_", then "_BACKOFFICE_SECRET" —
// MEUSITE_BACKOFFICE_SECRET for a project named meusite. It follows the name the project is built under, so
// renaming the project renames the variable. It is never a flag — every user
// of the machine reads a command line — nor a file, which can end up committed
// with the code.
func SecretEnv(sandbox *api.Sandbox) string {
	name := []byte(sandbox.Deps.StringsDeps.ToUpper(sandbox.Config.ProjectName))
	for i, char := range name {
		if !(char >= 'A' && char <= 'Z') && !(char >= '0' && char <= '9') {
			name[i] = '_'
		}
	}
	return string(name) + "_BACKOFFICE_SECRET"
}

// MinSecretLength is the fewest characters the secret may have.
const MinSecretLength = 32

// GeneratedSecretBytes is how many random bytes the secret ReadSecret
// generates carries, in hex — twice MinSecretLength.
const GeneratedSecretBytes = 32

// Role is what the role column of a backofficeuser stands for.
type Role int64

const (
	// RoleRoot is the root role.
	RoleRoot Role = iota
	// RoleViewer is the viewer role.
	RoleViewer
)

// Roles is every role a backoffice user may hold, in the order they are offered.
func Roles(sandbox *api.Sandbox) []Role {
	return []Role{RoleRoot, RoleViewer}
}

// RoleName is the display name of a role, "unknown" for one Roles does not hold.
func RoleName(sandbox *api.Sandbox, role Role) string {
	switch role {
	case RoleRoot:
		return "root"
	case RoleViewer:
		return "viewer"
	}
	return "unknown"
}

// ParseRole is the role RoleName spells as name, or false when none does.
func ParseRole(sandbox *api.Sandbox, name string) (Role, bool) {
	for _, role := range Roles(sandbox) {
		if RoleName(sandbox, role) == name {
			return role, true
		}
	}
	return RoleRoot, false
}

// ValidRole tells whether role is one of Roles.
func ValidRole(sandbox *api.Sandbox, role Role) bool {
	for _, known := range Roles(sandbox) {
		if known == role {
			return true
		}
	}
	return false
}

// nowSeconds is the current time in seconds since the Unix epoch.
func nowSeconds(sandbox *api.Sandbox) int64 {
	return sandbox.Deps.StdDeps.Now() / 1_000_000_000
}

// ReadSecret is the secret that signs session tokens, read from the SecretEnv
// environment variable. Unset or empty, it is GeneratedSecretBytes random
// bytes generated for this run alone, and the bool is true: it lives only in
// memory, so every session ends when the server restarts. A value shorter
// than MinSecretLength is refused, saying how to set it, rather than replaced,
// so a mistyped secret never turns into a generated one in silence.
func ReadSecret(sandbox *api.Sandbox) (string, bool, error) {
	env := SecretEnv(sandbox)
	secret := sandbox.Deps.EnvDeps.Getenv(env)
	if secret == "" {
		generated, err := sandbox.Deps.RandDeps.Hex(GeneratedSecretBytes)
		return generated, true, err
	}
	if len(secret) < MinSecretLength {
		return "", false, sandbox.Deps.StdDeps.Errorf("the %s environment variable holds fewer than %d characters: set it to a random secret (openssl rand -hex 32), or unset it to have one generated for each run; it signs the backoffice sessions", env, MinSecretLength)
	}
	return secret, false, nil
}

// HashPassword is how a backoffice password is stored: a salted, deliberately
// slow hash of it, with a salt of its own, so equal passwords never share a
// hash and a leaked database is slow to guess at.
func HashPassword(sandbox *api.Sandbox, password string) (string, error) {
	return sandbox.Deps.PasswordDeps.Hash(password)
}

// FindUserByUsernameOrEmail looks a backoffice user up by username or email. It reads every
// user once whichever the login is, so a username, an email and a login that
// names nobody take the same time.
func FindUserByUsernameOrEmail(sandbox *api.Sandbox, login string) (backoffice_db.BackofficeUserRecord, bool, error) {
	users, err := backoffice_db.New(sandbox).ListBackofficeUsers(backoffice_db.BackofficeUserFilter{})
	if err != nil {
		return backoffice_db.BackofficeUserRecord{}, false, err
	}
	for _, user := range users {
		if user.Username == login || user.Email == login {
			return user, true, nil
		}
	}
	return backoffice_db.BackofficeUserRecord{}, false, nil
}

// Authenticate answers the user whose login (username or email) and password
// match, or false when either does not. A login that names nobody still costs
// one password hash, so how long a refusal takes never tells an unknown login
// from a wrong password.
func Authenticate(sandbox *api.Sandbox, login string, password string) (backoffice_db.BackofficeUserRecord, bool, error) {
	none := backoffice_db.BackofficeUserRecord{}
	user, ok, err := FindUserByUsernameOrEmail(sandbox, login)
	if err != nil {
		return none, false, err
	}
	if !ok {
		_, err = HashPassword(sandbox, password)
		return none, false, err
	}
	match, err := sandbox.Deps.PasswordDeps.Verify(user.PasswordHash, password)
	if err != nil || !match {
		return none, false, err
	}
	return user, true, nil
}

// IssueSessionJWT opens a session for user, whose request came from the client ip
// ip, and signs its token, valid for SessionSeconds and bound to ip. The session is a sessions record under user, living as long
// as the token; its id travels as the token's `jti`, so the authentication
// middleware can tell whether that one session is still open. The user's
// expired sessions are dropped first, so they never pile up.
func IssueSessionJWT(sandbox *api.Sandbox, user backoffice_db.BackofficeUserRecord, ip string) (string, error) {
	now := nowSeconds(sandbox)
	err := dropExpired(sandbox, user.Id, now)
	if err != nil {
		return "", err
	}

	session, err := backoffice_db.New(sandbox).AddBackofficeUserSession(user.Id, backoffice_db.BackofficeUserSessionInput{
		ExpiresAt: now + SessionSeconds,
	})
	if err != nil {
		return "", err
	}

	return sandbox.Deps.JwtDeps.Sign(jwtdeps.Claims{
		Id:        sandbox.Deps.StringsDeps.FormatInt(session.Id, 10),
		Subject:   sandbox.Deps.StringsDeps.FormatInt(user.Id, 10),
		IssuedAt:  now,
		ExpiresAt: now + SessionSeconds,
		Ip:        ip,
	}, sandbox.Config.SessionSecret)
}

// SessionCookie is the Set-Cookie value carrying token: HttpOnly,
// SameSite=Strict, Secure unless start-server serves plain http, on every
// path, expiring with the token.
func SessionCookie(sandbox *api.Sandbox, token string) string {
	return sandbox.Deps.StdDeps.Sprintf("%s=%s; Path=/; Max-Age=%d; HttpOnly; SameSite=Strict%s", CookieName, token, SessionSeconds, secureAttribute(sandbox))
}

// ClearedCookie is the Set-Cookie value that removes the session cookie.
func ClearedCookie(sandbox *api.Sandbox) string {
	return sandbox.Deps.StdDeps.Sprintf("%s=; Path=/; Max-Age=0; HttpOnly; SameSite=Strict%s", CookieName, secureAttribute(sandbox))
}

// secureAttribute is the Secure attribute of the session cookie, which keeps
// the browser from sending it over plain http; "" when start-server serves
// plain http, for local development.
func secureAttribute(sandbox *api.Sandbox) string {
	if sandbox.Config.InsecureHttp {
		return ""
	}
	return "; Secure"
}

// BearerToken is the token an Authorization header carries in the Bearer
// scheme, the scheme matched regardless of case, or "" when it carries none.
// The api-authentication middleware reads the API token from there, where the
// authentication one reads the session token from the session cookie.
func BearerToken(sandbox *api.Sandbox, authorization string) string {
	fields := sandbox.Deps.StringsDeps.Fields(authorization)
	if len(fields) != 2 || sandbox.Deps.StringsDeps.ToLower(fields[0]) != "bearer" {
		return ""
	}
	return fields[1]
}

// dropExpired deletes every session of the user with id userId that expired
// by now: its token is refused anyway, so the record is only garbage.
func dropExpired(sandbox *api.Sandbox, userId int64, now int64) error {
	db := backoffice_db.New(sandbox)
	sessions, err := db.ListBackofficeUserSessions(userId)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		if session.ExpiresAt <= now {
			err = backoffice_db.RemoveBackofficeUserSession(sandbox, db, userId, session.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// findSession is the session with id sessionId of the user with id userId.
func findSession(sandbox *api.Sandbox, userId int64, sessionId int64) (backoffice_db.BackofficeUserSessionRecord, bool) {
	sessions, err := backoffice_db.New(sandbox).ListBackofficeUserSessions(userId)
	if err != nil {
		return backoffice_db.BackofficeUserSessionRecord{}, false
	}
	for _, session := range sessions {
		if session.Id == sessionId {
			return session, true
		}
	}
	return backoffice_db.BackofficeUserSessionRecord{}, false
}

// ResolveSession answers the user a session token was issued for and the
// session it names, or false when the token is invalid or expired, was issued
// to another client ip than ip, names a session that was closed by a logout,
// or its user no longer exists.
func ResolveSession(sandbox *api.Sandbox, token string, ip string) (backoffice_db.BackofficeUserRecord, backoffice_db.BackofficeUserSessionRecord, bool) {
	none := func() (backoffice_db.BackofficeUserRecord, backoffice_db.BackofficeUserSessionRecord, bool) {
		return backoffice_db.BackofficeUserRecord{}, backoffice_db.BackofficeUserSessionRecord{}, false
	}
	if token == "" || ip == "" {
		return none()
	}
	claims, err := sandbox.Deps.JwtDeps.Parse(token, sandbox.Config.SessionSecret)
	if err != nil || claims.Ip != ip {
		return none()
	}
	userId, err := sandbox.Deps.StringsDeps.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return none()
	}
	sessionId, err := sandbox.Deps.StringsDeps.ParseInt(claims.Id, 10, 64)
	if err != nil {
		return none()
	}
	user, ok := backoffice_db.New(sandbox).FindBackofficeUserById(userId)
	if !ok {
		return none()
	}
	session, ok := findSession(sandbox, user.Id, sessionId)
	if !ok || session.ExpiresAt <= nowSeconds(sandbox) {
		return none()
	}
	return user, session, true
}

// CloseSessions closes every session of the user with id userId but keep, so
// their tokens are refused from here on; nil keep closes them all. It runs
// when the user's password changes, so a session opened with the old one
// does not outlive it.
func CloseSessions(sandbox *api.Sandbox, userId int64, keep *backoffice_db.BackofficeUserSessionRecord) error {
	db := backoffice_db.New(sandbox)
	sessions, err := db.ListBackofficeUserSessions(userId)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		if keep != nil && session.Id == keep.Id {
			continue
		}
		err = backoffice_db.RemoveBackofficeUserSession(sandbox, db, userId, session.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Logout closes session of user: its record is deleted, so its token is
// refused from here on, along with every other session of user that expired.
func Logout(sandbox *api.Sandbox, user backoffice_db.BackofficeUserRecord, session backoffice_db.BackofficeUserSessionRecord) error {
	err := backoffice_db.RemoveBackofficeUserSession(sandbox, backoffice_db.New(sandbox), user.Id, session.Id)
	if err != nil {
		return err
	}
	return dropExpired(sandbox, user.Id, nowSeconds(sandbox))
}
