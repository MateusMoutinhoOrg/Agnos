package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/moduleconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// InstallAdapter renders one adapter of the embedded catalog into the target
// project and writes its declaration beside the package, so the installed tree
// answers by itself which dep that adapter fills. The versioned module the
// adapter imports — "" for one that needs nothing beyond the stdlib — is
// pinned in go.mod here, where the import actually is, rather than under the
// name of the contract. The adapter is then enrolled in every available, since
// a contract nothing binds is a nil func waiting to panic.
func InstallAdapter(deps *deps.Deps, io *smartio.SmartIO, adapter_conf *adapterconf.AdapterConf, module_conf *moduleconf.ModuleConf, vars map[string]interface{}) error {

	group := utils.AdapterlistGroup + "/" + adapter_conf.Name
	if err := utils.RenderGroupExcept(deps, io, group, vars, []string{utils.AdapterConfFile}); err != nil {
		return err
	}

	if err := io.WriteFileOverwrite(utils.AdapterConfPath(adapter_conf.Name), []byte(adapter_conf.Render())); err != nil {
		return err
	}

	if err := utils.EnrollAdapter(deps, io, adapter_conf.Name); err != nil {
		return err
	}

	module, version, ok := adapter_conf.ModuleSpec()
	if !ok {
		return nil
	}

	module_conf.AddRequire(module + " " + version)
	return io.WriteFileOverwrite("go.mod", []byte(module_conf.Render()))
}
