package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// Start scaffolds a project and leaves it built. The first build runs inside
// the same transaction as the skeleton, because the declarations it reads are
// not on disk yet; the follow-up build runs against the persisted tree, which
// is what collects the contracts that first build wrote into sandbox/api/ —
// SmartIO listings read disk, so a contract rendered into the transaction is
// not one the build that rendered it can already see.
//
// The Go toolchain is looked for before anything is written: the follow-up
// build runs it, and a start that fails there would leave a project that was
// scaffolded but never built.
func Start(sandbox *api.Sandbox, props api.StartProps) error {
	if _, err := sandbox.Deps.Rundeps.Run(rundeps.RunProps{Program: "go", Args: []string{"version"}}); err != nil {
		return sandbox.Deps.Std.Errorf("start needs the Go toolchain (go %s or later) on PATH to build the project it scaffolds: %w", utils.GoFloor, err)
	}
	io := smartio.New(sandbox, props.Path, sandbox.Config.ProjectName)
	if err := StartInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	return buildAction.PersistAndBuild(sandbox, io, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
