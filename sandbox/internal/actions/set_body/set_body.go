package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetBody rewrites the `body` keys of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step so ReadBody and the dispatch arm pick the new envelope up.
func SetBody(deps *deps.Deps, props api.RouteBodyProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	if err := SetBodyInternal(deps, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
