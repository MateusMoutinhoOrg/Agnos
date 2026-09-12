package add_header

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddHeader inserts one header declaration into the `headers` of
// sandbox/internal/routes/<route>/route.yaml, then runs build as a follow-up
// step so entries.go and the dispatch arm pick it up.
func AddHeader(sandbox *api.Sandbox, props api.RouteFieldProps) error {
	io := smartio.New(sandbox, props.Path, config.ProjectName)
	if err := AddHeaderInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
