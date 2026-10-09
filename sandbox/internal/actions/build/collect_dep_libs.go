package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// collectLibDirs lists the immediate sub-directories of dir and returns one
// {"Name": <dir>, "Title": <Field>} entry per sub-directory, in listing order,
// Title spelled the way utils.DepField spells a Deps field. Both the
// sandbox/deps and adapters/impls trees have this fixed shape: one
// sub-directory per sub-contract.
func collectLibDirs(sandbox *api.Sandbox, io *stagedfs.StagedFS, dir string) []map[string]string {

	dirs := io.ListDirs(dir)

	var libs []map[string]string
	for _, d := range dirs {
		parts := sandbox.Deps.StringsDeps.Split(d, "/")
		name := parts[len(parts)-1]

		if len(name) == 0 {
			continue
		}

		libs = append(libs, map[string]string{
			"Name":  name,
			"Title": utils.DepField(sandbox, name),
		})
	}

	return libs
}

// CollectDepLibs returns one entry per sandbox/deps sub-contract directory,
// for the {{range .DepLibs}} loop in sandbox/deps/generated.deps.go. Type is the root
// type of the contract: Contract for every one the catalog or a hand writes,
// Sandbox for a remote dep, which keeps the name of the api it was copied
// from.
func CollectDepLibs(sandbox *api.Sandbox, io *stagedfs.StagedFS) []map[string]string {
	libs := collectLibDirs(sandbox, io, "sandbox/deps")
	for _, lib := range libs {
		lib["Type"] = utils.ContractType
		if conf, err := utils.LoadAdapterConf(sandbox, io, lib["Name"]); err == nil && conf.Origin == adapterconf.OriginGenerated {
			lib["Type"] = utils.RemoteContractType
		}
	}
	return libs
}
