package front_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// serverDir is what tells a project that already has a server layer from one
// that has to be given one first — the same directory `build` reads hasServer
// from.
const serverDir = "sandbox/internal/server"

// staticStyle and staticScript are the two files the skeleton carries. They
// are not decoration: //go:embed keeps no empty directory, so a `dirref
// "styles"` against a directory that does not exist is an error at render
// time, and the scaffolded page names both directories.
const (
	staticStyle  = utils.StaticAssetsDir + "/styles/main.css"
	staticScript = utils.StaticAssetsDir + "/scripts/main.js"
)

// FrontInitInternal renders every embedded asset under assets/front into the
// target project at the path it holds inside that group, then writes the two
// halves that are the project's from the moment they exist: the static route
// and the skeleton of assets/frontend/.
//
// A project with no server layer is given one first, on this same open
// SmartIO: actions compose by sharing one transaction, so there is no
// intermediate Persist and no intermediate build between the two halves.
func FrontInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("front-init started with path %s \n", path)

	if !io.IsDir(serverDir) {
		if err := serverInitAction.ServerInitInternal(sandbox, io, path); err != nil {
			return err
		}
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module": module_conf.Module,
	}

	if err := writeStaticRoute(sandbox, io, vars); err != nil {
		return err
	}

	if err := writeFrontSkeleton(sandbox, io, vars); err != nil {
		return err
	}

	// The static route is written before the group is rendered because
	// pageio's StaticMount is read off that declaration. The build that
	// follows renders it again from the same source, so the only thing this
	// ordering buys is a first render that is already correct rather than one
	// pointing at a mount the project may have renamed.
	vars["StaticMount"] = buildAction.CollectFrontMount(sandbox, io)

	return utils.RenderGroup(sandbox, io, "front", vars)
}

// writeStaticRoute scaffolds the route serving assets/frontend/static, leaving
// an existing one alone: like any route's handler.go it is written once and
// then the project's — and its safeSegments check is what stands between a
// caller's path and the rest of the asset tree, so re-rendering over an edited
// copy would undo a deliberate change without saying so.
func writeStaticRoute(sandbox *api.Sandbox, io *smartio.SmartIO, vars map[string]interface{}) error {
	dir := utils.RouteDir(sandbox, utils.StaticRouteName)

	if io.IsDir(dir) {
		sandbox.Deps.Std.Log("front-init: %s already exists, keeping it \n", dir)
		return nil
	}

	if err := utils.RenderTemplateToDest(sandbox, io, "templates/static_route.yaml", vars, dir+"/route.yaml"); err != nil {
		return err
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/static_handler.go", vars, dir+"/handler.go")
}

// writeFrontSkeleton writes the two starting files of assets/frontend/ through
// io.WriteFile, which refuses to overwrite: everything under that tree is the
// project's content, so a second front-init — after a front-purge, say — finds
// them already there and leaves what was written in between untouched.
//
// assets/frontend/pages/ is deliberately left empty; add-page is what fills it.
func writeFrontSkeleton(sandbox *api.Sandbox, io *smartio.SmartIO, vars map[string]interface{}) error {
	skeleton := []struct {
		template string
		dest     string
	}{
		{"templates/front_main.css", staticStyle},
		{"templates/front_main.js", staticScript},
	}

	for _, file := range skeleton {
		if io.IsFile(file.dest) {
			sandbox.Deps.Std.Log("front-init: %s already exists, keeping it \n", file.dest)
			continue
		}

		content, err := sandbox.Deps.Embeddeps.RenderTemplate(file.template, vars)
		if err != nil {
			return err
		}
		if err := io.WriteFile(file.dest, content); err != nil {
			return err
		}
	}

	return nil
}
