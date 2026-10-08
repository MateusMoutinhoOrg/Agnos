package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeusers"
)

// BackofficeUserFormPage is what backoffice/backoffice_user_form.html is
// rendered with: the form that adds a backoffice user, or the one that edits
// one.
type BackofficeUserFormPage struct {
	Viewer Viewer
	// Title heads the page.
	Title string
	// Action is where the form posts.
	Action string
	// Submit labels the submit button.
	Submit string
	// Error is shown above the form, "" for none.
	Error    string
	Username string
	Email    string
	Roles    []RoleOption
	// PasswordRequired is true when adding: an edit keeps the current
	// password when the field is left blank.
	PasswordRequired bool
}

// RoleOption is one role the form offers.
type RoleOption struct {
	Value    int64
	Name     string
	Selected bool
}

// roleOptions are every role, the one of fields selected.
func roleOptions(sandbox *api.Sandbox, fields backofficeusers.Fields) []RoleOption {
	options := []RoleOption{}
	for _, role := range backofficeauth.Roles(sandbox) {
		options = append(options, RoleOption{
			Value:    int64(role),
			Name:     backofficeauth.RoleName(sandbox, role),
			Selected: int64(role) == fields.Role,
		})
	}
	return options
}

// RenderAddUserPage answers, under status, the form that adds a
// backoffice user, filled with fields — never their password — and with
// message above it.
func RenderAddUserPage(sandbox *api.Sandbox, response *serverdeps.Response, status int, user *backoffice_db.BackofficeUserRecord, fields backofficeusers.Fields, message string) error {
	return RenderHTML(sandbox, response, status, "backoffice/backoffice_user_form.html", BackofficeUserFormPage{
		Viewer:           viewerOf(sandbox, user),
		Title:            "Add user",
		Action:           "/admin/root/add-backoffice-user",
		Submit:           "Add user",
		Error:            message,
		Username:         fields.Username,
		Email:            fields.Email,
		Roles:            roleOptions(sandbox, fields),
		PasswordRequired: true,
	})
}

// RenderSetUserPage answers, under status, the form that edits the
// backoffice user with id id, filled with fields — never their password — and
// with message above it.
func RenderSetUserPage(sandbox *api.Sandbox, response *serverdeps.Response, status int, user *backoffice_db.BackofficeUserRecord, id int64, fields backofficeusers.Fields, message string) error {
	return RenderHTML(sandbox, response, status, "backoffice/backoffice_user_form.html", BackofficeUserFormPage{
		Viewer:   viewerOf(sandbox, user),
		Title:    "Edit user #" + sandbox.Deps.StringsDeps.FormatInt(id, 10),
		Action:   "/admin/root/set-backoffice-user/" + sandbox.Deps.StringsDeps.FormatInt(id, 10),
		Submit:   "Save changes",
		Error:    message,
		Username: fields.Username,
		Email:    fields.Email,
		Roles:    roleOptions(sandbox, fields),
	})
}
