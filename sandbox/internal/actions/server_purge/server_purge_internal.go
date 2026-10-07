package server_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	frontPurgeAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_purge"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// serverDirs are the directories the server layer owns whole. The asset group
// only names the files it installs, so removing those one by one would leave
// the generated neighbours behind — a route's new.go with no route.yaml
// and no handler.go next to it. The server layer is generated from end to end,
// so purging it means dropping these directories entirely.
//
// The start_server command goes with them, from whatever folder of
// sandbox/internal/commands it sits in: it is the entry point server-init
// writes, and it would not compile without the layer it starts.
// sandbox/constructors/server goes with them for the same reason: it fills
// Sandbox.Server by naming the package being removed.
//
// docs/Routes is one of them: the asset group installs its doc.md and
// props.yaml, but the build writes one page per route beside them. Removing
// only the installed two would leave a directory of pages with no props.yaml,
// which every later build reads as a doc that fails to load.
//
// sandbox/internal/routeprops is one of them: every route names it, so it goes
// with them.
var serverDirs = []string{
	"sandbox/internal/server",
	utils.RoutesDir,
	utils.RoutePropsDir,
	utils.GeneratedDir + "/server",
	"docs/Routes",
	utils.ConstructorDir("server"),
}

// ServerPurgeInternal removes from the target project every file that the
// "server" asset group would have installed, at the path it holds inside that
// group, plus the directories the server layer owns whole, then drops any
// directory the removal left empty.
//
// The front layer goes with it, since it has nothing to be served by without
// it. The cli layer is deliberately left in place: server-init may have installed
// it, but a cli, once there, is the project's. The deps the server layer
// pulled in are left too — other code may use them.
func ServerPurgeInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("server-purge started with path %s \n", path)

	// The front layer is served by a route of this one, so it goes first, on
	// this same open SmartIO: a front left on with no server to serve it is a
	// declaration verify refuses. Its pages under assets/frontend/ stay.
	has_front, err := utils.ExtensionEnabled(sandbox, io, utils.ExtensionSandboxFront)
	if err != nil {
		return err
	}
	if has_front {
		if err := frontPurgeAction.FrontPurgeInternal(sandbox, io, path); err != nil {
			return err
		}
	}

	files, err := utils.ExtensionFiles(sandbox, utils.ExtensionSandboxServer)
	if err != nil {
		return err
	}

	for _, file := range files {
		io.RemoveDir(file)
	}

	// sandbox/api/routeprops.go is where an older build declared
	// RouteProps; a tree no build has moved it out of yet still carries it.
	io.RemoveDir("sandbox/api/" + utils.RoutePropsFile)

	for _, dir := range append(serverDirs, utils.CommandDir(sandbox, io, "start_server")) {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
	}

	utils.RemoveRetiredGenerated(sandbox, io, utils.ExtensionSandboxServer)

	for _, dir := range ancestorDirs(sandbox, files) {
		if len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxServer, false)
}

// ancestorDirs returns every directory that contains one of the given files,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(sandbox *api.Sandbox, files []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		for i := 1; i < len(parts); i++ {
			dir := sandbox.Deps.Stringsdeps.Join(parts[:i], "/")
			if !seen[dir] {
				seen[dir] = true
				dirs = append(dirs, dir)
			}
		}
	}
	sandbox.Deps.Sortdeps.Slice(dirs, func(i, j int) bool {
		return sandbox.Deps.Stringsdeps.Count(dirs[i], "/") > sandbox.Deps.Stringsdeps.Count(dirs[j], "/")
	})
	return dirs
}
