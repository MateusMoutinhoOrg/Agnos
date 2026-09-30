package rename_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// generatedRouteFiles are the files of a route package every build writes, so
// a rename leaves them behind for the follow-up build to write again.
var generatedRouteFiles = []string{"new.go", "entries.go"}

// RenameRouteInternal moves every hand-written file of the route's directory
// to <folder>/<name>/ — its own folder, or the one --dir names — rewriting the
// package clause of each Go file, and removes the old directory and every
// folder it leaves empty. Name may be the current one when only the folder
// changes. The route.yaml moves as it is: nothing in it names the package.
//
// The generated health route is refused.
func RenameRouteInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RenameRouteProps) error {
	if err := utils.ValidateRouteName(sandbox, props.Route); err != nil {
		return err
	}
	if err := utils.ValidateRouteName(sandbox, props.Name); err != nil {
		return err
	}

	old_pkg := utils.RoutePackage(sandbox, props.Route)
	new_pkg := utils.RoutePackage(sandbox, props.Name)

	if old_pkg == "health" || new_pkg == "health" {
		return sandbox.Deps.Std.Errorf("the health route is generated and cannot be renamed")
	}
	old_dir, found := utils.FindUnitDir(sandbox, io, utils.RoutesDir, utils.RouteConfFile, old_pkg)
	if !found {
		return sandbox.Deps.Std.Errorf("route %q not found in %s", utils.RouteIdentifier(sandbox, props.Route), utils.RoutesDir)
	}

	group := utils.UnitGroupOf(sandbox, utils.RoutesDir, old_dir)
	if props.HasDir {
		moved, err := utils.UnitGroup(sandbox, props.Dir)
		if err != nil {
			return err
		}
		group = moved
	}
	new_dir := utils.UnitDirIn(utils.RoutesDir, group, new_pkg)

	if old_dir == new_dir {
		return sandbox.Deps.Std.Errorf("route %q is already named %q in %s", utils.RouteIdentifier(sandbox, props.Route), utils.RouteIdentifier(sandbox, props.Name), old_dir)
	}
	if old_pkg != new_pkg {
		if existing, taken := utils.FindUnitDir(sandbox, io, utils.RoutesDir, utils.RouteConfFile, new_pkg); taken {
			return sandbox.Deps.Std.Errorf("route %q already exists in %s", utils.RouteIdentifier(sandbox, props.Name), existing)
		}
	}
	if io.IsDir(new_dir) || sandbox.Deps.Stringsdeps.HasPrefix(new_dir, old_dir+"/") {
		return sandbox.Deps.Std.Errorf("%s is taken: pick another --dir", new_dir)
	}
	if utils.HoldsOtherUnit(sandbox, io, old_dir, utils.RouteConfFile) {
		return sandbox.Deps.Std.Errorf("route %q holds another route under %s: move that one first", utils.RouteIdentifier(sandbox, props.Route), old_dir)
	}

	sandbox.Deps.Std.Log("rename-route moving %s to %s \n", old_dir, new_dir)

	for _, file := range io.ListFilesRecursively(old_dir) {
		relative := sandbox.Deps.Stringsdeps.TrimPrefix(file, old_dir+"/")
		if isGenerated(relative) {
			continue
		}
		content, err := io.ReadFile(file)
		if err != nil {
			return err
		}
		if sandbox.Deps.Stringsdeps.HasSuffix(relative, ".go") {
			content = []byte(renamePackage(sandbox, string(content), old_pkg, new_pkg))
		}
		if err := io.WriteFile(new_dir+"/"+relative, content); err != nil {
			return err
		}
	}

	for _, file := range io.ListAllRecursively(old_dir) {
		io.RemoveDir(file)
	}
	io.RemoveDir(old_dir)
	utils.PruneEmptyGroups(sandbox, io, utils.RoutesDir, old_dir)
	return nil
}

// isGenerated reports a file of the package the follow-up build writes again.
func isGenerated(relative string) bool {
	for _, generated := range generatedRouteFiles {
		if relative == generated {
			return true
		}
	}
	return false
}

// renamePackage rewrites the package clause of one Go file, and nothing else.
func renamePackage(sandbox *api.Sandbox, content string, old_pkg string, new_pkg string) string {
	lines := sandbox.Deps.Stringsdeps.Split(content, "\n")
	for i, line := range lines {
		if line == "package "+old_pkg {
			lines[i] = "package " + new_pkg
			break
		}
	}
	return sandbox.Deps.Stringsdeps.Join(lines, "\n")
}
