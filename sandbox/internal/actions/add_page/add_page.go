package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddPage scaffolds a new page: the route package answering it under
// sandbox/internal/routes/<name>/ and the html template it renders under
// assets/frontend/pages/, then runs build as a follow-up step so entries.go
// and the dispatch arm are generated for it.
func AddPage(deps *deps.Deps, props api.PageProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	if err := AddPageInternal(deps, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
