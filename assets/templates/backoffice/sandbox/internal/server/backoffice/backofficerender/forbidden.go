package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/databases/backoffice_db"
)

// ForbiddenPage is what backoffice/forbidden.html is rendered with.
type ForbiddenPage struct {
	Viewer Viewer
}

// RenderForbiddenPage answers, under a 403, the page telling user that what they asked
// for needs the root role.
func RenderForbiddenPage(sandbox *api.Sandbox, response *serverdeps.Response, user *backoffice_db.BackofficeUserRecord) error {
	return RenderHTML(sandbox, response, api.StatusForbidden, "backoffice/forbidden.html", ForbiddenPage{Viewer: viewerOf(sandbox, user)})
}
