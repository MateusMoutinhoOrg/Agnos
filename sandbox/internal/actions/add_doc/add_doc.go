package add_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddDoc scaffolds a new doc directory under docs/ — a doc.md stub and the
// props.yaml declaring it — then runs build as a follow-up step so the theme
// indexes and the parent's Index.md list it.
func AddDoc(sandbox *api.Sandbox, props api.DocProps) error {
	io := smartio.New(sandbox, props.Path, config.ProjectName)
	if err := AddDocInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
