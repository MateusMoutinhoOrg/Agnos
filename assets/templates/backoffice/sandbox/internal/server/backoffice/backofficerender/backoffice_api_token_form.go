package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backofficedb"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficetokens"
)

// BackofficeApiTokenFormPage is what backoffice/backoffice_api_token_form.html
// is rendered with: the form that creates an API token.
type BackofficeApiTokenFormPage struct {
	Viewer Viewer
	// Error is shown above the form, "" for none.
	Error       string
	Name        string
	Expirations []ExpirationOption
	// IsCustom shows the date field, for the custom expiration.
	IsCustom bool
	Date     string
	// MinDate is the first day the date field may name: today, in UTC.
	MinDate string
	Ips     string
	// ClientIp is the ip the browser's request came from, offered as one
	// click to fill the ips field.
	ClientIp string
}

// ExpirationOption is one expiration the form offers.
type ExpirationOption struct {
	Value    string
	Label    string
	Selected bool
}

// expirationOptions are every expiration, the one of fields selected.
func expirationOptions(sandbox *api.Sandbox, fields backofficetokens.Fields) []ExpirationOption {
	options := []ExpirationOption{}
	for _, expiration := range backofficetokens.Expirations(sandbox) {
		options = append(options, ExpirationOption{
			Value:    expiration.Value,
			Label:    expiration.Label,
			Selected: expiration.Value == fields.Expiration,
		})
	}
	return options
}

// CreateBackofficeApiTokenForm answers, under status, the form that creates
// an API token for user, filled with fields and with message above it, the
// client ip clientIp offered for the ips field.
func CreateBackofficeApiTokenForm(sandbox *api.Sandbox, response *serverdeps.Response, status int, user *backofficedb.BackofficeuserItem, fields backofficetokens.Fields, clientIp string, message string) error {
	now := sandbox.Deps.Std.Now() / 1_000_000_000
	return Html(sandbox, response, status, "backoffice/backoffice_api_token_form.html", BackofficeApiTokenFormPage{
		Viewer:      viewerOf(sandbox, user),
		Error:       message,
		Name:        fields.Name,
		Expirations: expirationOptions(sandbox, fields),
		IsCustom:    fields.Expiration == backofficetokens.ExpirationCustom,
		Date:        fields.Date,
		MinDate:     sandbox.Deps.Timedeps.FormatUnix(now, backofficetokens.DateLayout),
		Ips:         fields.Ips,
		ClientIp:    clientIp,
	})
}
