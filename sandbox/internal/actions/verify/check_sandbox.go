package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// sandboxAllowedDirs is the fixed set of sub-directories the sandbox/ tree may
// contain; sandboxAllowedFiles is the fixed set of loose files.
var sandboxAllowedDirs = []string{"api", "binds", "deps", "internal"}
var sandboxAllowedFiles = []string{"new.go"}

// CheckSandbox runs every sandbox-layer rule and returns one string per
// violation, in a stable order.
func CheckSandbox(sandbox *api.Sandbox, io *smartio.SmartIO, module string) []string {
	var violations []string

	if !io.IsDir("sandbox") {
		return violations
	}

	violations = append(violations, checkSandboxContents(sandbox, io)...)
	violations = append(violations, checkSandboxImports(sandbox, io, module)...)
	violations = append(violations, checkSandboxApi(sandbox, io, module)...)
	violations = append(violations, checkSandboxDeps(sandbox, io, module)...)
	violations = append(violations, checkSandboxBinds(sandbox, io)...)

	return violations
}

// checkSandboxContents enforces that sandbox/ holds only the api, binds, deps
// and internal directories plus a loose new.go.
func checkSandboxContents(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string

	for _, dir := range io.ListDirs("sandbox") {
		name := lastSegment(sandbox, dir)
		if !contains(sandboxAllowedDirs, name) {
			violations = append(violations, "sandbox/ contains unexpected directory "+name+
				" (allowed: api, binds, deps, internal)")
		}
	}

	for _, file := range io.ListFiles("sandbox") {
		name := lastSegment(sandbox, file)
		if !contains(sandboxAllowedFiles, name) {
			violations = append(violations, "sandbox/ contains unexpected file "+name+
				" (allowed: new.go)")
		}
	}

	return violations
}

// checkSandboxImports enforces that a file under sandbox/ imports nothing but
// another sandbox package. The standard library is banned along with every
// third-party module: a capability the sandbox needs from outside — reading a
// file, formatting a string, sorting a slice — is restated as a contract under
// sandbox/deps/ and filled by an adapter, and that indirection is the whole of
// what keeps the core closed and testable.
//
// sandbox/deps/ is the one exception, being the contracts themselves;
// checkSandboxDeps is what constrains those.
func checkSandboxImports(sandbox *api.Sandbox, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(sandbox, io, "sandbox") {
		if isUnder(file, "sandbox/deps") {
			continue
		}
		for _, imp := range fileImports(sandbox, io, file) {
			if isUnder(imp, module+"/sandbox") {
				continue
			}
			violations = append(violations, file+" imports "+imp+
				"; sandbox/ may import only sandbox packages"+
				" (restate an outside capability as a contract under sandbox/deps/)")
		}
	}

	return violations
}

// checkSandboxApi enforces that sandbox/api/* imports nothing but the loose
// sandbox/deps package — no stdlib, no external module, no other sandbox
// package. api is pure contract: every type it declares is written in Go's own
// builtin types, so any other import is a contract leaking a shape it does not
// own.
//
// sandbox/deps is the one exception, and it buys exactly one field:
// Sandbox.Deps, the capability set every function of the sandbox is handed
// along with the api itself. That field is this repo's own wiring rather than
// part of what it offers, so apishape.DepsField keeps it out of the copy a
// consumer installs — which is what leaves the contract as portable as it was.
func checkSandboxApi(sandbox *api.Sandbox, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(sandbox, io, "sandbox/api") {
		for _, imp := range fileImports(sandbox, io, file) {
			if imp == module+"/sandbox/deps" {
				continue
			}
			violations = append(violations,
				file+" imports "+imp+"; sandbox/api/* may import nothing but "+module+"/sandbox/deps")
		}
	}

	return violations
}

// checkSandboxDeps enforces that a contract package under sandbox/deps/<x>/
// imports nothing at all — the stdlib included. A contract restates an
// outside api, so it may not name a type it borrows from one: a shape the
// sandbox needs is spelled out in Go's own builtin types, and the adapter,
// which lives outside the sandbox, is what converts.
//
// The loose files directly in sandbox/deps/ are the one exception. deps.go
// composes the contracts into deps.Deps, so it names them and nothing else.
func checkSandboxDeps(sandbox *api.Sandbox, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(sandbox, io, "sandbox/deps") {
		for _, imp := range fileImports(sandbox, io, file) {
			if isDirectChild(sandbox, file, "sandbox/deps") {
				if isUnder(imp, module+"/sandbox/deps") {
					continue
				}
				violations = append(violations,
					file+" imports "+imp+"; sandbox/deps/*.go may import only sandbox/deps packages")
				continue
			}
			violations = append(violations,
				file+" imports "+imp+"; sandbox/deps/<x>/ may import nothing at all")
		}
	}

	return violations
}

// checkSandboxBinds enforces that every file in sandbox/binds mirrors a file
// of the same name in sandbox/api and declares only functions.
func checkSandboxBinds(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	var violations []string

	apiFiles := map[string]bool{}
	for _, file := range io.ListFiles("sandbox/api") {
		apiFiles[lastSegment(sandbox, file)] = true
	}

	for _, file := range io.ListFiles("sandbox/binds") {
		name := lastSegment(sandbox, file)
		if !sandbox.Deps.Stringsdeps.HasSuffix(name, ".go") {
			continue
		}
		if !apiFiles[name] {
			violations = append(violations,
				"sandbox/binds/"+name+" has no matching file in sandbox/api/")
		}
		for _, decl := range topLevelNonFuncDecls(sandbox, io, file) {
			violations = append(violations,
				"sandbox/binds/"+name+" declares "+decl+"; sandbox/binds/ may contain only functions")
		}
	}

	return violations
}

// topLevelNonFuncDecls returns a label for every top-level declaration in the
// file that is not a function, types first, then constants, then variables.
func topLevelNonFuncDecls(sandbox *api.Sandbox, io *smartio.SmartIO, file string) []string {
	parsed := parseFile(sandbox, io, file)
	if parsed == nil {
		return nil
	}

	var found []string
	for range parsed.Types {
		found = append(found, "a type")
	}
	for range parsed.Constants {
		found = append(found, "a const")
	}
	for range parsed.Variables {
		found = append(found, "a var")
	}
	return found
}

// fileImports returns the import paths of one Go file, sorted.
func fileImports(sandbox *api.Sandbox, io *smartio.SmartIO, file string) []string {
	parsed := parseFile(sandbox, io, file)
	if parsed == nil {
		return nil
	}

	var imports []string
	for _, spec := range parsed.Imports {
		imports = append(imports, spec.Path)
	}
	sandbox.Deps.Sortdeps.Strings(imports)
	return imports
}

// parseFile reads file through the transaction-aware io and parses it. A file
// that cannot be read or parsed yields nil (the compiler reports those).
func parseFile(sandbox *api.Sandbox, io *smartio.SmartIO, file string) *goimportsdeps.File {
	content, err := io.ReadFile(file)
	if err != nil {
		return nil
	}
	parsed, err := sandbox.Deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return nil
	}
	return parsed
}

// goFilesUnder lists every .go file at or below dir, in listing order.
func goFilesUnder(sandbox *api.Sandbox, io *smartio.SmartIO, dir string) []string {
	var files []string
	for _, file := range io.ListFilesRecursively(dir) {
		if sandbox.Deps.Stringsdeps.HasSuffix(file, ".go") {
			files = append(files, file)
		}
	}
	return files
}

func lastSegment(sandbox *api.Sandbox, path string) string {
	parts := sandbox.Deps.Stringsdeps.Split(path, "/")
	return parts[len(parts)-1]
}

func contains(list []string, want string) bool {
	for _, item := range list {
		if item == want {
			return true
		}
	}
	return false
}

// isUnder reports whether imp is prefix itself or a package below it.
func isUnder(imp string, prefix string) bool {
	return imp == prefix || len(imp) > len(prefix) && imp[:len(prefix)+1] == prefix+"/"
}

// isDirectChild reports whether path is a loose entry of dir rather than of a
// directory below it.
func isDirectChild(sandbox *api.Sandbox, path string, dir string) bool {
	return isUnder(path, dir) &&
		len(sandbox.Deps.Stringsdeps.Split(path, "/")) == len(sandbox.Deps.Stringsdeps.Split(dir, "/"))+1
}
