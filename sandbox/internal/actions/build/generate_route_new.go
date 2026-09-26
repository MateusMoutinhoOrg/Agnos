package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateRouteNew renders assets/templates/route_new.go and
// assets/templates/route_entries.go once per route into
// sandbox/internal/routeslist/<name>/: new.go, the api.Route that package
// declares — a 1:1 image of its route.yaml — which
// sandbox/internal/generated/server/server/new.go collects into Server.Routes, and
// entries.go, the Entries struct its InternalPureHandler is handed plus the
// ReadBody its body declaration calls for. It is the server layer's
// GenerateCommandNew; the module path is merged in because both files import
// the project's own packages.
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

		dir := routesDir + "/" + name
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/route_new.go", vars, dir+"/new.go"); err != nil {
			return err
		}
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/route_entries.go", vars, dir+"/entries.go"); err != nil {
			return err
		}
	}
	return nil
}
