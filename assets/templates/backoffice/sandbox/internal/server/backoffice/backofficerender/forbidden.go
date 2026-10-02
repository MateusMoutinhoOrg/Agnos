package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backofficedb"
)

// ForbiddenPage is what backoffice/forbidden.html is rendered with.
type ForbiddenPage struct {
	Viewer Viewer
}

// Forbidden answers, under a 403, the page telling user that what they asked
// for needs the root role.
func Forbidden(sandbox *api.Sandbox, response *serverdeps.Response, user *backofficedb.BackofficeuserItem) error {
	return Html(sandbox, response, api.StatusForbidden, "backoffice/forbidden.html", ForbiddenPage{Viewer: viewerOf(sandbox, user)})
}
