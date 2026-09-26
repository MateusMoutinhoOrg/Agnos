package rename_route

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// generatedRouteFiles are the files of a route package every build writes, so
// a rename leaves them behind for the follow-up build to write again.
var generatedRouteFiles = []string{"new.go", "entries.go"}

// RenameRouteInternal moves every hand-written file of
// sandbox/internal/routeslist/<route>/ to sandbox/internal/routeslist/<name>/,
// rewriting the package clause of each Go file, and removes the old
// directory. The route.yaml moves as it is: nothing in it names the package.
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
	old_dir := utils.RouteDir(sandbox, props.Route)
	new_dir := utils.RouteDir(sandbox, props.Name)

	if old_pkg == "health" || new_pkg == "health" {
		return sandbox.Deps.Std.Errorf("the health route is generated and cannot be renamed")
	}
	if !io.IsDir(old_dir) {
		return sandbox.Deps.Std.Errorf("route %q not found in %s", utils.RouteIdentifier(sandbox, props.Route), old_dir)
	}
	if old_pkg == new_pkg {
		return sandbox.Deps.Std.Errorf("route %q is already named %q", utils.RouteIdentifier(sandbox, props.Route), utils.RouteIdentifier(sandbox, props.Name))
	}
	if io.IsDir(new_dir) {
		return sandbox.Deps.Std.Errorf("route %q already exists in %s", utils.RouteIdentifier(sandbox, props.Name), new_dir)
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
