package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// adaptersAllowedDirs is the fixed set of sub-directories the adapters/ tree
// may contain.
var adaptersAllowedDirs = []string{"bindings", "impls"}

// adapterLibsDir holds one package per contract, each exporting the binder
// adapters/bindings/standard/new.go is generated to call.
const adapterLibsDir = utils.AdaptersDir

// binderParamType is the single parameter every adapter binder takes. The
// generated standard adapter calls <lib>.Bind(&deps) for every directory it
// finds, so a package that spells this differently only breaks at compile
// time, after build has already rewritten the tree.
const binderParamType = "*deps.Deps"

// CheckAdapters enforces the adapters/ tree: it holds only the bindings and
// libs directories, every lib exports Bind(sandbox *api.Sandbox), and every field
// of Deps is filled by one of them — an unassigned field is a nil func that
// panics on first use, which no compiler catches.
func CheckAdapters(sandbox *api.Sandbox, io *stagedfs.StagedFS) []string {
	var violations []string

	if !io.IsDir("adapters") {
		return violations
	}

	for _, dir := range io.ListDirs("adapters") {
		name := lastSegment(sandbox, dir)
		if !contains(adaptersAllowedDirs, name) {
			violations = append(violations, "adapters/ contains unexpected directory "+name+
				" (allowed: bindings, libs)")
		}
	}

	for _, file := range io.ListFiles("adapters") {
		violations = append(violations, "adapters/ contains unexpected file "+lastSegment(sandbox, file)+
			" (adapters/ may hold only the bindings and libs directories)")
	}

	violations = append(violations, checkAdapterBinders(sandbox, io)...)
	violations = append(violations, checkAdapterCoverage(sandbox, io)...)

	return violations
}

// checkAdapterBinders reports every adapters/impls package that does not export
// the binder the generated standard adapter calls.
func checkAdapterBinders(sandbox *api.Sandbox, io *stagedfs.StagedFS) []string {
	var violations []string

	if !io.IsDir(adapterLibsDir) {
		return violations
	}

	for _, lib := range io.ListDirs(adapterLibsDir) {
		if hasBinder(sandbox, io, lib) {
			continue
		}
		violations = append(violations, lib+" exports no Bind(deps "+binderParamType+")"+
			" (adapters/bindings/standard/new.go is generated to call it)")
	}

	return violations
}

// hasBinder reports whether any .go file of lib declares the binder. A file
// that does not parse is reported by CheckContractDocs only for the contract
// trees, so here an unparsable file simply carries no binder.
func hasBinder(sandbox *api.Sandbox, io *stagedfs.StagedFS, lib string) bool {
	for _, file := range goFilesUnder(sandbox, io, lib) {
		content, err := io.ReadFile(file)
		if err != nil {
			continue
		}

		parsed, err := sandbox.Deps.GoimportsDeps.Parse(string(content))
		if err != nil {
			continue
		}

		for _, function := range parsed.Functions {
			if isBinder(function) {
				return true
			}
		}
	}

	return false
}

// isBinder reports whether one parsed declaration is the adapter binder:
// a plain exported func named Bind taking a single *deps.Deps.
func isBinder(function goimportsdeps.Function) bool {
	return function.Name == "Bind" &&
		function.Receiver == "" &&
		len(function.Params) == 1 &&
		function.Params[0].Type == binderParamType
}

// checkAdapterCoverage enforces the one invariant a binding exists for:
// every field of Deps is filled exactly once. Zero is a nil func that panics
// on first use; two is a silent overwrite in which the last binder wins, which
// no compiler and no test catches. Which adapter fills which field is read
// from the adapter's own declaration, never from the body of its Bind.
//
// A project whose bindings are all hand-written declares no selection, so
// there is nothing to resolve; it falls back to the weaker question — is every
// contract mentioned by some adapter at all.
func checkAdapterCoverage(sandbox *api.Sandbox, io *stagedfs.StagedFS) []string {
	if !io.IsDir(adapterLibsDir) || !io.IsDir(utils.ContractsDir) {
		return nil
	}

	bindings := utils.DeclaredBindings(sandbox, io)
	if len(bindings) == 0 {
		return checkAdapterMentions(sandbox, io)
	}

	var violations []string
	for _, binding := range bindings {
		violations = append(violations, checkBindingCoverage(sandbox, io, binding)...)
	}

	return violations
}

// checkBindingCoverage resolves one binding's selection into the fields it
// fills and reports both ways it can be wrong.
func checkBindingCoverage(sandbox *api.Sandbox, io *stagedfs.StagedFS, binding string) []string {
	var violations []string

	conf, err := utils.LoadBindingConf(sandbox, io, binding)
	if err != nil {
		return []string{err.Error()}
	}

	filled := map[string][]string{}
	for _, adapter := range conf.Adapters {
		adapter_conf, err := utils.LoadAdapterConf(sandbox, io, adapter)
		if err != nil {
			violations = append(violations, utils.BindingConfPath(binding)+" binds "+adapter+
				", which is not installed under "+adapterLibsDir+"/")
			continue
		}
		field := utils.DepField(sandbox, adapter_conf.Dep)
		filled[field] = append(filled[field], adapter)
	}

	for _, contract := range io.ListDirs(utils.ContractsDir) {
		name := lastSegment(sandbox, contract)
		if len(name) == 0 {
			continue
		}
		field := utils.DepField(sandbox, name)

		switch len(filled[field]) {
		case 1:
			continue
		case 0:
			violations = append(violations, "binding "+binding+" fills deps."+field+
				" with no adapter (an unfilled Deps field is a nil func that panics on first use)")
		default:
			violations = append(violations, "binding "+binding+" fills deps."+field+
				" with "+sandbox.Deps.StringsDeps.Join(filled[field], " and ")+
				" (two binders of one field is a silent overwrite: the last one bound wins)")
		}
	}

	return violations
}

// checkAdapterMentions is the coverage question a project with no declared
// binding can still answer: does some adapter mention each contract's Deps
// field at all. The field name is the title-cased contract directory, the same
// spelling sandbox/deps/generated.deps.go is generated with, and a binder fills it either
// whole or field by field — both mention it.
func checkAdapterMentions(sandbox *api.Sandbox, io *stagedfs.StagedFS) []string {
	var violations []string

	bound := adapterSources(sandbox, io)

	for _, contract := range io.ListDirs(utils.ContractsDir) {
		name := lastSegment(sandbox, contract)
		if len(name) == 0 {
			continue
		}
		field := utils.DepField(sandbox, name)

		if sandbox.Deps.StringsDeps.Contains(bound, "deps."+field) {
			continue
		}

		violations = append(violations, "contract "+utils.ContractsDir+"/"+name+
			" has no adapter filling deps."+field+
			" (an unfilled Deps field is a nil func that panics on first use)")
	}

	return violations
}

// adapterSources concatenates every Go source under adapters/impls, so one
// scan answers the question for all contracts at once.
func adapterSources(sandbox *api.Sandbox, io *stagedfs.StagedFS) string {
	sources := ""

	for _, file := range goFilesUnder(sandbox, io, adapterLibsDir) {
		content, err := io.ReadFile(file)
		if err != nil {
			continue
		}
		sources += string(content)
	}

	return sources
}
