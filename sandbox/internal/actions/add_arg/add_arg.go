package add_arg

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/config"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// AddArg appends (or inserts at --position) one positional arg declaration
// into sandbox/internal/commands/<command>/entries.yaml, then runs build as a
// follow-up step so entries.go and the dispatch layer pick it up.
func AddArg(deps *deps.Deps, props api.FieldProps) error {
	io := smartio.New(deps, props.Path, config.ProjectName)
	if err := AddArgInternal(deps, io, props); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(deps, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
