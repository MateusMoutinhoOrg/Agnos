package backofficerender

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// RenderHTML renders the template at path with vars and answers it under status.
// Templates run on html/template, so every value is escaped for the context
// it lands in — text, an attribute, a url — whether or not the template
// spells {{"{{"}}html ...}} too; a field added without it is still safe.
func RenderHTML(sandbox *api.Sandbox, response *serverdeps.Response, status int, path string, vars any) error {
	content, err := sandbox.Deps.EmbedDeps.RenderHTMLTemplate(path, vars)
	if err != nil {
		return err
	}
	response.SetHeader("Content-Type", "text/html; charset=utf-8")
	response.SetHeader("Cache-Control", "no-store")
	response.SetStatus(status)
	return response.Write(content)
}
