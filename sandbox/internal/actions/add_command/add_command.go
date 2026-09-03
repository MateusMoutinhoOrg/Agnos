package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddCommand scaffolds a new command package under
// sandbox/internal/commands/<name>/ — a hand-written entries.yaml and a stub
// handler.go — then runs build as a follow-up step so entries.go and the
// dispatch layer are generated for it.
func AddCommand(deps *deps.Deps, path string, name string, help string, category string) error {
	io := smartio.New(deps, path, config.ProjectName)
	if err := AddCommandInternal(deps, io, name, help, category); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
