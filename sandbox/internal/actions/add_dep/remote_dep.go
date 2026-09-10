package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/apishape"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoteApiDir is the contract half of an agnos repo, the directory a consumer
// copies. It imports nothing, by the rule every agnos repo lives under, so the
// copy is self-contained: only the package clause changes.
const RemoteApiDir = "sandbox/api"

// AddRemoteDepInternal installs another agnos repo as a dep. The repo's own
// `sandbox/api/` becomes this project's `sandbox/deps/<name>/`, and the adapter
// filling it is generated: it builds the remote sandbox out of the remote
// repo's own adapters and converts the result into the copied contract.
//
// The conversion cannot come from the remote side. `<name>.Sandbox` lives at an
// import path inside *this* module, one the remote repo does not know when it
// publishes, so any Go that names that type has to compile here. And a
// top-level cast will not do instead: Go's type identity does not reach through
// named types, so two copies of a struct whose fields are named types are never
// the same type.
func AddRemoteDepInternal(deps *deps.Deps, io *smartio.SmartIO, props api.AddDepProps) error {
	module, version := splitModuleSpec(deps, props.Dep)

	name := props.As
	if name == "" {
		name = defaultDepName(deps, module)
	}

	deps.Std.Log("add-dep started with path %s module %s as %s \n", props.Path, module, name)

	if err := utils.ValidateDepName(deps, name); err != nil {
		return err
	}

	resolved, err := resolveModule(deps, props.Path, module, version)
	if err != nil {
		return err
	}

	remote, err := ReadRemoteApi(deps, resolved.Dir)
	if err != nil {
		return err
	}

	if violations := apishape.Violations(deps, remote); len(violations) > 0 {
		return deps.Std.Errorf("%s@%s cannot be installed as a dep, its api is not convertible:\n  - %s",
			resolved.Path, resolved.Version, deps.Stringsdeps.Join(violations, "\n  - "))
	}

	if err := CopyRemoteApi(deps, io, remote, name); err != nil {
		return err
	}

	available := props.RemoteAvailable
	if available == "" {
		available = utils.StandardAvailable
	}

	if err := GenerateShim(deps, io, remote, ShimProps{
		Dep:       name,
		Module:    resolved.Path,
		Available: available,
		HasDeps:   deps.Iodeps.IsDir(deps.Iodeps.Join(resolved.Dir, utils.AvailableDir(available))),
	}); err != nil {
		return err
	}

	adapter_conf := adapterconf.NewEmpty(deps)
	adapter_conf.Name = name
	adapter_conf.Dep = name
	adapter_conf.Help = "Generated shim over " + resolved.Path
	adapter_conf.Module = resolved.Path + "@" + resolved.Version
	adapter_conf.Origin = adapterconf.OriginGenerated

	if err := io.WriteFileOverwrite(utils.AdapterConfPath(name), []byte(adapter_conf.Render())); err != nil {
		return err
	}

	module_conf, err := utils.LoadModuleConf(deps, io)
	if err != nil {
		return err
	}
	module_conf.AddRequire(resolved.Path + " " + resolved.Version)
	if err := io.WriteFileOverwrite("go.mod", []byte(module_conf.Render())); err != nil {
		return err
	}

	return utils.EnrollAdapter(deps, io, name)
}

// ReadRemoteApi parses the sandbox/api/ of a resolved module, straight off
// disk: the module cache is outside the project, so it is read through the
// filesystem contract rather than the project's transactional io.
func ReadRemoteApi(deps *deps.Deps, dir string) (*apishape.Api, error) {
	api_dir := deps.Iodeps.Join(dir, RemoteApiDir)

	if !deps.Iodeps.IsDir(api_dir) {
		return nil, deps.Std.Errorf("%s has no %s: it is not an agnos repo", dir, RemoteApiDir)
	}

	var sources []apishape.Source
	for _, file := range deps.Iodeps.ListFiles(api_dir) {
		if !deps.Stringsdeps.HasSuffix(file, ".go") {
			continue
		}
		content, err := deps.Iodeps.ReadFile(file)
		if err != nil {
			return nil, err
		}
		sources = append(sources, apishape.Source{Name: baseName(deps, file), Content: string(content)})
	}

	if len(sources) == 0 {
		return nil, deps.Std.Errorf("%s is empty: there is no contract to copy", api_dir)
	}

	return apishape.New(deps, sources)
}

// CopyRemoteApi writes the remote contract into sandbox/deps/<name>/, one file
// per file, with nothing changed but the package clause — which is what makes
// the copy legible as the thing it is, doc comments and all.
func CopyRemoteApi(deps *deps.Deps, io *smartio.SmartIO, remote *apishape.Api, name string) error {
	for _, file := range remote.Files {
		content := deps.Stringsdeps.ReplaceAll(file.Content,
			"package "+file.Parsed.Package+"\n", "package "+name+"\n")

		formatted, err := deps.Goimportsdeps.Format(content)
		if err != nil {
			return deps.Std.Errorf("%s could not be formatted after the package clause was rewritten: %w", file.Name, err)
		}

		if err := io.WriteFileOverwrite(utils.ContractsDir+"/"+name+"/"+file.Name, []byte(formatted)); err != nil {
			return err
		}
	}

	return nil
}

// splitModuleSpec cuts "<module>@<version>" apart. A spec with no version
// leaves the version empty, which resolves to whatever the build list already
// holds, or to latest.
func splitModuleSpec(deps *deps.Deps, spec string) (string, string) {
	at := deps.Stringsdeps.LastIndex(spec, "@")
	if at < 0 {
		return spec, ""
	}
	return spec[:at], spec[at+1:]
}

// defaultDepName is the name a remote dep takes when --as names none: the last
// segment of the module path, lower-cased, because a dep is a directory and a
// Go package.
func defaultDepName(deps *deps.Deps, module string) string {
	segments := deps.Stringsdeps.Split(module, "/")
	return deps.Stringsdeps.ToLower(segments[len(segments)-1])
}

// baseName is the last segment of a host path.
func baseName(deps *deps.Deps, path string) string {
	segments := deps.Stringsdeps.Split(deps.Stringsdeps.ReplaceAll(path, "\\", "/"), "/")
	return segments[len(segments)-1]
}
