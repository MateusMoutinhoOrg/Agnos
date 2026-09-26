package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveRouteInternal deletes every file under
// sandbox/internal/routeslist/<name>/ plus the directory itself. The generated
// health route is refused: it is rendered by build, not declared.
func RemoveRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateRouteName(sandbox, name); err != nil {
		return err
	}
	pkg := utils.RoutePackage(sandbox, name)
	if pkg == "health" {
		return sandbox.Deps.Std.Errorf("the health route is generated and cannot be removed")
	}

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
