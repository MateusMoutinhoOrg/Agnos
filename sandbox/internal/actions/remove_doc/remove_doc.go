package remove_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// RemoveDoc deletes a whole doc directory of docs/ — its sub-docs and assets
// included — then runs build so the indexes that listed it are rewritten
// without it.
func RemoveDoc(sandbox *api.Sandbox, props api.RemoveDocProps) error {
	io := stagedfs.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := RemoveDocInternal(sandbox, io, props.Name); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeNone})
}
