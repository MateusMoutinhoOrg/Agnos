package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateBuiltinRouteYamls renders the route.yaml of every route the server
// group writes itself — health and openapi — before the routes are collected:
// the server layer's GenerateHelpCommandYaml. The group renders the same bytes
// again at the end of the build; rendering them first is what lets a project
// whose server group just gained a route collect it in the same build, and
// get its generated.new.go and generated.input.go along with its handler.go.
func GenerateBuiltinRouteYamls(sandbox *api.Sandbox, io *stagedfs.StagedFS, vars map[string]interface{}) error {
	for _, name := range utils.GeneratedRoutes() {
		file := utils.RoutesDir + "/" + name + "/" + utils.RouteConfFile
		if err := utils.RenderTemplateToDest(sandbox, io, utils.ExtensionServer+"/"+file, vars, file); err != nil {
			return err
		}
	}
	return nil
}
