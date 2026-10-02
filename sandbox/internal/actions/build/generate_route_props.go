package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// retiredServerFiles are generated files an older build wrote that nothing
// renders any more. Each one names a symbol the current contract dropped, so a
// tree still carrying it would not compile; it is removed on every build.
var retiredServerFiles = []string{
	utils.GeneratedDir + "/routeio/locals.go",
}

// GenerateRouteProps rewrites sandbox/internal/routeprops/routeprops.go — the
// RouteProps one request's chain of routes shares, handed to every
// InternalPureHandler as its first argument — as the embedding of every part
// the package declares beside it: project.go, the project's own, and one file
// per mechanic that hands something on (backoffice.go). See
// generatePropsAggregate.
func GenerateRouteProps(sandbox *api.Sandbox, io *smartio.SmartIO, module string) error {
	for _, retired := range retiredServerFiles {
		if io.IsFile(retired) {
			io.RemoveDir(retired)
		}
	}

	return generatePropsAggregate(sandbox, io, propsAggregate{
		Dir:             utils.RoutePropsDir,
		File:            utils.RoutePropsFile,
		Type:            "RouteProps",
		Template:        "templates/routeprops.go",
		ProjectTemplate: "templates/routeprops_project.go",
	}, module)
}
