package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateRouteProps rewrites sandbox/internal/routeprops/routeprops.go — the
// RouteProps one request's chain of routes shares, handed to every
// Handle as its first argument — as the embedding of every part
// the package declares beside it: project.go, the project's own, and one file
// per mechanic that hands something on (backoffice.go). See
// generatePropsAggregate. What an older build wrote under
// sandbox/internal/generated/routeio — locals.go among it — goes with the rest
// of that retired package, through utils.RemoveRetiredGenerated.
func GenerateRouteProps(sandbox *api.Sandbox, io *stagedfs.StagedFS, module string) error {
	return generatePropsAggregate(sandbox, io, propsAggregate{
		Dir:             utils.RoutePropsDir,
		File:            utils.RoutePropsFile,
		Type:            "RouteProps",
		Template:        "templates/routeprops.go",
		ProjectTemplate: "templates/routeprops_project.go",
	}, module)
}
