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
// It also enforces what each mechanic requires (utils.ExtensionRequires):
// every sandbox-<x> mechanic renders into the sandbox, and the front layer is
// served by the server layer, so none can be on while what it needs is off.
//
// And it enforces that a mechanic with code has the opinated lib that code
// lives in (utils.OpinatedLib): its -init installs it, but enable-extension
// turns the key on alone, and the build then renders api aliases of a
// contract that is not there.
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

	for _, spec := range utils.ExtensionCatalog() {
		if !extensions_conf.IsEnabled(spec.Name) {
			continue
		}
		for _, required := range utils.ExtensionRequires(spec.Name) {
			if !extensions_conf.IsEnabled(required) {
				violations = append(violations, utils.ExtensionsConfPath(sandbox)+" has "+spec.Name+
					" on with "+required+" off (it has nothing to render into)")
			}
		}
	}

	for _, spec := range utils.ExtensionCatalog() {
		lib := utils.OpinatedLib(spec.Name)
		if lib == "" || !extensions_conf.IsEnabled(spec.Name) || io.IsDir(utils.ContractsDir+"/"+lib) {
			continue
		}
		violations = append(violations, utils.ExtensionsConfPath(sandbox)+" has "+spec.Name+
			" on with no "+utils.ContractsDir+"/"+lib+" (the layer's code is that lib: add-dep "+lib+")")
	}

	return violations
}
