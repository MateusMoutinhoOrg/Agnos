package front_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// FrontPurge removes from the project every file the "front" asset group
// would have installed, then runs build as a follow-up step.
func FrontPurge(sandbox *api.Sandbox, path string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := FrontPurgeInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
