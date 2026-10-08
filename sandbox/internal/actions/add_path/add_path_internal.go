package add_path

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddPathInternal parses the target route's route.yaml, inserts the new entry
// into `paths` and writes the file back. Every entry of `paths` has to match
// for the route to run, so where it lands only changes the order route.yaml
// and the docs read in.
func AddPathInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddPathProps) error {
	conf, err := utils.LoadRouteConf(sandbox, io, props.Route)
	if err != nil {
		return err
	}

	position, err := utils.CheckRoutePosition(sandbox, "path", props.Position, len(conf.Paths))
	if err != nil {
		return err
	}

	path, err := utils.NewRoutePath(sandbox, props)
	if err != nil {
		return err
	}
	if utils.RouteIdTaken(conf, path.Id, "") {
		return sandbox.Deps.StdDeps.Errorf("route %q already has an Input field named %q", props.Route, path.Id)
	}

	sandbox.Deps.StdDeps.Logf("add-path adding %s to %s \n", path.Id, utils.RouteConfPath(sandbox, io, props.Route))

	conf.Paths = utils.InsertAt(conf.Paths, path, position)
	return utils.SaveRouteConf(sandbox, io, props.Route, conf)
}
