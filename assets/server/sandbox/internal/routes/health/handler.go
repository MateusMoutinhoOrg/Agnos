package health

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
)

// Handle answers the built-in health route with a fixed JSON
// object. It is the server layer's `version` command: a route agnos writes
// itself, so a freshly initialized server already answers something.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	body := sandbox.Deps.SerializableDeps.CreateObject()
	body.AddItemToObject("status", "ok")

	response.SetStatus(api.StatusOK)
	response.Write([]byte(sandbox.Deps.SerializableDeps.SerializeToJson(body)))

	return nil
}
