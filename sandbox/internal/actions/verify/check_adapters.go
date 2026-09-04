package verify

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	goimportsdeps "github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// adaptersAllowedDirs is the fixed set of sub-directories the adapters/ tree
// may contain.
var adaptersAllowedDirs = []string{"availables", "libs"}

// adapterLibsDir holds one package per contract, each exporting the binder
// adapters/availables/standard/new.go is generated to call.
const adapterLibsDir = "adapters/libs"

// binderParamType is the single parameter every adapter binder takes. The
// generated standard adapter calls <lib>.Bind(&deps) for every directory it
// finds, so a package that spells this differently only breaks at compile
// time, after build has already rewritten the tree.
const binderParamType = "*deps.Deps"

// CheckAdapters enforces the adapters/ tree: it holds only the availables and
// libs directories, every lib exports Bind(deps *deps.Deps), and every field
// of Deps is filled by one of them — an unassigned field is a nil func that
// panics on first use, which no compiler catches.
func CheckAdapters(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir("adapters") {
		return violations
	}

	for _, dir := range io.ListDirs("adapters") {
		name := lastSegment(dir)
		if !contains(adaptersAllowedDirs, name) {
			violations = append(violations, "adapters/ contains unexpected directory "+name+
				" (allowed: availables, libs)")
		}
	}

	for _, file := range io.ListFiles("adapters") {
		violations = append(violations, "adapters/ contains unexpected file "+lastSegment(file)+
			" (adapters/ may hold only the availables and libs directories)")
	}

	violations = append(violations, checkAdapterBinders(deps, io)...)
	violations = append(violations, checkAdapterCoverage(io)...)

	return violations
}

// checkAdapterBinders reports every adapters/libs package that does not export
// the binder the generated standard adapter calls.
func checkAdapterBinders(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir(adapterLibsDir) {
		return violations
	}

	for _, lib := range io.ListDirs(adapterLibsDir) {
		if hasBinder(deps, io, lib) {
			continue
		}
		violations = append(violations, lib+" exports no Bind(deps "+binderParamType+")"+
			" (adapters/availables/standard/new.go is generated to call it)")
	}

	return violations
}

// hasBinder reports whether any .go file of lib declares the binder. A file
// that does not parse is reported by CheckContracts only for the contract
// trees, so here an unparsable file simply carries no binder.
func hasBinder(deps *deps.Deps, io *smartio.SmartIO, lib string) bool {
	for _, file := range goFilesUnder(io, lib) {
		content, err := io.ReadFile(file)
		if err != nil {
			continue
		}

		parsed, err := deps.Goimportsdeps.Parse(string(content))
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

// checkAdapterCoverage reports every sandbox/deps contract with no adapter lib
// mentioning its Deps field. The field name is the title-cased contract
// directory, the same spelling sandbox/deps/deps.go is generated with, and a
// binder fills it either whole or field by field — both mention it.
func checkAdapterCoverage(io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir(adapterLibsDir) || !io.IsDir("sandbox/deps") {
		return violations
	}

	bound := adapterSources(io)

	for _, contract := range io.ListDirs("sandbox/deps") {
		name := lastSegment(contract)
		if len(name) == 0 {
			continue
		}
		field := strings.ToUpper(name[:1]) + name[1:]

		if strings.Contains(bound, "deps."+field) {
			continue
		}

		violations = append(violations, "contract sandbox/deps/"+name+
			" has no adapter filling deps."+field+
			" (an unfilled Deps field is a nil func that panics on first use)")
	}

	return violations
}

// adapterSources concatenates every Go source under adapters/libs, so one
// scan answers the question for all contracts at once.
func adapterSources(io *smartio.SmartIO) string {
	var sources strings.Builder

	for _, file := range goFilesUnder(io, adapterLibsDir) {
		content, err := io.ReadFile(file)
		if err != nil {
			continue
		}
		sources.Write(content)
	}

	return sources.String()
}
