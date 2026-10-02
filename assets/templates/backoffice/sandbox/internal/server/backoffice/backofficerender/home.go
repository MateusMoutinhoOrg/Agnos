package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backofficedb"
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

// Home answers the home page for user, whose session lasts sessionMinutes.
func Home(sandbox *api.Sandbox, response *serverdeps.Response, user *backofficedb.BackofficeuserItem, sessionMinutes int) error {
	return Html(sandbox, response, api.StatusOk, "backoffice/home.html", HomePage{
		Id:             sandbox.Deps.Stringsdeps.FormatInt(user.Id, 10),
		Username:       user.Username,
		Initial:        initialOf(sandbox, user.Username),
		Email:          user.Email,
		Role:           backofficeauth.RoleName(sandbox, backofficeauth.Role(user.Role)),
		SessionMinutes: sessionMinutes,
		IsRoot:         backofficeauth.Role(user.Role) == backofficeauth.RoleRoot,
	})
}
