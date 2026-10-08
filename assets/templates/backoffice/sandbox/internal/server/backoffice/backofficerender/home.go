package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
)

// HomePage is what backoffice/home.html is rendered with.
type HomePage struct {
	Id       string
	Username string
	// Initial is the first letter of Username, shown in the avatar.
	Initial        string
	Email          string
	Role           string
	SessionMinutes int
	// IsRoot shows what only a root may do.
	IsRoot bool
}

// initialOf is the first letter of name, "?" for an empty one.
func initialOf(sandbox *api.Sandbox, name string) string {
	for _, letter := range name {
		return string(letter)
	}
	return "?"
}

// RenderHomePage answers the home page for user, whose session lasts sessionMinutes.
func RenderHomePage(sandbox *api.Sandbox, response *serverdeps.Response, user *backoffice_db.BackofficeUserRecord, sessionMinutes int) error {
	return RenderHTML(sandbox, response, api.StatusOK, "backoffice/home.html", HomePage{
		Id:             sandbox.Deps.StringsDeps.FormatInt(user.Id, 10),
		Username:       user.Username,
		Initial:        initialOf(sandbox, user.Username),
		Email:          user.Email,
		Role:           backofficeauth.RoleName(sandbox, backofficeauth.Role(user.Role)),
		SessionMinutes: sessionMinutes,
		IsRoot:         backofficeauth.Role(user.Role) == backofficeauth.RoleRoot,
	})
}
