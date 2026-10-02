package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// Html renders the template at path with vars and answers it under status.
// Templates run on text/template, so every value a user can influence is
// escaped in the template itself with {{"{{"}}html ...}}.
func Html(sandbox *api.Sandbox, response *serverdeps.Response, status int, path string, vars any) error {
	content, err := sandbox.Deps.Embeddeps.RenderTemplate(path, vars)
	if err != nil {
		return err
	}
	response.SetHeader("Content-Type", "text/html; charset=utf-8")
	response.SetHeader("Cache-Control", "no-store")
	response.SetStatus(status)
	return response.Write(content)
}
