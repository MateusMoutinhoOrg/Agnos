package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateRouteNew renders assets/templates/route_new.go once per route into
// sandbox/internal/routes/<name>/new.go — the api.Route that package declares,
// derived from its route.yaml, which sandbox/internal/server/new.go collects into
// Server.Routes, plus the ReadBody its body declaration calls for. It is the
// server layer's GenerateCommandNew; the module path is merged in because a
// route's new.go imports the project's own packages.
func GenerateRouteNew(sandbox *api.Sandbox, io *smartio.SmartIO, routes []map[string]any, module string) error {
	for _, route := range routes {
		name, _ := route["Name"].(string)
		if name == "" {
			continue
		}

		vars := map[string]any{"Module": module, "GeneratorName": generatorName(sandbox)}
		for key, value := range route {
			vars[key] = value
		}

		dest := routesDir + "/" + name + "/new.go"
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/route_new.go", vars, dest); err != nil {
			return err
		}
	}
	return nil
}
