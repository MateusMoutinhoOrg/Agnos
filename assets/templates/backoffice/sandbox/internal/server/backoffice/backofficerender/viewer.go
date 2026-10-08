package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
	"{{.Module}}/sandbox/internal/server/backoffice/backofficeauth"
)

// Viewer is the signed-in user, as the top bar of every admin page shows it.
type Viewer struct {
	Username string
	// Initial is the first letter of Username, shown in the avatar.
	Initial string
	// IsRoot shows what only a root may do.
	IsRoot bool
}

// viewerOf is the Viewer of user.
func viewerOf(sandbox *api.Sandbox, user *backoffice_db.BackofficeUserRecord) Viewer {
	return Viewer{
		Username: user.Username,
		Initial:  initialOf(sandbox, user.Username),
		IsRoot:   backofficeauth.Role(user.Role) == backofficeauth.RoleRoot,
	}
}
