package add_doc

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddDoc scaffolds a new doc directory under docs/ — a doc.md stub and the
// props.yaml declaring it — then runs build as a follow-up step so the theme
// indexes and the parent's Index.md list it.
func AddDoc(deps *deps.Deps, props api.DocProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	if err := AddDocInternal(deps, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
