package health

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// RouteHandler answers the built-in health route with a fixed JSON object. It
// is the server layer's `version` command: a route agnos writes itself, so a
// freshly initialized server already answers something.
func RouteHandler(sandbox *api.Sandbox, entries *Entries, response serverdeps.Response) int {
	body := sandbox.Deps.Serializables.CreateObject()
	body.AddItemToObject("status", "ok")

	response.SetHeader("Content-Type", "application/json")
	response.SetStatus(api.StatusOk)
	response.Write([]byte(sandbox.Deps.Serializables.SerializeToJson(body)))

	return api.StatusOk
}
