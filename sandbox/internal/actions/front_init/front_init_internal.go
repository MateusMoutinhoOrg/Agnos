package front_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// FrontInitInternal turns the front mechanic on in the project's declaration,
// then writes what is the project's from the moment it exists: the frontend
// route, and the index and 404 pages of assets/frontend/. The group
// itself is rendered by the follow-up build, like every other mechanic.
//
// A project with no server layer is given one first, on this same open
// SmartIO: actions compose by sharing one transaction, so there is no
// intermediate Persist and no intermediate build between the two halves.
func FrontInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("front-init started with path %s \n", path)

	has_server, err := utils.ExtensionEnabled(sandbox, io, utils.ExtensionSandboxServer)
	if err != nil {
		return err
	}
	if !has_server {
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

	if err := writeFrontendRoute(sandbox, io, vars); err != nil {
		return err
	}

	if err := writeIndexPage(sandbox, io); err != nil {
		return err
	}

	if err := writeNotFoundPage(sandbox, io); err != nil {
		return err
	}

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxFront, true)
}

// writeFrontendRoute scaffolds the route serving assets/frontend, leaving an
// existing one alone: like any route's InternalPureHandler.go it is written
// once and then the project's — turning spaFallback on is an edit to it — so
// re-rendering over an edited copy would undo a deliberate change without
// saying so. The path check it relies on lives in the generated frontio, so a
// fix to it reaches the project on the next build whatever this file holds.
func writeFrontendRoute(sandbox *api.Sandbox, io *smartio.SmartIO, vars map[string]interface{}) error {
	dir := utils.RouteDir(sandbox, utils.FrontendRouteName)

	if io.IsDir(dir) {
		sandbox.Deps.Std.Log("front-init: %s already exists, keeping it \n", dir)
		return nil
	}

	if err := utils.RenderTemplateToDest(sandbox, io, "templates/frontend_route.yaml", vars, dir+"/route.yaml"); err != nil {
		return err
	}
	return utils.RenderTemplateToDest(sandbox, io, "templates/frontend_handler.go", vars, dir+"/InternalPureHandler.go")
}

// writeIndexPage writes assets/frontend/index.html, the page "/" answers, the
// way add-page writes any other — and like it, keeps one already there:
// everything under that tree is the project's content, so a second front-init
// (after a front-purge, say) leaves what was written in between untouched.
// It is also what keeps the tree from being empty, which //go:embed would drop.
func writeIndexPage(sandbox *api.Sandbox, io *smartio.SmartIO) error {
	page := utils.PageAsset(sandbox, utils.FrontendIndexPage)
	if io.IsFile(page) {
		sandbox.Deps.Std.Log("front-init: %s already exists, keeping it \n", page)
		return nil
	}
	return utils.WritePage(sandbox, io, utils.FrontendIndexPage, "Home")
}

// writeNotFoundPage writes assets/frontend/404.html, the formatted page the
// frontend route answers with a 404 when a path names no file, and like
// writeIndexPage keeps one already there: it is the project's content from the
// moment it exists, so restyling it is an edit to that file.
func writeNotFoundPage(sandbox *api.Sandbox, io *smartio.SmartIO) error {
	page := utils.PageAsset(sandbox, utils.FrontendNotFoundPage)
	if io.IsFile(page) {
		sandbox.Deps.Std.Log("front-init: %s already exists, keeping it \n", page)
		return nil
	}
	return utils.WriteNotFoundPage(sandbox, io)
}
