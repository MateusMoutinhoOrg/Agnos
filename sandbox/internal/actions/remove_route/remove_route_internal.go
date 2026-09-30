package remove_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoveRouteInternal deletes every file of the route's directory — in whatever
// folder it sits — plus the directory itself, and every folder the removal
// leaves empty. One holding another route below it is refused, and so is the
// generated health route: it is rendered by build, not declared.
func RemoveRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, name string) error {
	if err := utils.ValidateRouteName(sandbox, name); err != nil {
		return err
	}
	pkg := utils.RoutePackage(sandbox, name)
	if pkg == "health" {
		return sandbox.Deps.Std.Errorf("the health route is generated and cannot be removed")
	}

	dir := utils.RouteDir(sandbox, io, name)
	if !io.IsDir(dir) {
		return sandbox.Deps.Std.Errorf("route %q not found", utils.RouteIdentifier(sandbox, name))
	}

	if utils.HoldsOtherUnit(sandbox, io, dir, utils.RouteConfFile) {
		return sandbox.Deps.Std.Errorf("route %q holds another route under %s: move or remove that one first", utils.RouteIdentifier(sandbox, name), dir)
	}

	sandbox.Deps.Std.Log("remove-route removing %s \n", dir)

	for _, file := range io.ListAllRecursively(dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(dir)
	utils.PruneEmptyGroups(sandbox, io, utils.RoutesDir, dir)
	return nil
}
