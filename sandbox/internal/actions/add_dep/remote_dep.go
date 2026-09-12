package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/apishape"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// RemoteApiDir is the contract half of an agnos repo, the directory a consumer
// copies. It imports nothing but its own sandbox/deps, by the rule every agnos
// repo lives under, and that one import exists for one field — Sandbox.Deps,
// which the copy drops. So what lands in the consumer is self-contained: the
// package clause changes, and the wiring stays behind.
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
func AddRemoteDepInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.AddDepProps) error {
	module, version := splitModuleSpec(sandbox, props.Dep)

	name := props.As
	if name == "" {
		name = defaultDepName(sandbox, module)
	}

	sandbox.Deps.Std.Log("add-dep started with path %s module %s as %s \n", props.Path, module, name)

	if err := utils.ValidateDepName(sandbox, name); err != nil {
		return err
	}

	resolved, err := resolveModule(sandbox, props.Path, module, version)
	if err != nil {
		return err
	}

	remote, err := ReadRemoteApi(sandbox, resolved.Dir)
	if err != nil {
		return err
	}

	if violations := apishape.Violations(sandbox, remote); len(violations) > 0 {
		return sandbox.Deps.Std.Errorf("%s@%s cannot be installed as a dep, its api is not convertible:\n  - %s",
			resolved.Path, resolved.Version, sandbox.Deps.Stringsdeps.Join(violations, "\n  - "))
	}

	if err := CopyRemoteApi(sandbox, io, remote, name); err != nil {
		return err
	}

	available := props.RemoteAvailable
	if available == "" {
		available = utils.StandardAvailable
	}

	if err := GenerateShim(sandbox, io, remote, ShimProps{
		Dep:       name,
		Module:    resolved.Path,
		Available: available,
		HasDeps:   sandbox.Deps.Iodeps.IsDir(sandbox.Deps.Iodeps.Join(resolved.Dir, utils.AvailableDir(available))),
	}); err != nil {
		return err
	}

	adapter_conf := adapterconf.NewEmpty(sandbox)
	adapter_conf.Name = name
	adapter_conf.Dep = name
	adapter_conf.Help = "Generated shim over " + resolved.Path
	adapter_conf.Module = resolved.Path + "@" + resolved.Version
	adapter_conf.Origin = adapterconf.OriginGenerated

	if err := io.WriteFileOverwrite(utils.AdapterConfPath(name), []byte(adapter_conf.Render())); err != nil {
		return err
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}
	module_conf.AddRequire(resolved.Path + " " + resolved.Version)
	if err := io.WriteFileOverwrite("go.mod", []byte(module_conf.Render())); err != nil {
		return err
	}

	return utils.EnrollAdapter(sandbox, io, name)
}

// ReadRemoteApi parses the sandbox/api/ of a resolved module, straight off
// disk: the module cache is outside the project, so it is read through the
// filesystem contract rather than the project's transactional io.
func ReadRemoteApi(sandbox *api.Sandbox, dir string) (*apishape.Api, error) {
	api_dir := sandbox.Deps.Iodeps.Join(dir, RemoteApiDir)

	if !sandbox.Deps.Iodeps.IsDir(api_dir) {
		return nil, sandbox.Deps.Std.Errorf("%s has no %s: it is not an agnos repo", dir, RemoteApiDir)
	}

	var sources []apishape.Source
	for _, file := range sandbox.Deps.Iodeps.ListFiles(api_dir) {
		if !sandbox.Deps.Stringsdeps.HasSuffix(file, ".go") {
			continue
		}
		content, err := sandbox.Deps.Iodeps.ReadFile(file)
		if err != nil {
			return nil, err
		}
		sources = append(sources, apishape.Source{Name: baseName(sandbox, file), Content: string(content)})
	}

	if len(sources) == 0 {
		return nil, sandbox.Deps.Std.Errorf("%s is empty: there is no contract to copy", api_dir)
	}

	return apishape.New(sandbox, sources)
}

// CopyRemoteApi writes the remote contract into sandbox/deps/<name>/, one file
// per file, with nothing changed but the package clause and the dependency
// wiring — which is what makes the copy legible as the thing it is, doc
// comments and all.
func CopyRemoteApi(sandbox *api.Sandbox, io *smartio.SmartIO, remote *apishape.Api, name string) error {
	for _, file := range remote.Files {
		content := sandbox.Deps.Stringsdeps.ReplaceAll(file.Content,
			"package "+file.Parsed.Package+"\n", "package "+name+"\n")

		content = stripDepsWiring(sandbox, content)

		formatted, err := sandbox.Deps.Goimportsdeps.Format(content)
		if err != nil {
			return sandbox.Deps.Std.Errorf("%s could not be formatted after the package clause was rewritten: %w", file.Name, err)
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
func splitModuleSpec(sandbox *api.Sandbox, spec string) (string, string) {
	at := sandbox.Deps.Stringsdeps.LastIndex(spec, "@")
	if at < 0 {
		return spec, ""
	}
	return spec[:at], spec[at+1:]
}

// defaultDepName is the name a remote dep takes when --as names none: the last
// segment of the module path, lower-cased, because a dep is a directory and a
// Go package.
func defaultDepName(sandbox *api.Sandbox, module string) string {
	segments := sandbox.Deps.Stringsdeps.Split(module, "/")
	return sandbox.Deps.Stringsdeps.ToLower(segments[len(segments)-1])
}

// baseName is the last segment of a host path.
func baseName(sandbox *api.Sandbox, path string) string {
	segments := sandbox.Deps.Stringsdeps.Split(sandbox.Deps.Stringsdeps.ReplaceAll(path, "\\", "/"), "/")
	return segments[len(segments)-1]
}

// stripDepsWiring removes the one part of an api package that does not cross
// into a consumer: the Sandbox.Deps field and the sandbox/deps import it is
// the only reason for. Deps is how the remote repo reaches the outside world,
// filled by that repo's own adapters and already bound in the compiled sandbox
// the shim builds — a consumer has its own, and naming the remote's would name
// a package it does not have.
//
// The removal is by line because that is the shape the field is written in:
// sandbox/api/sandbox.go is generated, so the field is always its own line and
// the import always its own. An import block left holding nothing goes with
// them, since gofmt keeps an empty one.
func stripDepsWiring(sandbox *api.Sandbox, content string) string {
	var kept []string

	for _, line := range sandbox.Deps.Stringsdeps.Split(content, "\n") {
		trimmed := sandbox.Deps.Stringsdeps.TrimSpace(line)

		if isDepsImportLine(sandbox, trimmed) {
			continue
		}

		if isDepsFieldLine(sandbox, trimmed) {
			kept = dropTrailingComment(sandbox, kept)
			continue
		}

		kept = append(kept, line)
	}

	return dropEmptyImportBlock(sandbox, kept)
}

// isDepsImportLine reports whether one trimmed line imports the loose
// sandbox/deps package, aliased or not. A contract package under
// sandbox/deps/<x>/ ends in another segment, so it does not match.
func isDepsImportLine(sandbox *api.Sandbox, trimmed string) bool {
	if !sandbox.Deps.Stringsdeps.HasSuffix(trimmed, "/sandbox/deps\"") {
		return false
	}
	return sandbox.Deps.Stringsdeps.HasPrefix(trimmed, "\"") ||
		sandbox.Deps.Stringsdeps.HasPrefix(trimmed, "deps \"")
}

// isDepsFieldLine reports whether one trimmed line declares the Sandbox.Deps
// field. It reads the line as fields rather than as text because the file is
// gofmt'ed, so the name and the type are padded apart by however wide the
// widest field of the struct is.
func isDepsFieldLine(sandbox *api.Sandbox, trimmed string) bool {
	parts := sandbox.Deps.Stringsdeps.Fields(trimmed)
	return len(parts) == 2 && parts[0] == apishape.DepsField && parts[1] == "*deps.Deps"
}

// dropTrailingComment removes the doc comment the dropped field was carrying,
// which is every line back to the first that is not one. A comment explaining
// Deps outlives it otherwise, describing a field the copy does not have.
func dropTrailingComment(sandbox *api.Sandbox, kept []string) []string {
	for len(kept) > 0 {
		last := sandbox.Deps.Stringsdeps.TrimSpace(kept[len(kept)-1])
		if !sandbox.Deps.Stringsdeps.HasPrefix(last, "//") {
			break
		}
		kept = kept[:len(kept)-1]
	}
	return kept
}

// dropEmptyImportBlock joins the kept lines back together, leaving out an
// "import (" ... ")" that stripDepsWiring emptied. gofmt keeps an empty block,
// and a contract that declares no import should not carry one.
func dropEmptyImportBlock(sandbox *api.Sandbox, lines []string) string {
	var kept []string

	for index := 0; index < len(lines); index++ {
		if sandbox.Deps.Stringsdeps.TrimSpace(lines[index]) != "import (" {
			kept = append(kept, lines[index])
			continue
		}

		end := index
		for end < len(lines) && sandbox.Deps.Stringsdeps.TrimSpace(lines[end]) != ")" {
			end++
		}

		empty := end < len(lines)
		for inner := index + 1; inner < end; inner++ {
			if sandbox.Deps.Stringsdeps.TrimSpace(lines[inner]) != "" {
				empty = false
			}
		}

		if !empty {
			kept = append(kept, lines[index])
			continue
		}

		index = end
	}

	return sandbox.Deps.Stringsdeps.Join(kept, "\n")
}
