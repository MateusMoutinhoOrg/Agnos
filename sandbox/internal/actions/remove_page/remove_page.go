package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemovePage deletes one page — assets/frontend/<name>.html — then runs build
// as a follow-up step. The build renders only: a page carries no Go.
func RemovePage(sandbox *api.Sandbox, path string, name string) error {
	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := RemovePageInternal(sandbox, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
