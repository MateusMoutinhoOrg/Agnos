package verify

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/internal/smartio"
)

// VerifyInternal runs every schema check against the transaction-aware io and
// returns a single error listing every violation found, or nil when the tree
// is well-formed.
func VerifyInternal(deps *deps.Deps, io *smartio.SmartIO, path string) error {
	deps.Std.Printf("verify started with path %s \n", path)

	gomod, err := io.ReadFile(path + "/go.mod")
	if err != nil {
		return err
	}
	module_conf, err := moduleconf.New(deps, string(gomod))
	if err != nil {
		return err
	}

	var violations []string
	violations = append(violations, CheckSandbox(deps, io, module_conf.Module)...)
	violations = append(violations, CheckAdapters(deps, io)...)

	if len(violations) == 0 {
		deps.Std.Printf("verify passed")
		return nil
	}

	return deps.Std.Errorf("verify found %d violation(s):\n  - %s",
		len(violations), strings.Join(violations, "\n  - "))
}
