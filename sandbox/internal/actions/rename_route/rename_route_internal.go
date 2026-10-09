package rename_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// generatedRouteFiles are the files of a route package every build writes,
// under the names a build before utils.GeneratedPrefix wrote them: a rename
// leaves them behind for the follow-up build to write again, as it does every
// generated.* file.
var generatedRouteFiles = []string{utils.UnitNewFile, utils.UnitInputFile}

// RenameRouteInternal moves every hand-written file of the route's directory
// to <folder>/<name>/ — its own folder, or the one --dir names — rewriting the
// package clause of each Go file, and removes the old directory and every
// folder it leaves empty. Name may be the current one when only the folder
// changes. The route.yaml moves as it is: nothing in it names the package.
//
// The generated health and openapi routes are refused.
func RenameRouteInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.RenameRouteProps) error {
	if err := utils.ValidateRouteName(sandbox, props.Route); err != nil {
		return err
	}
	if err := utils.ValidateRouteName(sandbox, props.Name); err != nil {
		return err
	}

	oldPkg := utils.RoutePackage(sandbox, props.Route)
	newPkg := utils.RoutePackage(sandbox, props.Name)

	for _, pkg := range []string{oldPkg, newPkg} {
		if utils.IsGeneratedRoute(sandbox, pkg) {
			return sandbox.Deps.StdDeps.Errorf("the %s route is generated and cannot be renamed", pkg)
		}
	}
	old_dir, found := utils.FindUnitDir(sandbox, io, utils.RoutesDir, utils.RouteConfFile, oldPkg)
	if !found {
		return sandbox.Deps.StdDeps.Errorf("route %q not found in %s", utils.RouteName(sandbox, props.Route), utils.RoutesDir)
	}

	group := utils.UnitGroupOf(sandbox, utils.RoutesDir, old_dir)
	if props.HasDir {
		moved, err := utils.UnitGroup(sandbox, props.Dir)
		if err != nil {
			return err
		}
		group = moved
	}
	new_dir := utils.UnitDirIn(utils.RoutesDir, group, newPkg)

	if old_dir == new_dir {
		return sandbox.Deps.StdDeps.Errorf("route %q is already named %q in %s", utils.RouteName(sandbox, props.Route), utils.RouteName(sandbox, props.Name), old_dir)
	}
	if oldPkg != newPkg {
		if existing, taken := utils.FindUnitDir(sandbox, io, utils.RoutesDir, utils.RouteConfFile, newPkg); taken {
			return sandbox.Deps.StdDeps.Errorf("route %q already exists in %s", utils.RouteName(sandbox, props.Name), existing)
		}
	}
	if io.IsDir(new_dir) || sandbox.Deps.StringsDeps.HasPrefix(new_dir, old_dir+"/") {
		return sandbox.Deps.StdDeps.Errorf("%s is taken: pick another --dir", new_dir)
	}
	if utils.HoldsOtherUnit(sandbox, io, old_dir, utils.RouteConfFile) {
		return sandbox.Deps.StdDeps.Errorf("route %q holds another route under %s: move that one first", utils.RouteName(sandbox, props.Route), old_dir)
	}

	sandbox.Deps.StdDeps.Logf("rename-route moving %s to %s \n", old_dir, new_dir)

	for _, file := range io.ListFilesRecursively(old_dir) {
		relative := sandbox.Deps.StringsDeps.TrimPrefix(file, old_dir+"/")
		if isGenerated(sandbox, relative) {
			continue
		}
		content, err := io.ReadFile(file)
		if err != nil {
			return err
		}
		if sandbox.Deps.StringsDeps.HasSuffix(relative, ".go") {
			content = []byte(renamePackage(sandbox, string(content), oldPkg, newPkg))
		}
		if err := io.CreateFile(new_dir+"/"+relative, content); err != nil {
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
func isGenerated(sandbox *api.Sandbox, relative string) bool {
	if utils.IsGeneratedFile(sandbox, relative) && !sandbox.Deps.StringsDeps.Contains(relative, "/") {
		return true
	}
	for _, generated := range generatedRouteFiles {
		if relative == generated {
			return true
		}
	}
	return false
}

// renamePackage rewrites the package clause of one Go file, and nothing else.
func renamePackage(sandbox *api.Sandbox, content string, oldPkg string, newPkg string) string {
	lines := sandbox.Deps.StringsDeps.Split(content, "\n")
	for i, line := range lines {
		if line == "package "+oldPkg {
			lines[i] = "package " + newPkg
			break
		}
	}
	return sandbox.Deps.StringsDeps.Join(lines, "\n")
}
