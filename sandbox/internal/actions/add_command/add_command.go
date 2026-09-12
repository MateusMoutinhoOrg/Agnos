package add_command

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// AddCommand scaffolds a new command package under
// sandbox/internal/commands/<name>/ — a hand-written entries.yaml and a stub
// handler.go — then runs build as a follow-up step so entries.go and the
// dispatch layer are generated for it.
func AddCommand(sandbox *api.Sandbox, path string, name string, help string, category string) error {
	io := smartio.New(sandbox, path, config.ProjectName)
	if err := AddCommandInternal(sandbox, io, name, help, category); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
