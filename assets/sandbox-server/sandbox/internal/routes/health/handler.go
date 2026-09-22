package health

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

// RouteHandler answers the built-in health route with a fixed JSON object. It
// is the server layer's `version` command: a route agnos writes itself, so a
// freshly initialized server already answers something.
func RouteHandler(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) error {
	body := sandbox.Deps.Serializables.CreateObject()
	body.AddItemToObject("status", "ok")

	response.SetHeader("Content-Type", "application/json")
	response.SetStatus(api.StatusOk)
	response.Write([]byte(sandbox.Deps.Serializables.SerializeToJson(body)))

	return nil
}
