package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backofficedb"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficetokens"
)

// instantLayout is how the token pages spell an instant.
const instantLayout = "2006-01-02 15:04 UTC"

// BackofficeApiTokensPage is what backoffice/backoffice_api_tokens.html is
// rendered with.
type BackofficeApiTokensPage struct {
	Viewer Viewer
	// Notice is shown above the list, its Text "" for none.
	Notice Notice
	// Created is the token just created, shown once above the list; its
	// Token is "" for none.
	Created CreatedToken
	Tokens  []BackofficeApiTokenRow
	// ShowOwner adds the owner column, for a root, who sees every user's
	// tokens.
	ShowOwner bool
}

// CreatedToken is a token the list page shows in full, right after it was
// created: the one time it is ever shown.
type CreatedToken struct {
	Name  string
	Token string
}

// BackofficeApiTokenRow is one token of the list.
type BackofficeApiTokenRow struct {
	Id     string
	Name   string
	Prefix string
	Owner  string
	// IsOwn marks a token of the signed-in user, on a root's list.
	IsOwn   bool
	Created string
	// Expires is the instant the token expires, or "Never".
	Expires   string
	IsExpired bool
	// Ips are the ips the token is accepted from, or "Any".
	Ips string
	// LastUsed is the instant and the ip of the token's last request, or
	// "Never".
	LastUsed string
}

// apiTokenNoticeOf is how the list page words a backofficetokens notice code;
// an unknown code shows nothing.
func apiTokenNoticeOf(sandbox *api.Sandbox, code string) Notice {
	switch code {
	case backofficetokens.NoticeRevoked:
		return Notice{Text: "Token revoked. It no longer authenticates any request.", Kind: "ok"}
	case backofficetokens.NoticeNotFound:
		return Notice{Text: "That token no longer exists.", Kind: "error"}
	}
	return Notice{}
}

// instantOf is the instant seconds as the token pages spell it.
func instantOf(sandbox *api.Sandbox, seconds int64) string {
	return sandbox.Deps.Timedeps.FormatUnix(seconds, instantLayout)
}

// apiTokenRow is the row of listed on the list page of user, at now.
func apiTokenRow(sandbox *api.Sandbox, user *backofficedb.BackofficeuserItem, listed backofficetokens.Listed, now int64) BackofficeApiTokenRow {
	token := listed.Token
	row := BackofficeApiTokenRow{
		Id:        sandbox.Deps.Stringsdeps.FormatInt(token.Id, 10),
		Name:      token.Name,
		Prefix:    token.Prefix,
		Owner:     listed.Owner.Username,
		IsOwn:     listed.Owner.Id == user.Id,
		Created:   instantOf(sandbox, token.Createdat),
		Expires:   "Never",
		IsExpired: backofficetokens.Expired(sandbox, token, now),
		Ips:       "Any",
		LastUsed:  "Never",
	}
	if token.Expiresat != 0 {
		row.Expires = instantOf(sandbox, token.Expiresat)
	}
	ips := backofficetokens.IpList(sandbox, token)
	if len(ips) > 0 {
		row.Ips = sandbox.Deps.Stringsdeps.Join(ips, ", ")
	}
	if token.Lastusedat != 0 {
		row.LastUsed = instantOf(sandbox, token.Lastusedat) + " from " + token.Lastusedip
	}
	return row
}

// BackofficeApiTokens answers, under status, the API token list page for
// user: listed, with the notice code notice above it, and created shown in
// full when its Token is not "".
func BackofficeApiTokens(sandbox *api.Sandbox, response *serverdeps.Response, status int, user *backofficedb.BackofficeuserItem, listed []backofficetokens.Listed, notice string, created CreatedToken) error {
	now := sandbox.Deps.Std.Now() / 1_000_000_000
	rows := []BackofficeApiTokenRow{}
	for _, item := range listed {
		rows = append(rows, apiTokenRow(sandbox, user, item, now))
	}

	viewer := viewerOf(sandbox, user)
	return Html(sandbox, response, status, "backoffice/backoffice_api_tokens.html", BackofficeApiTokensPage{
		Viewer:    viewer,
		Notice:    apiTokenNoticeOf(sandbox, notice),
		Created:   created,
		Tokens:    rows,
		ShowOwner: viewer.IsRoot,
	})
}
