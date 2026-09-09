package add_param

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddParam inserts one query-parameter declaration into the `params` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step so entries.go and the dispatch arm pick it up.
func AddParam(deps *deps.Deps, props api.RouteFieldProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	if err := AddParamInternal(deps, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
