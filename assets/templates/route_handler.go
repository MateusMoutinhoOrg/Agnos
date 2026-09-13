package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
)

func RouteHandler(sandbox *api.Sandbox, route *api.Route, response serverdeps.Response) int {
	response.SetHeader("Content-Type", "text/plain")
	response.SetStatus(api.StatusOk)
	response.Write([]byte("{{.Identifier}} called\n"))
	return api.StatusOk
}
