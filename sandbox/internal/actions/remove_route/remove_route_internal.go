package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveRouteInternal deletes every file under
// sandbox/internal/routes/<name>/ plus the directory itself. The generated
// health route is refused: it is rendered by build, not declared.
func RemoveRouteInternal(deps *deps.Deps, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateRouteName(deps, name); err != nil {
		return err
	}
	pkg := utils.RoutePackage(deps, name)
	if pkg == "health" {
		return deps.Std.Errorf("the health route is generated and cannot be removed")
	}

	dir := utils.RouteDir(deps, name)
	if !io.IsDir(dir) {
		return deps.Std.Errorf("route %q not found", utils.RouteIdentifier(deps, name))
	}

	deps.Std.Log("remove-route removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	return nil
}
