package backofficetokens

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backofficedb"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeguard"
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

// Listed is one token of the list, with the user it belongs to.
type Listed struct {
	Token backofficedb.ApitokenItem
	Owner backofficedb.BackofficeuserItem
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
	return sandbox.Deps.Std.Now() / 1_000_000_000
}

// tokenSha is how a token is stored and looked up: its SHA-256, lower-case
// hex. A token is SecretBytes of randomness, so no salt is needed.
func tokenSha(sandbox *api.Sandbox, token string) string {
	return sandbox.Deps.Hashdeps.Sha256Hex([]byte(token))
}

// isRoot tells whether user holds the root role.
func isRoot(sandbox *api.Sandbox, user backofficedb.BackofficeuserItem) bool {
	return backofficeauth.Role(user.Role) == backofficeauth.RoleRoot
}

// Expired tells whether token stopped being valid by now. A token without an
// expiration never does.
func Expired(sandbox *api.Sandbox, token backofficedb.ApitokenItem, now int64) bool {
	return token.Expiresat != 0 && token.Expiresat <= now
}

// IpList is the ips a token is accepted from, as stored in its ips field;
// empty when it is accepted from any ip.
func IpList(sandbox *api.Sandbox, token backofficedb.ApitokenItem) []string {
	if token.Ips == "" {
		return []string{}
	}
	return sandbox.Deps.Stringsdeps.Split(token.Ips, ",")
}

// Create issues a new token for owner, as fields describe it. It answers a
// message for the form when fields are refused; once the token is stored, it
// answers the token itself — the one time it is ever shown — and its record.
func Create(sandbox *api.Sandbox, owner backofficedb.BackofficeuserItem, fields Fields) (string, backofficedb.ApitokenItem, string, error) {
	none := backofficedb.ApitokenItem{}
	now := nowSeconds(sandbox)

	name := sandbox.Deps.Stringsdeps.TrimSpace(fields.Name)
	message, err := validateName(sandbox, owner, name)
	if err != nil || message != "" {
		return "", none, message, err
	}
	expiresAt, message := expirationOf(sandbox, fields, now)
	if message != "" {
		return "", none, message, nil
	}
	ips, message, err := normalizeIps(sandbox, fields.Ips)
	if err != nil || message != "" {
		return "", none, message, err
	}

	secret, err := sandbox.Deps.Randdeps.Hex(SecretBytes)
	if err != nil {
		return "", none, "", err
	}
	token := TokenPrefix + secret
	item, err := backofficedb.New(sandbox).AddApitoken(backofficedb.ApitokenNew{
		Tokensha:  tokenSha(sandbox, token),
		Name:      name,
		Prefix:    token[:PrefixLength],
		Ownerid:   owner.Id,
		Createdat: now,
		Expiresat: expiresAt,
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
func List(sandbox *api.Sandbox, viewer backofficedb.BackofficeuserItem) ([]Listed, error) {
	db := backofficedb.New(sandbox)
	tokens, err := db.ListApitoken(backofficedb.ApitokenFiltrage{})
	if err != nil {
		return nil, err
	}

	listed := []Listed{}
	for _, token := range tokens {
		if !isRoot(sandbox, viewer) && token.Ownerid != viewer.Id {
			continue
		}
		owner, ok := db.FindBackofficeuserById(token.Ownerid)
		if !ok {
			continue
		}
		listed = append(listed, Listed{Token: token, Owner: owner})
	}
	sandbox.Deps.Sortdeps.SliceStable(listed, func(i int, j int) bool {
		return listed[i].Token.Id > listed[j].Token.Id
	})
	return listed, nil
}

// Revoke deletes the token with id id on behalf of actor, so it is refused
// from the next request on. It answers the notice the list page shows next:
// a user revokes their own tokens, a root anyone's, and a token actor may not
// revoke reads as one that does not exist.
func Revoke(sandbox *api.Sandbox, actor backofficedb.BackofficeuserItem, id int64) (string, error) {
	db := backofficedb.New(sandbox)
	token, ok := db.FindApitokenById(id)
	if !ok || (token.Ownerid != actor.Id && !isRoot(sandbox, actor)) {
		return NoticeNotFound, nil
	}
	err := db.RemoveApitoken(id)
	if err != nil {
		return "", err
	}
	return NoticeRevoked, nil
}

// RemoveOfOwner deletes every token of the user with id ownerId. It runs when
// that user is removed, so no token outlives its user.
func RemoveOfOwner(sandbox *api.Sandbox, ownerId int64) error {
	db := backofficedb.New(sandbox)
	tokens, err := db.ListApitoken(backofficedb.ApitokenFiltrage{})
	if err != nil {
		return err
	}
	for _, token := range tokens {
		if token.Ownerid != ownerId {
			continue
		}
		err = db.RemoveApitoken(token.Id)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve answers the user token acts as and its record, for a request that
// came from the client ip ip, or false when token is not an API token, is
// unknown — revoked included —, has expired, is not accepted from ip, or its
// user no longer exists. A token it accepts is marked as last used now, from
// ip.
func Resolve(sandbox *api.Sandbox, token string, ip string) (backofficedb.BackofficeuserItem, backofficedb.ApitokenItem, bool, error) {
	noUser := backofficedb.BackofficeuserItem{}
	noToken := backofficedb.ApitokenItem{}
	if !sandbox.Deps.Stringsdeps.HasPrefix(token, TokenPrefix) {
		return noUser, noToken, false, nil
	}

	db := backofficedb.New(sandbox)
	item, ok := db.FindApitokenByTokensha(tokenSha(sandbox, token))
	if !ok {
		return noUser, noToken, false, nil
	}
	now := nowSeconds(sandbox)
	if Expired(sandbox, item, now) || !accepts(sandbox, item, ip) {
		return noUser, noToken, false, nil
	}
	user, ok := db.FindBackofficeuserById(item.Ownerid)
	if !ok {
		return noUser, noToken, false, nil
	}

	err := db.UpdateApitokenLastusedat(item.Id, now)
	if err != nil {
		return noUser, noToken, false, err
	}
	err = db.UpdateApitokenLastusedip(item.Id, ip)
	if err != nil {
		return noUser, noToken, false, err
	}
	item.Lastusedat = now
	item.Lastusedip = ip
	return user, item, true, nil
}

// accepts tells whether token may be used from the client ip ip: any ip when
// it lists none, one of its own otherwise.
func accepts(sandbox *api.Sandbox, token backofficedb.ApitokenItem, ip string) bool {
	allowed := IpList(sandbox, token)
	if len(allowed) == 0 {
		return true
	}
	ip = sandbox.Deps.Stringsdeps.ToLower(ip)
	for _, candidate := range allowed {
		if candidate == ip {
			return true
		}
	}
	return false
}

// validateName answers why name may not name a new token of owner, "" when it
// may: it is required, at most MaxNameLength characters, and no other token of
// owner bears it, regardless of case.
func validateName(sandbox *api.Sandbox, owner backofficedb.BackofficeuserItem, name string) (string, error) {
	if name == "" {
		return "The name is required.", nil
	}
	if len([]rune(name)) > MaxNameLength {
		return sandbox.Deps.Std.Sprintf("The name can have at most %d characters.", MaxNameLength), nil
	}

	tokens, err := backofficedb.New(sandbox).ListApitoken(backofficedb.ApitokenFiltrage{})
	if err != nil {
		return "", err
	}
	lowered := sandbox.Deps.Stringsdeps.ToLower(name)
	for _, token := range tokens {
		if token.Ownerid == owner.Id && sandbox.Deps.Stringsdeps.ToLower(token.Name) == lowered {
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
		day, err := sandbox.Deps.Timedeps.ParseUnix(DateLayout, sandbox.Deps.Stringsdeps.TrimSpace(fields.Date))
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
// message for the form when one entry is not an ip: entries are separated by
// commas, trimmed, lower-cased, stripped of empties and duplicates, and joined
// back by commas. "" is a token accepted from any ip.
func normalizeIps(sandbox *api.Sandbox, list string) (string, string, error) {
	strings := sandbox.Deps.Stringsdeps
	kept := []string{}
	seen := map[string]bool{}
	for _, entry := range strings.Split(list, ",") {
		ip := strings.ToLower(strings.TrimSpace(entry))
		if ip == "" || seen[ip] {
			continue
		}
		valid, err := backofficeguard.IsIp(sandbox, ip)
		if err != nil {
			return "", "", err
		}
		if !valid {
			return "", sandbox.Deps.Std.Sprintf("%s is not a valid IP address.", ip), nil
		}
		seen[ip] = true
		kept = append(kept, ip)
	}
	return strings.Join(kept, ","), "", nil
}
