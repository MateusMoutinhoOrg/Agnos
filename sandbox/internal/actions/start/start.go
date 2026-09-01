package start

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

func Start(deps *deps.Deps, props api.StartProps) error {
	io := smartio.New(deps, props.Path, props.ProjectName)
	if err := StartInternal(deps, io, props); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(deps, io, props.Path); err != nil {
		return err
	}
	return io.Persist()
}
