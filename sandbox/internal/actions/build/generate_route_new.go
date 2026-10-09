package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateRouteNew renders assets/templates/route_new.go and
// assets/templates/route_input.go once per route into
// the directory it sits in under sandbox/internal/routes: generated.new.go,
// the api.Route that package declares — a 1:1 image of its route.yaml — which
// sandbox/internal/server/generated.new.go collects into Server.Routes, and
// generated.input.go, the Input struct its Handle is handed plus the
// ReadBody its body declaration calls for. It is the server layer's
// GenerateCommandNew; the module path is merged in because both files import
// the project's own packages.
func GenerateRouteNew(sandbox *api.Sandbox, io *stagedfs.StagedFS, routes []map[string]any, module string) error {
	for _, route := range routes {
		name, _ := route["RouteName"].(string)
		if name == "" {
			continue
		}

		vars := map[string]any{"Module": module, "GeneratorName": generatorName(sandbox)}
		for key, value := range route {
			vars[key] = value
		}

		dir, _ := route["Dir"].(string)
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/route_new.go", vars, dir+"/"+utils.GeneratedFile(sandbox, utils.UnitNewFile)); err != nil {
			return err
		}
		if err := utils.RenderTemplateToDest(sandbox, io, "templates/route_input.go", vars, dir+"/"+utils.GeneratedFile(sandbox, utils.UnitInputFile)); err != nil {
			return err
		}
	}
	return nil
}
