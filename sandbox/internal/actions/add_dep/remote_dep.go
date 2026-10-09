package add_dep

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/apishape"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/declarations/adapterconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
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
func AddRemoteDepInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, props api.AddDepProps) error {
	module, version := splitModuleSpec(sandbox, props.Dep)

	name := props.As
	if name == "" {
		name = defaultDepName(sandbox, module)
	}

	sandbox.Deps.StdDeps.Logf("add-dep started with path %s module %s as %s \n", props.Path, module, name)

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
		return sandbox.Deps.StdDeps.Errorf("%s@%s cannot be installed as a dep, its api is not convertible:\n  - %s",
			resolved.Path, resolved.Version, sandbox.Deps.StringsDeps.Join(violations, "\n  - "))
	}

	if err := CopyRemoteApi(sandbox, io, remote, name); err != nil {
		return err
	}

	binding := props.RemoteBinding
	if binding == "" {
		binding = utils.StandardBinding
	}

	binding_dir, has_deps := utils.RemoteBindingDir(sandbox, resolved.Dir, binding)
	if err := GenerateShim(sandbox, io, remote, ShimProps{
		Dep:        name,
		Module:     resolved.Path,
		Binding:    binding,
		BindingDir: binding_dir,
		HasDeps:    has_deps,
	}); err != nil {
		return err
	}

	adapter_conf := adapterconf.NewEmpty(sandbox)
	adapter_conf.Name = name
	adapter_conf.Dep = name
	adapter_conf.Help = "Generated shim over " + resolved.Path
	adapter_conf.Module = resolved.Path + "@" + resolved.Version
	adapter_conf.Origin = adapterconf.OriginGenerated

	if err := io.WriteFile(utils.AdapterConfPath(name), []byte(adapter_conf.Render())); err != nil {
		return err
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}
	module_conf.AddRequire(resolved.Path + " " + resolved.Version)
	if err := io.WriteFile("go.mod", []byte(module_conf.Render())); err != nil {
		return err
	}

	return utils.EnrollAdapter(sandbox, io, name)
}

// ReadRemoteApi parses the sandbox/api/ of a resolved module, straight off
// disk: the module cache is outside the project, so it is read through the
// filesystem contract rather than the project's transactional io.
func ReadRemoteApi(sandbox *api.Sandbox, dir string) (*apishape.Api, error) {
	api_dir := sandbox.Deps.IoDeps.Join(dir, RemoteApiDir)

	if !sandbox.Deps.IoDeps.IsDir(api_dir) {
		return nil, sandbox.Deps.StdDeps.Errorf("%s has no %s: it is not an agnos repo", dir, RemoteApiDir)
	}

	var sources []apishape.Source
	for _, file := range sandbox.Deps.IoDeps.ListFiles(api_dir) {
		if !sandbox.Deps.StringsDeps.HasSuffix(file, ".go") {
			continue
		}
		content, err := sandbox.Deps.IoDeps.ReadFile(file)
		if err != nil {
			return nil, err
		}
		sources = append(sources, apishape.Source{Name: baseName(sandbox, file), Content: string(content)})
	}

	if len(sources) == 0 {
		return nil, sandbox.Deps.StdDeps.Errorf("%s is empty: there is no contract to copy", api_dir)
	}

	// A mechanic's surface — the api files aliasing an opinionated lib — is the
	// lib's, not the repo's, so it never crosses into a consumer.
	shape, err := apishape.New(sandbox, sources)
	if err != nil {
		return nil, err
	}
	return apishape.WithoutMechanic(shape), nil
}

// CopyRemoteApi writes the remote contract into sandbox/deps/<name>/, one file
// per file (RemoteCopyPath), with nothing changed but the package clause, the
// dependency wiring and the generated header — which is what makes the copy legible as the thing it is, doc
// comments and all.
func CopyRemoteApi(sandbox *api.Sandbox, io *stagedfs.StagedFS, remote *apishape.Api, name string) error {
	for _, file := range remote.Files {
		formatted, err := RenderRemoteFile(sandbox, remote, file, name)
		if err != nil {
			return sandbox.Deps.StdDeps.Errorf("%s could not be formatted after the package clause was rewritten: %w", file.Name, err)
		}

		if err := utils.WriteGenerated(sandbox, io, RemoteCopyPath(sandbox, name, file.Name), []byte(formatted)); err != nil {
			return err
		}
	}

	return nil
}

// RenderRemoteFile is the single rendering of one remote contract file into a
// consumer's sandbox/deps/<name>/: the package clause rewritten, the dependency
// wiring stripped, the generated header put on top, then formatted. The install writes what it returns and the
// drift check compares against what it returns, so the two can never disagree
// about what a copy of that file looks like — a check that rendered the file
// its own way would report every copy as drifted the moment stripDepsWiring
// had anything to remove.
func RenderRemoteFile(sandbox *api.Sandbox, remote *apishape.Api, file apishape.ParsedFile, name string) (string, error) {
	content := sandbox.Deps.StringsDeps.ReplaceAll(file.Content,
		"package "+file.Parsed.Package+"\n", "package "+name+"\n")

	content = stripDepsWiring(sandbox, content)
	content = stripMechanicFields(sandbox, remote, content)

	return sandbox.Deps.GoimportsDeps.Format(utils.WithGeneratedHeader(sandbox, content))
}

// RemoteCopyPath is where one file of a remote contract is copied to:
// sandbox/deps/<name>/generated.<file>, since set-dep rewrites it whole.
func RemoteCopyPath(sandbox *api.Sandbox, name string, file string) string {
	return utils.ContractsDir + "/" + name + "/" + utils.GeneratedFile(sandbox, file)
}

// splitModuleSpec cuts "<module>@<version>" apart. A spec with no version
// leaves the version empty, which resolves to whatever the build list already
// holds, or to latest.
func splitModuleSpec(sandbox *api.Sandbox, spec string) (string, string) {
	at := sandbox.Deps.StringsDeps.LastIndex(spec, "@")
	if at < 0 {
		return spec, ""
	}
	return spec[:at], spec[at+1:]
}

// defaultDepName is the name a remote dep takes when --as names none: the last
// segment of the module path, lower-cased, because a dep is a directory and a
// Go package.
func defaultDepName(sandbox *api.Sandbox, module string) string {
	segments := sandbox.Deps.StringsDeps.Split(module, "/")
	return sandbox.Deps.StringsDeps.ToLower(segments[len(segments)-1])
}

// baseName is the last segment of a host path.
func baseName(sandbox *api.Sandbox, path string) string {
	segments := sandbox.Deps.StringsDeps.Split(sandbox.Deps.StringsDeps.ReplaceAll(path, "\\", "/"), "/")
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
// sandbox/api/generated.sandbox.go is generated, so the field is always its own line and
// the import always its own. An import block left holding nothing goes with
// them, since gofmt keeps an empty one.
func stripDepsWiring(sandbox *api.Sandbox, content string) string {
	var kept []string

	for _, line := range sandbox.Deps.StringsDeps.Split(content, "\n") {
		trimmed := sandbox.Deps.StringsDeps.TrimSpace(line)

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

// stripMechanicFields removes every field typed with a mechanic's surface —
// CliSandbox embedded in Sandbox — together with its doc comment: the type it
// names stayed behind with the mechanic files (apishape.WithoutMechanic), so
// the copy would otherwise name a type it does not declare. Like the Deps
// field, each one is its own line of a generated, gofmt'ed struct.
func stripMechanicFields(sandbox *api.Sandbox, remote *apishape.Api, content string) string {
	var kept []string

	for _, line := range sandbox.Deps.StringsDeps.Split(content, "\n") {
		parts := sandbox.Deps.StringsDeps.Fields(line)
		embedded := len(parts) == 1 && apishape.IsMechanic(remote, parts[0])
		named := len(parts) == 2 && apishape.IsMechanic(remote, parts[1])
		if embedded || named {
			kept = dropTrailingComment(sandbox, kept)
			continue
		}
		kept = append(kept, line)
	}

	return sandbox.Deps.StringsDeps.Join(kept, "\n")
}

// isDepsImportLine reports whether one trimmed line imports the loose
// sandbox/deps package, aliased or not. A contract package under
// sandbox/deps/<x>/ ends in another segment, so it does not match.
func isDepsImportLine(sandbox *api.Sandbox, trimmed string) bool {
	if !sandbox.Deps.StringsDeps.HasSuffix(trimmed, "/sandbox/deps\"") {
		return false
	}
	return sandbox.Deps.StringsDeps.HasPrefix(trimmed, "\"") ||
		sandbox.Deps.StringsDeps.HasPrefix(trimmed, "deps \"")
}

// isDepsFieldLine reports whether one trimmed line declares the Sandbox.Deps
// field. It reads the line as fields rather than as text because the file is
// gofmt'ed, so the name and the type are padded apart by however wide the
// widest field of the struct is.
func isDepsFieldLine(sandbox *api.Sandbox, trimmed string) bool {
	parts := sandbox.Deps.StringsDeps.Fields(trimmed)
	return len(parts) == 2 && parts[0] == apishape.DepsField && parts[1] == "*deps.Deps"
}

// dropTrailingComment removes the doc comment the dropped field was carrying,
// which is every line back to the first that is not one. A comment explaining
// Deps outlives it otherwise, describing a field the copy does not have.
func dropTrailingComment(sandbox *api.Sandbox, kept []string) []string {
	for len(kept) > 0 {
		last := sandbox.Deps.StringsDeps.TrimSpace(kept[len(kept)-1])
		if !sandbox.Deps.StringsDeps.HasPrefix(last, "//") {
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
		if sandbox.Deps.StringsDeps.TrimSpace(lines[index]) != "import (" {
			kept = append(kept, lines[index])
			continue
		}

		end := index
		for end < len(lines) && sandbox.Deps.StringsDeps.TrimSpace(lines[end]) != ")" {
			end++
		}

		empty := end < len(lines)
		for inner := index + 1; inner < end; inner++ {
			if sandbox.Deps.StringsDeps.TrimSpace(lines[inner]) != "" {
				empty = false
			}
		}

		if !empty {
			kept = append(kept, lines[index])
			continue
		}

		index = end
	}

	return sandbox.Deps.StringsDeps.Join(kept, "\n")
}
