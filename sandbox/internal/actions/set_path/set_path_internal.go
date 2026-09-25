package set_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// SetPathInternal parses the target route's route.yaml, rebuilds the named
// entry of `paths` with the changes applied and writes the file back.
func SetPathInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.RoutePathEditProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	index := utils.FindRoutePath(sandbox, conf.Paths, props.Id)
	if index < 0 {
		return sandbox.Deps.Std.Errorf("route %q declares no path with id %q", props.Route, utils.RouteEntryId(sandbox, props.Id))
	}
	if utils.RoutePathEditEmpty(sandbox, props) {
		return sandbox.Deps.Std.Errorf("set-path was given nothing to change")
	}

	current := conf.Paths[index]
	path, err := utils.RoutePathEdited(sandbox, current, props)
	if err != nil {
		return err
	}
	if utils.RouteIdTaken(conf, path.Id, current.Id) {
		return sandbox.Deps.Std.Errorf("route %q already has an Entries field named %q", props.Route, path.Id)
	}

	sandbox.Deps.Std.Log("set-path rewriting %s of %s \n", current.Id, utils.RouteConfPath(sandbox, props.Route))

	conf.Paths[index] = path
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
