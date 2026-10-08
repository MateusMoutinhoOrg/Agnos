package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// VerifyInternal runs every schema check against the transaction-aware io and
// returns a single error listing every violation found, or nil when the tree
// is well-formed.
func VerifyInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("verify started with path %s \n", path)

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	var violations []string
	violations = append(violations, CheckExtensions(sandbox, io)...)
	violations = append(violations, CheckSandbox(sandbox, io, module_conf.Module)...)
	violations = append(violations, CheckContractDocs(sandbox, io)...)
	violations = append(violations, CheckApiShape(sandbox, io)...)
	violations = append(violations, CheckProps(sandbox, io)...)
	violations = append(violations, CheckContractFields(sandbox, io)...)
	violations = append(violations, CheckAdapters(sandbox, io)...)
	violations = append(violations, CheckDepCatalog(sandbox, io, module_conf.Module)...)
	violations = append(violations, CheckAdapterCatalog(sandbox, io, module_conf.Module)...)
	violations = append(violations, CheckRemoteDeps(sandbox, io, path)...)
	violations = append(violations, CheckCommands(sandbox, io)...)
	violations = append(violations, CheckRoutes(sandbox, io)...)
	violations = append(violations, CheckDatabases(sandbox, io)...)
	violations = append(violations, CheckDocs(sandbox, io)...)
	violations = append(violations, CheckStructure(sandbox, io)...)

	if len(violations) == 0 {
		sandbox.Deps.StdDeps.Logf("verify passed\n")
		return nil
	}

	return sandbox.Deps.StdDeps.Errorf("verify found %d violation(s):\n  - %s",
		len(violations), sandbox.Deps.StringsDeps.Join(violations, "\n  - "))
}
