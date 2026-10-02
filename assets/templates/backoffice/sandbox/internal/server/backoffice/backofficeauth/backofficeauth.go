package backofficeauth

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/jwtdeps"
	"{{.Module}}/sandbox/internal/databases/backofficedb"
)

// CookieName is the cookie the session token travels in. The authentication
// route declares a cookie parameter under the same key.
const CookieName = "admin_token"

// SessionSeconds is how long a session token is valid after login.
const SessionSeconds = 30 * 60

// SecretEnv is the environment variable start-server reads the secret that
// signs session tokens from: the project's name upper-cased, every character
// but a letter or a digit turned into "_", then "_SECRET" — MEUSITE_SECRET for
// a project named meusite. It follows the name the project is built under, so
// renaming the project renames the variable. It is never a flag — every user
// of the machine reads a command line — nor a file, which can end up committed
// with the code.
func SecretEnv(sandbox *api.Sandbox) string {
	name := []byte(sandbox.Deps.Stringsdeps.ToUpper(sandbox.Config.ProjectName))
	for i, char := range name {
		if !(char >= 'A' && char <= 'Z') && !(char >= '0' && char <= '9') {
			name[i] = '_'
		}
	}
	return string(name) + "_SECRET"
}

// MinSecretLength is the fewest characters the secret may have.
const MinSecretLength = 32

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
	return sandbox.Deps.Std.Now() / 1_000_000_000
}

// ReadSecret is the secret that signs session tokens, read from the SecretEnv
// environment variable. It refuses one that is unset or shorter than
// MinSecretLength, saying how to set it.
func ReadSecret(sandbox *api.Sandbox) (string, error) {
	env := SecretEnv(sandbox)
	secret := sandbox.Deps.Envdeps.Getenv(env)
	if len(secret) < MinSecretLength {
		return "", sandbox.Deps.Std.Errorf("set the %s environment variable to a random secret of at least %d characters (openssl rand -hex 32); it signs the backoffice sessions", env, MinSecretLength)
	}
	return secret, nil
}

// HashPassword is how a backoffice password is stored: a salted, deliberately
// slow hash of it, with a salt of its own, so equal passwords never share a
// hash and a leaked database is slow to guess at.
func HashPassword(sandbox *api.Sandbox, password string) (string, error) {
	return sandbox.Deps.Passworddeps.Hash(password)
}

// FindByLogin looks a backoffice user up by username or email. It reads every
// user once whichever the login is, so a username, an email and a login that
// names nobody take the same time.
func FindByLogin(sandbox *api.Sandbox, login string) (backofficedb.BackofficeuserItem, bool, error) {
	users, err := backofficedb.New(sandbox).ListBackofficeuser(backofficedb.BackofficeuserFiltrage{})
	if err != nil {
		return backofficedb.BackofficeuserItem{}, false, err
	}
	for _, user := range users {
		if user.Username == login || user.Email == login {
			return user, true, nil
		}
	}
	return backofficedb.BackofficeuserItem{}, false, nil
}

// Authenticate answers the user whose login (username or email) and password
// match, or false when either does not. A login that names nobody still costs
// one password hash, so how long a refusal takes never tells an unknown login
// from a wrong password.
func Authenticate(sandbox *api.Sandbox, login string, password string) (backofficedb.BackofficeuserItem, bool, error) {
	none := backofficedb.BackofficeuserItem{}
	user, ok, err := FindByLogin(sandbox, login)
	if err != nil {
		return none, false, err
	}
	if !ok {
		_, err = HashPassword(sandbox, password)
		return none, false, err
	}
	match, err := sandbox.Deps.Passworddeps.Verify(user.Passwordhash, password)
	if err != nil || !match {
		return none, false, err
	}
	return user, true, nil
}

// IssueToken opens a session for user, whose request came from the client ip
// ip, and signs its token, valid for SessionSeconds and bound to ip. The session is a sessions record under user, living as long
// as the token; its id travels as the token's `jti`, so the authentication
// middleware can tell whether that one session is still open. The user's
// expired sessions are dropped first, so they never pile up.
func IssueToken(sandbox *api.Sandbox, user backofficedb.BackofficeuserItem, ip string) (string, error) {
	now := nowSeconds(sandbox)
	err := dropExpired(sandbox, user.Id, now)
	if err != nil {
		return "", err
	}

	session, err := backofficedb.New(sandbox).AddBackofficeuserSessions(user.Id, backofficedb.SessionsNew{
		Expiresat: now + SessionSeconds,
	})
	if err != nil {
		return "", err
	}

	return sandbox.Deps.Jwtdeps.Sign(jwtdeps.Claims{
		Id:        sandbox.Deps.Stringsdeps.FormatInt(session.Id, 10),
		Subject:   sandbox.Deps.Stringsdeps.FormatInt(user.Id, 10),
		IssuedAt:  now,
		ExpiresAt: now + SessionSeconds,
		Ip:        ip,
	}, sandbox.Config.Secret)
}

// SessionCookie is the Set-Cookie value carrying token: HttpOnly,
// SameSite=Strict, Secure unless start-server serves plain http, on every
// path, expiring with the token.
func SessionCookie(sandbox *api.Sandbox, token string) string {
	return sandbox.Deps.Std.Sprintf("%s=%s; Path=/; Max-Age=%d; HttpOnly; SameSite=Strict%s", CookieName, token, SessionSeconds, secureAttribute(sandbox))
}

// ClearedCookie is the Set-Cookie value that removes the session cookie.
func ClearedCookie(sandbox *api.Sandbox) string {
	return sandbox.Deps.Std.Sprintf("%s=; Path=/; Max-Age=0; HttpOnly; SameSite=Strict%s", CookieName, secureAttribute(sandbox))
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
	fields := sandbox.Deps.Stringsdeps.Fields(authorization)
	if len(fields) != 2 || sandbox.Deps.Stringsdeps.ToLower(fields[0]) != "bearer" {
		return ""
	}
	return fields[1]
}

// dropExpired deletes every session of the user with id userId that expired
// by now: its token is refused anyway, so the record is only garbage.
func dropExpired(sandbox *api.Sandbox, userId int64, now int64) error {
	db := backofficedb.New(sandbox)
	sessions, err := db.ListBackofficeuserSessions(userId)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		if session.Expiresat <= now {
			err = backofficedb.RemoveBackofficeuserSessions(sandbox, db, userId, session.Id)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// findSession is the session with id sessionId of the user with id userId.
func findSession(sandbox *api.Sandbox, userId int64, sessionId int64) (backofficedb.SessionsItem, bool) {
	sessions, err := backofficedb.New(sandbox).ListBackofficeuserSessions(userId)
	if err != nil {
		return backofficedb.SessionsItem{}, false
	}
	for _, session := range sessions {
		if session.Id == sessionId {
			return session, true
		}
	}
	return backofficedb.SessionsItem{}, false
}

// SessionOfToken answers the user a session token was issued for and the
// session it names, or false when the token is invalid or expired, was issued
// to another client ip than ip, names a session that was closed by a logout,
// or its user no longer exists.
func SessionOfToken(sandbox *api.Sandbox, token string, ip string) (backofficedb.BackofficeuserItem, backofficedb.SessionsItem, bool) {
	none := func() (backofficedb.BackofficeuserItem, backofficedb.SessionsItem, bool) {
		return backofficedb.BackofficeuserItem{}, backofficedb.SessionsItem{}, false
	}
	if token == "" || ip == "" {
		return none()
	}
	claims, err := sandbox.Deps.Jwtdeps.Parse(token, sandbox.Config.Secret)
	if err != nil || claims.Ip != ip {
		return none()
	}
	userId, err := sandbox.Deps.Stringsdeps.ParseInt(claims.Subject, 10, 64)
	if err != nil {
		return none()
	}
	sessionId, err := sandbox.Deps.Stringsdeps.ParseInt(claims.Id, 10, 64)
	if err != nil {
		return none()
	}
	user, ok := backofficedb.New(sandbox).FindBackofficeuserById(userId)
	if !ok {
		return none()
	}
	session, ok := findSession(sandbox, user.Id, sessionId)
	if !ok || session.Expiresat <= nowSeconds(sandbox) {
		return none()
	}
	return user, session, true
}

// CloseSessions closes every session of the user with id userId but keep, so
// their tokens are refused from here on; nil keep closes them all. It runs
// when the user's password changes, so a session opened with the old one
// does not outlive it.
func CloseSessions(sandbox *api.Sandbox, userId int64, keep *backofficedb.SessionsItem) error {
	db := backofficedb.New(sandbox)
	sessions, err := db.ListBackofficeuserSessions(userId)
	if err != nil {
		return err
	}
	for _, session := range sessions {
		if keep != nil && session.Id == keep.Id {
			continue
		}
		err = backofficedb.RemoveBackofficeuserSessions(sandbox, db, userId, session.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Logout closes session of user: its record is deleted, so its token is
// refused from here on, along with every other session of user that expired.
func Logout(sandbox *api.Sandbox, user backofficedb.BackofficeuserItem, session backofficedb.SessionsItem) error {
	err := backofficedb.RemoveBackofficeuserSessions(sandbox, backofficedb.New(sandbox), user.Id, session.Id)
	if err != nil {
		return err
	}
	return dropExpired(sandbox, user.Id, nowSeconds(sandbox))
}
