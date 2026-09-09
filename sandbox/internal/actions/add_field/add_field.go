package add_field

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddField appends (or inserts) one field declaration into
// sandbox/internal/routes/<route>/route.yaml, in the origin named by
// props.In, then runs build as a follow-up step so entries.go and the dispatch
// arm pick it up.
func AddField(deps *deps.Deps, props api.RouteFieldProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	if err := AddFieldInternal(deps, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
