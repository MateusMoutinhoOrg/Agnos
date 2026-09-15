package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// VerifyInternal runs every schema check against the transaction-aware io and
// returns a single error listing every violation found, or nil when the tree
// is well-formed.
func VerifyInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("verify started with path %s \n", path)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	var violations []string
	violations = append(violations, CheckSandbox(sandbox, io, module_conf.Module)...)
	violations = append(violations, CheckContracts(sandbox, io)...)
	violations = append(violations, CheckApiShape(sandbox, io)...)
	violations = append(violations, CheckAdapters(sandbox, io)...)
	violations = append(violations, CheckDeplist(sandbox, io, module_conf.Module)...)
	violations = append(violations, CheckAdapterlist(sandbox, io, module_conf.Module)...)
	violations = append(violations, CheckRemoteDeps(sandbox, io, path)...)
	violations = append(violations, CheckRoutes(sandbox, io)...)
	violations = append(violations, CheckDocs(sandbox, io)...)
	violations = append(violations, CheckStructure(sandbox, io)...)

	if len(violations) == 0 {
		sandbox.Deps.Std.Log("verify passed\n")
		return nil
	}

	return sandbox.Deps.Std.Errorf("verify found %d violation(s):\n  - %s",
		len(violations), sandbox.Deps.Stringsdeps.Join(violations, "\n  - "))
}
