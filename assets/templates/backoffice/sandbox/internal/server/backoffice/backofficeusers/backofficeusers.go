package backofficeusers

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeapitokens"
)

// ListPath is the page every add, edit and remove sends the browser back to.
const ListPath = "/admin/list-backoffice-users"

// DefaultLimit is how many users a page of the list shows when none is asked.
const DefaultLimit = 20

// MaxLimit is the most users one page of the list shows.
const MaxLimit = 100

// MinPasswordLength is the fewest characters a new password may have.
const MinPasswordLength = 12

// MaxPasswordLength is the most bytes a new password may have: the most a
// sign-in checks (backofficethrottle.MaxPasswordLength).
const MaxPasswordLength = 1024

// MinDistinctCharacters is the fewest different characters a new password
// may hold, so a run of one key ("aaaaaaaaaaaa") or two is refused.
const MinDistinctCharacters = 5

// commonPasswords are passwords of MinPasswordLength or more that leaked
// lists show among the most used, lower-cased: a new password equal to one,
// regardless of case, is refused.
var commonPasswords = []string{
	"123456789012", "1234567890123", "12345678901234", "123456789123",
	"1q2w3e4r5t6y", "1qaz2wsx3edc", "abcdefghijkl", "abcd1234abcd",
	"iloveyou1234", "letmein12345", "password1234", "password12345",
	"password123!", "passwordpassword", "qwertyuiop12", "qwerty123456",
	"qwertyuiop123", "qwertyuiopasdfgh", "q1w2e3r4t5y6", "welcome12345",
	"zaq12wsxcde3", "administrator", "admin1234567", "changeme1234",
}

// GeneratedPasswordBytes is how many random bytes a password GeneratePassword
// makes carries, spelled in hex: twice as many characters.
const GeneratedPasswordBytes = 16

// emailPattern is the shape an email has to have: something, an @, a domain
// with a dot. The mailbox itself is never checked.
const emailPattern = `^[^@\s]+@[^@\s]+\.[^@\s]+$`

// The notices an add, an edit or a remove hands the list page through
// ListPath?notice=<code>. Each one is a fixed word the page words itself, so
// nothing the client sent is ever shown back.
const (
	// NoticeAdded follows a user added.
	NoticeAdded = "added"
	// NoticeUpdated follows a user edited.
	NoticeUpdated = "updated"
	// NoticePasswordChanged follows a user edited with a new password, which
	// ended their sessions and revoked their API tokens.
	NoticePasswordChanged = "password-changed"
	// NoticeRemoved follows a user removed.
	NoticeRemoved = "removed"
	// NoticeNotFound follows an edit or a remove of a user that does not exist.
	NoticeNotFound = "not-found"
	// NoticeSelf follows a root trying to remove its own account.
	NoticeSelf = "self"
	// NoticeNotRoot follows a remove by a user who stopped being a root — or
	// was removed — while the request ran.
	NoticeNotRoot = "not-root"
)

// Query is what the list page is asked for.
type Query struct {
	// Search keeps the users whose username or email holds it, regardless of
	// case; "" keeps every user.
	Search string
	// Role keeps the users of one role, by its backofficeauth.RoleName; ""
	// or a name no role has keeps every role.
	Role string
	// Page is the page to show, counted from 1.
	Page int
	// Limit is how many users a page shows.
	Limit int
	// Viewer is who asks. One who is not a root sees every other user's
	// email masked (Masked), and Search matches only usernames and their own
	// email, so a search cannot spell out an email a letter at a time; nil
	// shows every field.
	Viewer *backoffice_db.BackofficeUserRecord
}

// Listing is one page of the users a Query keeps.
type Listing struct {
	// Users is the page itself.
	Users []backoffice_db.BackofficeUserRecord
	// Search is the Query's Search, trimmed.
	Search string
	// Role is the Query's Role, "" when it names no role.
	Role string
	// Total is how many users the Query keeps, across every page.
	Total int
	// Page is the page shown, within 1 and Pages.
	Page int
	// Pages is how many pages Total fills, at least 1.
	Pages int
	// Limit is how many users a page shows, within 1 and MaxLimit.
	Limit int
}

// writing holds a token while a user is added, edited or removed: the
// checks those make — a username or an email taken, the last root — and the
// write they guard run as one, so two requests at the same time can neither
// both take one name nor leave no root.
var writing = make(chan struct{}, 1)

// lock waits for the writing token; unlock gives it back.
func lock(sandbox *api.Sandbox) {
	writing <- struct{}{}
}

// unlock gives the writing token back.
func unlock(sandbox *api.Sandbox) {
	<-writing
}

// Fields is what the add and edit forms send for one user.
type Fields struct {
	Username string
	Email    string
	// Password is the plain password; "" on an edit keeps the current one.
	Password string
	Role     int64
}

// ListLocation is ListPath carrying notice for the list page to show.
func ListLocation(sandbox *api.Sandbox, notice string) string {
	return ListPath + "?notice=" + notice
}

// Find is the user with id id, or false when there is none.
func Find(sandbox *api.Sandbox, id int64) (backoffice_db.BackofficeUserRecord, bool) {
	return backoffice_db.New(sandbox).FindBackofficeUserById(id)
}

// List answers the page of users query asks for. The filters run here rather
// than through the generated filter: the search matches anywhere in either
// field, and a zero RoleMin/RoleMax turns its filter off, so the root role (0)
// could never be asked for there.
func List(sandbox *api.Sandbox, query Query) (Listing, error) {
	users, err := backoffice_db.New(sandbox).ListBackofficeUsers(backoffice_db.BackofficeUserFilter{})
	if err != nil {
		return Listing{}, err
	}

	listing := Listing{Search: sandbox.Deps.StringsDeps.TrimSpace(query.Search)}
	role, byRole := backofficeauth.ParseRole(sandbox, query.Role)
	if byRole {
		listing.Role = query.Role
	}
	search := sandbox.Deps.StringsDeps.ToLower(listing.Search)

	kept := []backoffice_db.BackofficeUserRecord{}
	for _, user := range users {
		if byRole && backofficeauth.Role(user.Role) != role {
			continue
		}
		user = Masked(sandbox, query.Viewer, user)
		if search != "" && !holds(sandbox, user.Username, search) && !holds(sandbox, user.Email, search) {
			continue
		}
		kept = append(kept, user)
	}

	listing.Total = len(kept)
	listing.Limit = query.Limit
	if listing.Limit < 1 {
		listing.Limit = DefaultLimit
	}
	listing.Limit = min(listing.Limit, MaxLimit)
	listing.Pages = max(1, (listing.Total+listing.Limit-1)/listing.Limit)
	listing.Page = min(max(query.Page, 1), listing.Pages)

	from := (listing.Page - 1) * listing.Limit
	listing.Users = kept[from:min(from+listing.Limit, listing.Total)]
	return listing, nil
}

// Masked is user as viewer may see it: whole for a root, for themselves and
// for a nil viewer; with the email cut to its first character and its domain
// ("a***@example.com") for anyone else, so a viewer reading the list learns
// who has an account without collecting every address.
func Masked(sandbox *api.Sandbox, viewer *backoffice_db.BackofficeUserRecord, user backoffice_db.BackofficeUserRecord) backoffice_db.BackofficeUserRecord {
	if viewer == nil || viewer.Id == user.Id || backofficeauth.Role(viewer.Role) == backofficeauth.RoleRoot {
		return user
	}
	at := sandbox.Deps.StringsDeps.LastIndex(user.Email, "@")
	if at < 1 {
		user.Email = "***"
		return user
	}
	user.Email = user.Email[:1] + "***" + user.Email[at:]
	return user
}

// GeneratePassword is a random password for a new user, of
// GeneratedPasswordBytes random bytes: what add-backoffice-user gives the user
// it creates, so no password ever travels on a command line.
func GeneratePassword(sandbox *api.Sandbox) (string, error) {
	return sandbox.Deps.RandDeps.Hex(GeneratedPasswordBytes)
}

// Add inserts the user fields describe. It answers a message for the form when
// fields are refused, and the user added with "" once it is. Within the
// process the check and the insert run under one lock; another process — an
// add-backoffice-user beside a running server — is caught after the insert:
// of two users holding one username or email, the one inserted later is
// removed again and refused.
func Add(sandbox *api.Sandbox, fields Fields) (backoffice_db.BackofficeUserRecord, string, error) {
	none := backoffice_db.BackofficeUserRecord{}
	fields = trimmed(sandbox, fields)
	message, err := validateFields(sandbox, fields, true)
	if err != nil || message != "" {
		return none, message, err
	}
	hash, err := backofficeauth.HashPassword(sandbox, fields.Password)
	if err != nil {
		return none, "", err
	}

	lock(sandbox)
	defer unlock(sandbox)
	message, err = validateUnique(sandbox, fields, true, 0)
	if err != nil || message != "" {
		return none, message, err
	}
	db := backoffice_db.New(sandbox)
	user, err := db.AddBackofficeUser(backoffice_db.BackofficeUserInput{
		Username:     fields.Username,
		Email:        fields.Email,
		PasswordHash: hash,
		Role:         fields.Role,
	})
	if err != nil {
		return none, "", err
	}
	message, err = validateUnique(sandbox, fields, false, user.Id)
	if err != nil || message == "" {
		return user, "", err
	}
	if !insertedFirst(sandbox, fields, user.Id) {
		if err := db.RemoveBackofficeUser(user.Id); err != nil {
			return none, "", err
		}
		return none, message, nil
	}
	return user, "", nil
}

// insertedFirst tells whether the user with id id holds fields' username and
// email before every other user sharing either: the one with the lower id.
func insertedFirst(sandbox *api.Sandbox, fields Fields, id int64) bool {
	strings := sandbox.Deps.StringsDeps
	users, err := backoffice_db.New(sandbox).ListBackofficeUsers(backoffice_db.BackofficeUserFilter{})
	if err != nil {
		return false
	}
	names := map[string]bool{strings.ToLower(fields.Username): true, strings.ToLower(fields.Email): true}
	for _, user := range users {
		if user.Id < id && (names[strings.ToLower(user.Username)] || names[strings.ToLower(user.Email)]) {
			return false
		}
	}
	return true
}

// Set writes fields over the user with id id on behalf of actor, the
// password only when one was given. It answers a message for the form when
// fields are refused — demoting the last root among them — and, once the user
// is written, the notice the list page shows next. A new role holds from the
// user's next request, because the authentication middlewares read the user
// afresh on each one. A new password ends every session of the user and
// revokes every API token of theirs, so whoever held one opened with the old
// password is out; only session, the actor's own, is spared when actor edits
// their own account.
func Set(sandbox *api.Sandbox, actor backoffice_db.BackofficeUserRecord, session *backoffice_db.BackofficeUserSessionRecord, id int64, fields Fields) (string, string, error) {
	fields = trimmed(sandbox, fields)
	message, err := validateFields(sandbox, fields, false)
	if err != nil || message != "" {
		return message, "", err
	}

	lock(sandbox)
	defer unlock(sandbox)
	if current, ok := Find(sandbox, actor.Id); !ok || backofficeauth.Role(current.Role) != backofficeauth.RoleRoot {
		return "You are no longer a root user, so nothing was changed.", "", nil
	}
	message, err = validateUnique(sandbox, fields, false, id)
	if err != nil || message != "" {
		return message, "", err
	}
	user, ok := Find(sandbox, id)
	if !ok {
		return "That user no longer exists.", "", nil
	}
	if backofficeauth.Role(user.Role) == backofficeauth.RoleRoot && backofficeauth.Role(fields.Role) != backofficeauth.RoleRoot {
		roots, err := countRoots(sandbox)
		if err != nil {
			return "", "", err
		}
		if roots <= 1 {
			return "This is the last root user, so it has to stay root.", "", nil
		}
	}

	db := backoffice_db.New(sandbox)
	err = db.SetBackofficeUserUsername(id, fields.Username)
	if err != nil {
		return "", "", err
	}
	err = db.SetBackofficeUserEmail(id, fields.Email)
	if err != nil {
		return "", "", err
	}
	err = db.SetBackofficeUserRole(id, fields.Role)
	if err != nil {
		return "", "", err
	}
	if fields.Password == "" {
		return "", NoticeUpdated, nil
	}

	hash, err := backofficeauth.HashPassword(sandbox, fields.Password)
	if err != nil {
		return "", "", err
	}
	err = db.SetBackofficeUserPasswordHash(id, hash)
	if err != nil {
		return "", "", err
	}
	var keep *backoffice_db.BackofficeUserSessionRecord
	if actor.Id == id {
		keep = session
	}
	err = backofficeauth.CloseSessions(sandbox, id, keep)
	if err != nil {
		return "", "", err
	}
	err = backofficeapitokens.RemoveOfOwner(sandbox, id)
	if err != nil {
		return "", "", err
	}
	return "", NoticePasswordChanged, nil
}

// Remove deletes the user with id id, and every session and API token of it
// with it, on behalf of actor. It answers the notice the list page shows next:
// a root may not remove its own account, and since actor is checked to still
// be a root under the same lock as the removal, another root always remains.
func Remove(sandbox *api.Sandbox, actor backoffice_db.BackofficeUserRecord, id int64) (string, error) {
	if id == actor.Id {
		return NoticeSelf, nil
	}
	lock(sandbox)
	defer unlock(sandbox)
	current, ok := Find(sandbox, actor.Id)
	if !ok || backofficeauth.Role(current.Role) != backofficeauth.RoleRoot {
		return NoticeNotRoot, nil
	}
	_, ok = Find(sandbox, id)
	if !ok {
		return NoticeNotFound, nil
	}
	err := backofficeapitokens.RemoveOfOwner(sandbox, id)
	if err != nil {
		return "", err
	}
	err = backoffice_db.New(sandbox).RemoveBackofficeUser(id)
	if err != nil {
		return "", err
	}
	return NoticeRemoved, nil
}

// holds tells whether text holds search, which is already lower case,
// regardless of case.
func holds(sandbox *api.Sandbox, text string, search string) bool {
	return sandbox.Deps.StringsDeps.Contains(sandbox.Deps.StringsDeps.ToLower(text), search)
}

// trimmed is fields with the spaces around the username and the email cut.
func trimmed(sandbox *api.Sandbox, fields Fields) Fields {
	fields.Username = sandbox.Deps.StringsDeps.TrimSpace(fields.Username)
	fields.Email = sandbox.Deps.StringsDeps.TrimSpace(fields.Email)
	return fields
}

// validateFields answers why fields may not be stored, "" when they may,
// leaving out whether the username or the email is taken: validateUnique
// checks that, under the lock. A new user (creating) needs a password.
func validateFields(sandbox *api.Sandbox, fields Fields, creating bool) (string, error) {
	strings := sandbox.Deps.StringsDeps
	if fields.Username == "" {
		return "The username is required.", nil
	}
	if strings.ContainsAny(fields.Username, " \t\r\n") {
		return "The username cannot contain spaces.", nil
	}
	email, err := strings.MatchPattern(emailPattern, fields.Email)
	if err != nil {
		return "", err
	}
	if !email {
		return "Enter a valid email address.", nil
	}
	if !backofficeauth.ValidRole(sandbox, backofficeauth.Role(fields.Role)) {
		return "Choose a valid role.", nil
	}
	if creating && fields.Password == "" {
		return "The password is required.", nil
	}
	if fields.Password != "" {
		return passwordProblem(sandbox, fields), nil
	}
	return "", nil
}

// passwordProblem answers why the password of fields is too weak to be
// stored, "" when it is not: shorter than MinPasswordLength characters,
// longer than MaxPasswordLength bytes, holding fewer than
// MinDistinctCharacters different characters, one of commonPasswords, or the
// user's own username or email.
func passwordProblem(sandbox *api.Sandbox, fields Fields) string {
	strings := sandbox.Deps.StringsDeps
	password := fields.Password
	if len([]rune(password)) < MinPasswordLength {
		return sandbox.Deps.StdDeps.Sprintf("The password needs at least %d characters.", MinPasswordLength)
	}
	if len(password) > MaxPasswordLength {
		return sandbox.Deps.StdDeps.Sprintf("The password can have at most %d characters.", MaxPasswordLength)
	}
	distinct := map[rune]bool{}
	for _, char := range password {
		distinct[char] = true
	}
	if len(distinct) < MinDistinctCharacters {
		return sandbox.Deps.StdDeps.Sprintf("The password needs at least %d different characters.", MinDistinctCharacters)
	}
	lowered := strings.ToLower(password)
	for _, common := range commonPasswords {
		if lowered == common {
			return "That password is one of the most used ones. Choose another."
		}
	}
	if lowered == strings.ToLower(fields.Username) || lowered == strings.ToLower(fields.Email) {
		return "The password cannot be the username or the email."
	}
	return ""
}

// validateUnique answers why the username or the email of fields is taken,
// "" when neither is. A new user (creating) is checked against every user; an
// edited one, with id id, keeps its own username and email without them
// counting as taken. Login reads its field as a username or an email, so
// neither may be any other user's username or email, regardless of case.
func validateUnique(sandbox *api.Sandbox, fields Fields, creating bool, id int64) (string, error) {
	strings := sandbox.Deps.StringsDeps
	users, err := backoffice_db.New(sandbox).ListBackofficeUsers(backoffice_db.BackofficeUserFilter{})
	if err != nil {
		return "", err
	}
	username := strings.ToLower(fields.Username)
	address := strings.ToLower(fields.Email)
	for _, user := range users {
		if !creating && user.Id == id {
			continue
		}
		taken := []string{strings.ToLower(user.Username), strings.ToLower(user.Email)}
		for _, login := range taken {
			if login == username {
				return "That username is already in use.", nil
			}
			if login == address {
				return "That email is already in use.", nil
			}
		}
	}
	return "", nil
}

// countRoots is how many users hold the root role.
func countRoots(sandbox *api.Sandbox) (int, error) {
	users, err := backoffice_db.New(sandbox).ListBackofficeUsers(backoffice_db.BackofficeUserFilter{})
	if err != nil {
		return 0, err
	}
	roots := 0
	for _, user := range users {
		if backofficeauth.Role(user.Role) == backofficeauth.RoleRoot {
			roots++
		}
	}
	return roots, nil
}
