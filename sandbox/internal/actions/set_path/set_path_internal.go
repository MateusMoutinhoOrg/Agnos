package set_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetPathInternal parses the target route's route.yaml, rebuilds the named
// entry of `paths` with the changes applied and writes the file back.
func SetPathInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.SetPathProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	index := utils.FindRoutePath(sandbox, conf.Paths, props.Name)
	if index < 0 {
		return sandbox.Deps.StdDeps.Errorf("route %q declares no path with id %q", props.Route, utils.GoIdentifier(sandbox, props.Name))
	}
	if utils.IsRoutePathEditEmpty(sandbox, props) {
		return sandbox.Deps.StdDeps.Errorf("set-path was given nothing to change")
	}

	current := conf.Paths[index]
	path, err := utils.EditRoutePath(sandbox, current, props)
	if err != nil {
		return err
	}
	if utils.RouteIdTaken(conf, path.Id, current.Id) {
		return sandbox.Deps.StdDeps.Errorf("route %q already has an Input field named %q", props.Route, path.Id)
	}

	sandbox.Deps.StdDeps.Logf("set-path rewriting %s of %s \n", current.Id, utils.RouteConfPath(sandbox, io, props.Route))

	conf.Paths[index] = path
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
