package backoffice_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// BackofficePurge removes from the project everything backoffice-init wrote,
// then runs build as a follow-up step.
func BackofficePurge(sandbox *api.Sandbox, path string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := BackofficePurgeInternal(sandbox, io, path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
