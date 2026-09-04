package verify

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// VerifyInternal runs every schema check against the transaction-aware io and
// returns a single error listing every violation found, or nil when the tree
// is well-formed.
func VerifyInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Log("verify started with path %s \n", path)

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}

	var violations []string
	violations = append(violations, CheckSandbox(deps, io, module_conf.Module)...)
	violations = append(violations, CheckContracts(deps, io)...)
	violations = append(violations, CheckAdapters(deps, io)...)
	violations = append(violations, CheckDeplist(deps, io, module_conf.Module)...)
	violations = append(violations, CheckDocs(deps, io)...)
	violations = append(violations, CheckStructure(deps, io)...)
	violations = append(violations, CheckCommandsDoc(deps, io)...)

	if len(violations) == 0 {
		deps.Std.Log("verify passed\n")
		return nil
	}

	return deps.Std.Errorf("verify found %d violation(s):\n  - %s",
		len(violations), strings.Join(violations, "\n  - "))
}
