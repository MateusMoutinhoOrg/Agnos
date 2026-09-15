package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveRouteInternal deletes every file under
// sandbox/internal/routes/<name>/ plus the directory itself. The generated
// health route is refused: it is rendered by build, not declared. So is a
// route with an html template beside it — that is a page, and dropping its
// route alone would leave the html orphaned, so remove-page is the editor for
// it.
func RemoveRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateRouteName(sandbox, name); err != nil {
		return err
	}
	pkg := utils.RoutePackage(sandbox, name)
	if pkg == "health" {
		return sandbox.Deps.Std.Errorf("the health route is generated and cannot be removed")
	}
	if utils.IsPage(sandbox, io, name) {
		return sandbox.Deps.Std.Errorf("route %q is a page (%s is beside it): remove it with remove-page, which drops the html too",
			utils.RouteIdentifier(sandbox, name), utils.PageAsset(sandbox, name))
	}

	return RemoveRoutePackage(sandbox, io, name)
}

// RemoveRoutePackage is the deletion itself, without the checks that decide
// whether this route is remove-route's to drop. It is exported for remove-page,
// which owns the routes RemoveRouteInternal refuses and removes the html in the
// same transaction.
func RemoveRoutePackage(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	dir := utils.RouteDir(sandbox, name)
	if !io.IsDir(dir) {
		return sandbox.Deps.Std.Errorf("route %q not found", utils.RouteIdentifier(sandbox, name))
	}

	sandbox.Deps.Std.Log("remove-route removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	return nil
}
