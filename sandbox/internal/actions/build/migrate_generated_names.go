package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// movedGeneratedPackages pairs each package an older build rendered under
// utils.LegacyGeneratedDir with the mechanic that renders it now, at
// sandbox/internal/<name>/generated.new.go.
var movedGeneratedPackages = [][2]string{
	{"cli", utils.ExtensionCli},
	{"server", utils.ExtensionServer},
	{"config", utils.ExtensionSandbox},
}

// migratedRoots are the trees whose Go files may import a moved package: the
// sandbox, the adapters and the entry point.
var migratedRoots = []string{"sandbox", "adapters", "cmd"}

// MigrateGeneratedNames moves a tree an older build wrote, before
// utils.GeneratedPrefix named a generated file, onto the current layout. For
// every moved package whose mechanic is on, it rewrites the import of
// <module>/sandbox/internal/generated/<name> to <module>/sandbox/internal/<name>
// in every file the project owns — the constructor.go start wrote once
// imports it — and drops the old package, which the mechanic's group renders
// again at its new path in this same build. utils.LegacyGeneratedDir goes when
// nothing is left in it.
//
// The files an older build wrote beside the project's under their old names
// (a command's new.go) go as the build writes their generated.<name>
// (utils.WriteGenerated); until then the collectors count them once
// (utils.DropSuperseded). The files of a remote dep, which no build renders,
// are moved here (migrateRemoteDeps).
func MigrateGeneratedNames(sandbox *api.Sandbox, io *stagedfs.StagedFS, module string, enabled func(extension string) bool) error {
	if err := migrateRemoteDeps(sandbox, io); err != nil {
		return err
	}
	if !io.IsDir(utils.LegacyGeneratedDir) {
		return nil
	}

	replacements := map[string]string{}
	for _, moved := range movedGeneratedPackages {
		legacy := utils.LegacyGeneratedDir + "/" + moved[0]
		if !enabled(moved[1]) || !io.IsDir(legacy) {
			continue
		}
		replacements[`"`+module+"/"+legacy+`"`] = `"` + module + "/sandbox/internal/" + moved[0] + `"`
	}
	if len(replacements) == 0 {
		return nil
	}

	for _, root := range migratedRoots {
		for _, file := range io.ListFilesRecursively(root) {
			if !sandbox.Deps.StringsDeps.HasSuffix(file, ".go") || sandbox.Deps.StringsDeps.HasPrefix(file, utils.LegacyGeneratedDir+"/") {
				continue
			}
			content, err := io.ReadFile(file)
			if err != nil {
				return err
			}
			text := string(content)
			for old, current := range replacements {
				text = sandbox.Deps.StringsDeps.ReplaceAll(text, old, current)
			}
			if text == string(content) {
				continue
			}
			if err := io.WriteFile(file, []byte(text)); err != nil {
				return err
			}
			sandbox.Deps.StdDeps.Logf("build: %s now imports the generated packages from sandbox/internal/\n", file)
		}
	}

	for old := range replacements {
		dir := sandbox.Deps.StringsDeps.TrimPrefix(sandbox.Deps.StringsDeps.Trim(old, `"`), module+"/")
		io.RemoveDir(dir)
	}
	return nil
}

// RemoveEmptyLegacyGeneratedDir drops utils.LegacyGeneratedDir once nothing
// is left in it: after MigrateGeneratedNames moved its packages and
// utils.RemoveRetiredGenerated dropped the ones the libs replaced.
func RemoveEmptyLegacyGeneratedDir(sandbox *api.Sandbox, io *stagedfs.StagedFS) {
	if !io.IsDir(utils.LegacyGeneratedDir) {
		return
	}
	if len(io.ListFilesRecursively(utils.LegacyGeneratedDir)) > 0 {
		return
	}
	io.RemoveDir(utils.LegacyGeneratedDir)
}

// migrateRemoteDeps moves what an add-dep <module> before utils.GeneratedPrefix
// wrote — the copy of the module's api under sandbox/deps/<dep>/ and the shim
// adapters/impls/<dep>/<dep>.go — to generated.<name>, the header on top. Only
// set-dep rewrites those files, so a build that waited for it would leave them
// under the old names until the next version bump.
func migrateRemoteDeps(sandbox *api.Sandbox, io *stagedfs.StagedFS) error {
	for _, adapter := range utils.InstalledAdapters(sandbox, io) {
		adapter_conf, err := utils.LoadAdapterConf(sandbox, io, adapter)
		if err != nil || adapter_conf.Origin != adapterconf.OriginGenerated {
			continue
		}

		files := []string{utils.AdapterDir(adapter) + "/" + adapter_conf.Dep + ".go"}
		for _, file := range io.ListFiles(utils.ContractsDir + "/" + adapter_conf.Dep) {
			if sandbox.Deps.StringsDeps.HasSuffix(file, ".go") {
				files = append(files, file)
			}
		}

		for _, file := range files {
			if utils.IsGeneratedFile(sandbox, file) || !io.IsFile(file) {
				continue
			}
			content, err := io.ReadFile(file)
			if err != nil {
				return err
			}
			name := lastSegmentOf(sandbox, file)
			dest := sandbox.Deps.StringsDeps.TrimSuffix(file, name) + utils.GeneratedFile(sandbox, name)
			if err := utils.WriteGenerated(sandbox, io, dest, content); err != nil {
				return err
			}
		}
	}
	return nil
}
