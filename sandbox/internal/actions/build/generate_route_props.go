package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// routePropsDest is where routeprops.RouteProps is declared: a package of its
// own under sandbox/internal, since every route package names it, none of them
// may import another, and its fields may name any type of the project.
const routePropsDest = utils.RoutePropsDir + "/" + utils.RoutePropsFile

// retiredServerFiles are generated files an older build wrote that nothing
// renders any more. Each one names a symbol the current contract dropped, so a
// tree still carrying it would not compile; it is removed on every build.
var retiredServerFiles = []string{
	utils.GeneratedDir + "/routeio/locals.go",
}

// GenerateRouteProps renders assets/templates/routeprops.go into
// sandbox/internal/routeprops/routeprops.go — the RouteProps one request's chain of routes
// shares, handed to every InternalPureHandler as its first argument.
//
// It is written **once**, like the Handle* files: what a request carries from
// a middleware to the routes after it — the user it authenticated, say — is
// the project's to type, so a struct already on disk is left as it is. Writing
// it here rather than in server-init carries a project that ran server-init
// before it existed.
func GenerateRouteProps(sandbox *api.Sandbox, io *smartio.SmartIO, module string) error {
	for _, retired := range retiredServerFiles {
		if io.IsFile(retired) {
			io.RemoveDir(retired)
		}
	}

	// Read rather than IsFile: a struct MigrateLegacyProps moved here this
	// build is still pending, and IsFile only sees the disk.
	if _, err := io.ReadFile(routePropsDest); err == nil {
		return nil
	}

	vars := map[string]any{
		"Module":        module,
		"GeneratorName": generatorName(sandbox),
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/routeprops.go", vars, routePropsDest)
}
