package backofficeapitokens

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficehttp"
)

// An API token is what the /api/admin routes authenticate with, sent as
// `Authorization: Bearer <token>`. A backoffice user creates and revokes their
// tokens on the backoffice pages; the api only reads them. A token acts as the
// user it belongs to, with the role that user holds at the time of each
// request.
//
// The token itself is TokenPrefix followed by SecretBytes random bytes in hex,
// and is shown once, when it is created: only its SHA-256 is stored, in the
// indexed tokensha field, so a request finds its token in one lookup. Its
// first PrefixLength characters are stored as they are, so the list can tell
// one token from another.

// TokenPrefix starts every API token, so one is told apart from any other
// secret at a glance.
const TokenPrefix = "bo_"

// SecretBytes is how many random bytes a token carries after TokenPrefix.
const SecretBytes = 32

// PrefixLength is how many leading characters of a token the list shows.
const PrefixLength = len(TokenPrefix) + 8

// MaxNameLength is the most characters a token name may have.
const MaxNameLength = 100

// LastUsedResolution is how many seconds a token's last use may lag: its
// last-used-at is written at most once in that time.
const LastUsedResolution = 60

// MaxDaysForViewer is the longest a token of a user who is not a root may
// last: such a user cannot create one that never expires.
const MaxDaysForViewer = 90

// DaySeconds is how many seconds a day of an expiration lasts.
const DaySeconds = 24 * 60 * 60

// DateLayout is how the custom expiration date of the form is spelled, the
// value of an `<input type="date">`.
const DateLayout = "2006-01-02"

// ListPath is the page every revoke sends the browser back to.
const ListPath = "/admin/list-backoffice-api-tokens"

// The expiration values the form offers besides a number of days.
const (
	// ExpirationCustom expires the token at the end of the day the form's
	// date field names.
	ExpirationCustom = "custom"
	// ExpirationNever keeps the token valid until it is revoked.
	ExpirationNever = "never"
	// DefaultExpiration is the expiration the form starts with.
	DefaultExpiration = "30"
)

// The notices a revoke hands the list page through ListPath?notice=<code>.
// Each one is a fixed word the page words itself, so nothing the client sent
// is ever shown back.
const (
	// NoticeRevoked follows a token revoked.
	NoticeRevoked = "revoked"
	// NoticeNotFound follows a revoke of a token that does not exist, or that
	// the user may not revoke.
	NoticeNotFound = "not-found"
)

// Expiration is one choice of the form's expiration field.
type Expiration struct {
	// Value is what the form sends: a number of days, ExpirationCustom or
	// ExpirationNever.
	Value string
	// Label is how the form shows it.
	Label string
	// Days is how long a token chosen with it lasts, 0 for the two that are
	// not a number of days.
	Days int64
}

// Fields is what the create form sends.
type Fields struct {
	Name string
	// Expiration is the Value of one of Expirations.
	Expiration string
	// Date is the last day the token is valid, spelled by DateLayout; read
	// only when Expiration is ExpirationCustom.
	Date string
	// Ips are the client ips the token is accepted from, separated by
	// commas; "" accepts it from any ip.
	Ips string
}

// ApiTokenRow is one token of the list, with the user it belongs to.
type ApiTokenRow struct {
	Token backoffice_db.ApiTokenRecord
	Owner backoffice_db.BackofficeUserRecord
}

// Expirations is every expiration the form offers, in the order it offers them.
func Expirations(sandbox *api.Sandbox) []Expiration {
	return []Expiration{
		{Value: "7", Label: "7 days", Days: 7},
		{Value: "30", Label: "30 days", Days: 30},
		{Value: "60", Label: "60 days", Days: 60},
		{Value: "90", Label: "90 days", Days: 90},
		{Value: "365", Label: "1 year", Days: 365},
		{Value: ExpirationCustom, Label: "Custom date"},
		{Value: ExpirationNever, Label: "No expiration"},
	}
}

// ListLocation is ListPath carrying notice for the list page to show.
func ListLocation(sandbox *api.Sandbox, notice string) string {
	return ListPath + "?notice=" + notice
}

// nowSeconds is the current time in seconds since the Unix epoch.
func nowSeconds(sandbox *api.Sandbox) int64 {
	return sandbox.Deps.StdDeps.Now() / 1_000_000_000
}

// tokenSha is how a token is stored and looked up: its SHA-256, lower-case
// hex. A token is SecretBytes of randomness, so no salt is needed.
func tokenSha(sandbox *api.Sandbox, token string) string {
	return sandbox.Deps.HashDeps.Sha256Hex([]byte(token))
}

// isRoot tells whether user holds the root role.
func isRoot(sandbox *api.Sandbox, user backoffice_db.BackofficeUserRecord) bool {
	return backofficeauth.Role(user.Role) == backofficeauth.RoleRoot
}

// Expired tells whether token stopped being valid by now. A token without an
// expiration never does.
func Expired(sandbox *api.Sandbox, token backoffice_db.ApiTokenRecord, now int64) bool {
	return token.ExpiresAt != 0 && token.ExpiresAt <= now
}

// IpList is the ips a token is accepted from, as stored in its ips field;
// empty when it is accepted from any ip.
func IpList(sandbox *api.Sandbox, token backoffice_db.ApiTokenRecord) []string {
	if token.Ips == "" {
		return []string{}
	}
	return sandbox.Deps.StringsDeps.Split(token.Ips, ",")
}

// issuing holds a token while a token is named and stored, so two requests at
// the same time cannot both give one owner's tokens the same name.
var issuing = make(chan struct{}, 1)

// Add issues a new token for owner, as fields describe it. It answers a
// message for the form when fields are refused; once the token is stored, it
// answers the token itself — the one time it is ever shown — and its record.
func Add(sandbox *api.Sandbox, owner backoffice_db.BackofficeUserRecord, fields Fields) (string, backoffice_db.ApiTokenRecord, string, error) {
	none := backoffice_db.ApiTokenRecord{}
	now := nowSeconds(sandbox)

	issuing <- struct{}{}
	defer func() { <-issuing }()

	name := sandbox.Deps.StringsDeps.TrimSpace(fields.Name)
	message, err := validateName(sandbox, owner, name)
	if err != nil || message != "" {
		return "", none, message, err
	}
	expiresAt, message := expirationOf(sandbox, fields, now)
	if message != "" {
		return "", none, message, nil
	}
	if !isRoot(sandbox, owner) && (expiresAt == 0 || expiresAt > now+MaxDaysForViewer*DaySeconds) {
		return "", none, sandbox.Deps.StdDeps.Sprintf("Only a root user can create a token that lasts more than %d days or never expires.", MaxDaysForViewer), nil
	}
	ips, message, err := normalizeIps(sandbox, fields.Ips)
	if err != nil || message != "" {
		return "", none, message, err
	}

	secret, err := sandbox.Deps.RandDeps.Hex(SecretBytes)
	if err != nil {
		return "", none, "", err
	}
	token := TokenPrefix + secret
	item, err := backoffice_db.New(sandbox).AddApiToken(backoffice_db.ApiTokenInput{
		TokenSha256:  tokenSha(sandbox, token),
		Name:      name,
		Prefix:    token[:PrefixLength],
		OwnerId:   owner.Id,
		CreatedAt: now,
		ExpiresAt: expiresAt,
		Ips:       ips,
	})
	if err != nil {
		return "", none, "", err
	}
	return token, item, "", nil
}

// List is every token viewer may see, newest first: a root sees every user's,
// anyone else only their own. A token whose user no longer exists is left
// out.
func List(sandbox *api.Sandbox, viewer backoffice_db.BackofficeUserRecord) ([]ApiTokenRow, error) {
	db := backoffice_db.New(sandbox)
	tokens, err := db.ListApiTokens(backoffice_db.ApiTokenFilter{})
	if err != nil {
		return nil, err
	}

	listed := []ApiTokenRow{}
	for _, token := range tokens {
		if !isRoot(sandbox, viewer) && token.OwnerId != viewer.Id {
			continue
		}
		owner, ok := db.FindBackofficeUserById(token.OwnerId)
		if !ok {
			continue
		}
		listed = append(listed, ApiTokenRow{Token: token, Owner: owner})
	}
	sandbox.Deps.SortDeps.SliceStable(listed, func(i int, j int) bool {
		return listed[i].Token.Id > listed[j].Token.Id
	})
	return listed, nil
}

// Revoke deletes the token with id id on behalf of actor, so it is refused
// from the next request on. It answers the notice the list page shows next:
// a user revokes their own tokens, a root anyone's, and a token actor may not
// revoke reads as one that does not exist.
func Revoke(sandbox *api.Sandbox, actor backoffice_db.BackofficeUserRecord, id int64) (string, error) {
	db := backoffice_db.New(sandbox)
	token, ok := db.FindApiTokenById(id)
	if !ok || (token.OwnerId != actor.Id && !isRoot(sandbox, actor)) {
		return NoticeNotFound, nil
	}
	err := db.RemoveApiToken(id)
	if err != nil {
		return "", err
	}
	return NoticeRevoked, nil
}

// RemoveOfOwner deletes every token of the user with id ownerId. It runs when
// that user is removed, so no token outlives its user.
func RemoveOfOwner(sandbox *api.Sandbox, ownerId int64) error {
	db := backoffice_db.New(sandbox)
	tokens, err := db.ListApiTokens(backoffice_db.ApiTokenFilter{})
	if err != nil {
		return err
	}
	for _, token := range tokens {
		if token.OwnerId != ownerId {
			continue
		}
		err = db.RemoveApiToken(token.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve answers the user token acts as and its record, for a request that
// came from the client ip ip, or false when token is not an API token, is
// unknown — revoked included —, has expired, is not accepted from ip, or its
// user no longer exists. A token it accepts is marked as last used from ip,
// and now — written only when the ip changed or LastUsedResolution passed
// since the last mark, so a busy token costs no write per request.
func Resolve(sandbox *api.Sandbox, token string, ip string) (backoffice_db.BackofficeUserRecord, backoffice_db.ApiTokenRecord, bool, error) {
	noUser := backoffice_db.BackofficeUserRecord{}
	noToken := backoffice_db.ApiTokenRecord{}
	if !sandbox.Deps.StringsDeps.HasPrefix(token, TokenPrefix) {
		return noUser, noToken, false, nil
	}

	db := backoffice_db.New(sandbox)
	item, ok := db.FindApiTokenByTokenSha256(tokenSha(sandbox, token))
	if !ok {
		return noUser, noToken, false, nil
	}
	now := nowSeconds(sandbox)
	if Expired(sandbox, item, now) || !accepts(sandbox, item, ip) {
		return noUser, noToken, false, nil
	}
	user, ok := db.FindBackofficeUserById(item.OwnerId)
	if !ok {
		return noUser, noToken, false, nil
	}

	if item.LastUsedIp != ip {
		err := db.SetApiTokenLastUsedIp(item.Id, ip)
		if err != nil {
			return noUser, noToken, false, err
		}
		item.LastUsedIp = ip
	}
	if now-item.LastUsedAt >= LastUsedResolution {
		err := db.SetApiTokenLastUsedAt(item.Id, now)
		if err != nil {
			return noUser, noToken, false, err
		}
		item.LastUsedAt = now
	}
	return user, item, true, nil
}

// accepts tells whether token may be used from the client ip ip: any ip when
// it lists none, one it lists — an address, or a range in CIDR notation —
// otherwise. Addresses are compared as parsed, never as spelled.
func accepts(sandbox *api.Sandbox, token backoffice_db.ApiTokenRecord, ip string) bool {
	allowed := IpList(sandbox, token)
	if len(allowed) == 0 {
		return true
	}
	for _, candidate := range allowed {
		if backofficehttp.IpMatches(sandbox, candidate, ip) {
			return true
		}
	}
	return false
}

// validateName answers why name may not name a new token of owner, "" when it
// may: it is required, at most MaxNameLength characters, and no other token of
// owner bears it, regardless of case.
func validateName(sandbox *api.Sandbox, owner backoffice_db.BackofficeUserRecord, name string) (string, error) {
	if name == "" {
		return "The name is required.", nil
	}
	if len([]rune(name)) > MaxNameLength {
		return sandbox.Deps.StdDeps.Sprintf("The name can have at most %d characters.", MaxNameLength), nil
	}

	tokens, err := backoffice_db.New(sandbox).ListApiTokens(backoffice_db.ApiTokenFilter{})
	if err != nil {
		return "", err
	}
	lowered := sandbox.Deps.StringsDeps.ToLower(name)
	for _, token := range tokens {
		if token.OwnerId == owner.Id && sandbox.Deps.StringsDeps.ToLower(token.Name) == lowered {
			return "You already have a token with that name.", nil
		}
	}
	return "", nil
}

// expirationOf is the expiresat a token created now with fields gets, 0 for
// one that never expires, or a message for the form when the expiration is
// refused. A custom date keeps the token valid through the whole of that day,
// in UTC, so it has to be today or later.
func expirationOf(sandbox *api.Sandbox, fields Fields, now int64) (int64, string) {
	switch fields.Expiration {
	case ExpirationNever:
		return 0, ""
	case ExpirationCustom:
		day, err := sandbox.Deps.TimeDeps.ParseUnix(DateLayout, sandbox.Deps.StringsDeps.TrimSpace(fields.Date))
		if err != nil {
			return 0, "Choose the date the token expires on."
		}
		expiresAt := day + DaySeconds
		if expiresAt <= now {
			return 0, "The expiration date has to be in the future."
		}
		return expiresAt, ""
	}
	for _, expiration := range Expirations(sandbox) {
		if expiration.Days > 0 && expiration.Value == fields.Expiration {
			return now + expiration.Days*DaySeconds, ""
		}
	}
	return 0, "Choose a valid expiration."
}

// normalizeIps is the ips field of a token for the list the form sent, or a
// message for the form when one entry is neither an ip nor a CIDR range:
// entries are separated by commas, trimmed, written in their canonical form
// (2001:db8::1 for 2001:0DB8:0::1), stripped of empties and duplicates, and
// joined back by commas. "" is a token accepted from any ip.
func normalizeIps(sandbox *api.Sandbox, list string) (string, string, error) {
	strings := sandbox.Deps.StringsDeps
	kept := []string{}
	seen := map[string]bool{}
	for _, entry := range strings.Split(list, ",") {
		ip := strings.ToLower(strings.TrimSpace(entry))
		if ip == "" || seen[ip] {
			continue
		}
		canonical, valid := backofficehttp.CanonicalIpOrRange(sandbox, ip)
		if !valid {
			return "", sandbox.Deps.StdDeps.Sprintf("%s is not a valid IP address or CIDR range.", ip), nil
		}
		if seen[canonical] {
			continue
		}
		ip = canonical
		seen[ip] = true
		kept = append(kept, ip)
	}
	return strings.Join(kept, ","), "", nil
}
