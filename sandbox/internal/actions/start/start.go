package start

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// Start scaffolds a project and leaves it built. The first build runs inside
// the same transaction as the skeleton, because the declarations it reads are
// not on disk yet; the follow-up build runs against the persisted tree, which
// is what collects the contracts that first build wrote into sandbox/api/ —
// SmartIO listings read disk, so a contract rendered into the transaction is
// not one the build that rendered it can already see.
func Start(sandbox *api.Sandbox, props api.StartProps) error {
	io := smartio.New(sandbox, props.Path, props.ProjectName)
	if err := StartInternal(sandbox, io, props); err != nil {
		return err
	}
	if err := buildAction.BuildInternal(sandbox, io, props.Path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: props.Path, Runtime: api.RuntimeGo})
}
