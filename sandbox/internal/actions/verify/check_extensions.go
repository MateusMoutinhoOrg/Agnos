package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CheckExtensions enforces that <ProjectName>Config/extensions.yaml declares
// the mechanics this agnos knows, and only those. The declaration is the whole
// of what `build` reads to decide what to render, so a key nothing answers to
// is a mechanic the author believes is on while nothing generates it.
//
// It also enforces the one dependency the catalog has: every sandbox-<x>
// mechanic renders into the sandbox, so none of them can be on while the
// sandbox itself is off.
func CheckExtensions(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string

	extensions_conf, err := utils.LoadExtensionsConf(sandbox, io)
	if err != nil {
		return append(violations, err.Error())
	}

	for _, extension := range extensions_conf.Extensions {
		if !utils.IsExtensionName(extension.Name) {
			violations = append(violations, utils.ExtensionsConfPath(sandbox)+" declares "+extension.Name+
				", which is not an extension (known: "+
				sandbox.Deps.Stringsdeps.Join(utils.ExtensionNames(), ", ")+")")
		}
	}

	if extensions_conf.IsEnabled(utils.ExtensionSandbox) {
		return violations
	}

	for _, spec := range utils.ExtensionCatalog() {
		if spec.Name == utils.ExtensionSandbox {
			continue
		}
		if !sandbox.Deps.Stringsdeps.HasPrefix(spec.Name, utils.ExtensionSandbox+"-") {
			continue
		}
		if extensions_conf.IsEnabled(spec.Name) {
			violations = append(violations, utils.ExtensionsConfPath(sandbox)+" has "+spec.Name+
				" on with "+utils.ExtensionSandbox+" off (it renders into the sandbox and has nothing to render into)")
		}
	}

	return violations
}
