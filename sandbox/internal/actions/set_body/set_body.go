package set_body

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// SetBody rewrites the `body` keys of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step so ReadBody and the dispatch arm pick the new envelope up.
func SetBody(sandbox *api.Sandbox, props api.RouteBodyProps) error {
	io := smartio.New(sandbox, props.Path, config.ProjectName)
	if err := SetBodyInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
