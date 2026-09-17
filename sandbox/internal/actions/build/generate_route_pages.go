package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routePagesDir is the doc whose doc.md indexes one page per route.
const routePagesDir = utils.DocsDir + "/Routes"

// GenerateRoutePages is the server layer's GenerateCommandPages: it renders
// assets/templates/route_page.md once per visible route into
// docs/Routes/<name>.md, the page docs/Routes' own doc.md links to. A route is
// named by the package that declares it, so the page of a route is found from
// its directory and the other way round.
//
// A page whose route is gone is removed here, so the directory holds the routes
// that are declared now and no page nothing links to.
func GenerateRoutePages(sandbox *api.Sandbox, io *smartio.SmartIO, groups []RouteDocGroup) error {
	written := map[string]bool{}

	for _, group := range groups {
		for _, route := range group.Routes {
			file, err := routePageFile(sandbox, route.Name)
			if err != nil {
				return err
			}

			vars := map[string]any{
				"Category": group.Category,
				"Route":    route,
			}

			if err := utils.RenderTemplateToDest(sandbox, io, "templates/route_page.md", vars, file); err != nil {
				return err
			}
			written[file] = true
		}
	}

	removeStaleDocPages(sandbox, io, routePagesDir, written)
	return nil
}

// routePageFile is the page a route is written to. A name that spells one of
// the two reserved names of a doc directory is a hard error rather than a page
// silently overwriting the index it is linked from.
func routePageFile(sandbox *api.Sandbox, name string) (string, error) {
	file := name + docPageExt

	if file == utils.DocFile || file == utils.DocIndexFile {
		return "", sandbox.Deps.Std.Errorf(
			"route %s cannot be documented: its page would be %s/%s, which is the doc's own %s",
			name, routePagesDir, file, file)
	}

	return routePagesDir + "/" + file, nil
}
