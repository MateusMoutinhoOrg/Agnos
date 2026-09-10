package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CheckRemoteDeps is CheckDeplist for a dep copied from another agnos repo:
// sandbox/deps/<dep>/ must still be the module's own sandbox/api/, file for
// file. Nothing else keeps the two in step, so a contract that gained a field
// upstream would otherwise go on being called through the shape it had when it
// was copied.
//
// No hash is recorded anywhere: the module cache is immutable per version and
// go.sum already signs its content, so the comparison is against the cache
// itself. The check is skipped when the module is not in the cache — it runs
// with the proxy off, so a verify never reaches the network to answer it.
func CheckRemoteDeps(deps *deps.Deps, io *smartio.SmartIO, path string) []string {
	var violations []string

	for _, adapter := range utils.InstalledAdapters(deps, io) {
		adapter_conf, err := utils.LoadAdapterConf(deps, io, adapter)
		if err != nil || adapter_conf.Origin != adapterconf.OriginGenerated {
			continue
		}

		module, _, ok := adapter_conf.ModuleSpec()
		if !ok {
			continue
		}

		dir, ok := cachedModuleDir(deps, path, module)
		if !ok {
			continue
		}

		violations = append(violations, checkRemoteCopy(deps, io, adapter_conf.Dep, dir)...)
	}

	return violations
}

// cachedModuleDir asks the go toolchain where a module's source is, with the
// proxy off so the answer can only come from what is already on this machine.
func cachedModuleDir(deps *deps.Deps, path string, module string) (string, bool) {
	result, err := deps.Rundeps.Run(rundeps.RunProps{
		Dir:     path,
		Program: "go",
		Args:    []string{"list", "-m", "-json", module},
		Env:     []string{"GOPROXY=off"},
	})
	if err != nil || result.ExitCode != 0 {
		return "", false
	}

	parsed, err := deps.Serializables.ParseJson(result.Output)
	if err != nil {
		return "", false
	}

	item, _ := parsed.GetObjectItem("Dir")
	if item == nil || item.IsNull() {
		return "", false
	}

	dir, err := item.GetString()
	if err != nil || dir == "" {
		return "", false
	}

	return dir, true
}

// checkRemoteCopy compares one copied contract with the module it came from,
// rendering the module's own api the way the copy was written: the package
// clause rewritten, then formatted.
func checkRemoteCopy(deps *deps.Deps, io *smartio.SmartIO, dep string, dir string) []string {
	remote, err := add_dep.ReadRemoteApi(deps, dir)
	if err != nil {
		return []string{"dep " + dep + " could not be compared with the module it was copied from: " + err.Error()}
	}

	var violations []string

	for _, file := range remote.Files {
		target := utils.ContractsDir + "/" + dep + "/" + file.Name

		installed, err := io.ReadFile(target)
		if err != nil {
			violations = append(violations, target+" is missing, but the module it was copied from declares it"+
				" (run `agnos set-dep "+dep+" --version <version>`)")
			continue
		}

		expected := deps.Stringsdeps.ReplaceAll(file.Content,
			"package "+file.Parsed.Package+"\n", "package "+dep+"\n")

		formatted, err := deps.Goimportsdeps.Format(expected)
		if err != nil {
			violations = append(violations, target+" could not be compared: "+err.Error())
			continue
		}

		if formatted == string(installed) {
			continue
		}

		violations = append(violations, target+" has drifted from the module it was copied from"+
			" (run `agnos set-dep "+dep+" --version <version>`)")
	}

	return violations
}
