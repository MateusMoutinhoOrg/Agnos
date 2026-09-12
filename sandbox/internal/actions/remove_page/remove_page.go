package remove_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// RemovePage deletes one page whole — its route package and the html template
// it renders — then runs build as a follow-up step so the dispatch stops
// naming it. The build renders only: dropping a page may leave hand-written
// code referring to what is gone.
func RemovePage(sandbox *api.Sandbox, path string, name string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := RemovePageInternal(sandbox, io, name); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeNone})
}
