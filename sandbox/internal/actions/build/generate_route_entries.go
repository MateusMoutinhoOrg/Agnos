package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateRouteEntries renders assets/templates/route_entries.go once per route
// into sandbox/internal/routes/<name>/entries.go — the typed struct the user's
// handler.go receives, plus the ReadBody its body declaration calls for,
// derived from that route's route.yaml. It is the server layer's
// GenerateCommandEntries; the module path is merged in because a route's
// entries.go imports the project's own packages.
func GenerateRouteEntries(deps *deps.Deps, io *smartio.SmartIO, routes []map[string]any, module string) error {
	for _, route := range routes {
		name, _ := route["Name"].(string)
		if name == "" {
			continue
		}

		vars := map[string]any{"Module": module}
		for key, value := range route {
			vars[key] = value
		}

		dest := routesDir + "/" + name + "/entries.go"
		if err := utils.RenderTemplateToDest(deps, io, "templates/route_entries.go", vars, dest); err != nil {
			return err
		}
	}
	return nil
}
