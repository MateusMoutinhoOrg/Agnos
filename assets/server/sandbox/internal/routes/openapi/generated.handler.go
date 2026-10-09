package openapi

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
)

// document is the OpenAPI 3.0.3 document of every visible route — the bytes
// docs/Routes/openapi.json holds — rendered by `{{.GeneratorName}} build` from
// each route.yaml, so it is always the route set this binary was built with.
const document = {{printf "%q" (or .OpenApi "")}}

// Handle answers the built-in openapi route with document. Like health, it is
// a route the build writes itself; any origin may read it, so a Swagger UI served
// from elsewhere opens it by its address.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	response.SetHeader("Access-Control-Allow-Origin", "*")
	response.SetStatus(api.StatusOK)
	response.Write([]byte(document))

	return nil
}
