package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// LoginPage is what backoffice/login.html is rendered with.
type LoginPage struct {
	// Error is shown above the form, "" for none.
	Error string
	// Username refills the login field after a failed attempt.
	Username string
}

// Login answers the login page under status, with message above the form.
func Login(sandbox *api.Sandbox, response *serverdeps.Response, status int, message string, username string) error {
	return Html(sandbox, response, status, "backoffice/login.html", LoginPage{Error: message, Username: username})
}
