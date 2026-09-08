package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/goimportsdeps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// sandboxAllowedDirs is the fixed set of sub-directories the sandbox/ tree may
// contain; sandboxAllowedFiles is the fixed set of loose files.
var sandboxAllowedDirs = []string{"api", "binds", "deps", "internal"}
var sandboxAllowedFiles = []string{"new.go"}

// CheckSandbox runs every sandbox-layer rule and returns one string per
// violation, in a stable order.
func CheckSandbox(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	if !io.IsDir("sandbox") {
		return violations
	}

	violations = append(violations, checkSandboxContents(deps, io)...)
	violations = append(violations, checkSandboxImports(deps, io, module)...)
	violations = append(violations, checkSandboxApi(deps, io, module)...)
	violations = append(violations, checkSandboxDeps(deps, io, module)...)
	violations = append(violations, checkSandboxBinds(deps, io)...)

	return violations
}

// checkSandboxContents enforces that sandbox/ holds only the api, binds, deps
// and internal directories plus a loose new.go.
func checkSandboxContents(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	for _, dir := range io.ListDirs("sandbox") {
		name := lastSegment(deps, dir)
		if !contains(sandboxAllowedDirs, name) {
			violations = append(violations, "sandbox/ contains unexpected directory "+name+
				" (allowed: api, binds, deps, internal)")
		}
	}

	for _, file := range io.ListFiles("sandbox") {
		name := lastSegment(deps, file)
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
func checkSandboxImports(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(deps, io, "sandbox") {
		if isUnder(file, "sandbox/deps") {
			continue
		}
		for _, imp := range fileImports(deps, io, file) {
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

// checkSandboxApi enforces that sandbox/api/* imports nothing at all beyond
// other sandbox/api packages (no stdlib, no external modules, and in
// particular not sandbox/deps — api is pure contract).
func checkSandboxApi(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(deps, io, "sandbox/api") {
		for _, imp := range fileImports(deps, io, file) {
			if !isUnder(imp, module+"/sandbox/api") {
				violations = append(violations,
					file+" imports "+imp+"; sandbox/api/* may import only other sandbox/api packages")
			}
		}
	}

	return violations
}

// checkSandboxDeps enforces that sandbox/deps/* imports nothing but the
// standard library and other sandbox/deps packages. It is the boundary itself:
// a contract restates an outside api, so it names the standard library types
// that api is written in, and nothing more.
func checkSandboxDeps(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(deps, io, "sandbox/deps") {
		for _, imp := range fileImports(deps, io, file) {
			if isStdlib(deps, imp) {
				continue
			}
			if isUnder(imp, module+"/sandbox/deps") {
				continue
			}
			violations = append(violations,
				file+" imports "+imp+"; sandbox/deps/* may import only the standard library and sandbox/deps packages")
		}
	}

	return violations
}

// checkSandboxBinds enforces that every file in sandbox/binds mirrors a file
// of the same name in sandbox/api and declares only functions.
func checkSandboxBinds(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	apiFiles := map[string]bool{}
	for _, file := range io.ListFiles("sandbox/api") {
		apiFiles[lastSegment(deps, file)] = true
	}

	for _, file := range io.ListFiles("sandbox/binds") {
		name := lastSegment(deps, file)
		if !deps.Stringsdeps.HasSuffix(name, ".go") {
			continue
		}
		if !apiFiles[name] {
			violations = append(violations,
				"sandbox/binds/"+name+" has no matching file in sandbox/api/")
		}
		for _, decl := range topLevelNonFuncDecls(deps, io, file) {
			violations = append(violations,
				"sandbox/binds/"+name+" declares "+decl+"; sandbox/binds/ may contain only functions")
		}
	}

	return violations
}

// topLevelNonFuncDecls returns a label for every top-level declaration in the
// file that is not a function, types first, then constants, then variables.
func topLevelNonFuncDecls(deps *deps.Deps, io *smartio.SmartIO, file string) []string {
	parsed := parseFile(deps, io, file)
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
func fileImports(deps *deps.Deps, io *smartio.SmartIO, file string) []string {
	parsed := parseFile(deps, io, file)
	if parsed == nil {
		return nil
	}

	var imports []string
	for _, spec := range parsed.Imports {
		imports = append(imports, spec.Path)
	}
	deps.Sortdeps.Strings(imports)
	return imports
}

// parseFile reads file through the transaction-aware io and parses it. A file
// that cannot be read or parsed yields nil (the compiler reports those).
func parseFile(deps *deps.Deps, io *smartio.SmartIO, file string) *goimportsdeps.File {
	content, err := io.ReadFile(file)
	if err != nil {
		return nil
	}
	parsed, err := deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return nil
	}
	return parsed
}

// goFilesUnder lists every .go file at or below dir, in listing order.
func goFilesUnder(deps *deps.Deps, io *smartio.SmartIO, dir string) []string {
	var files []string
	for _, file := range io.ListFilesRecursively(dir) {
		if deps.Stringsdeps.HasSuffix(file, ".go") {
			files = append(files, file)
		}
	}
	return files
}

func lastSegment(deps *deps.Deps, path string) string {
	parts := deps.Stringsdeps.Split(path, "/")
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

// isStdlib reports whether imp is a standard-library package (its first path
// segment carries no dot, so it is not a domain).
func isStdlib(deps *deps.Deps, imp string) bool {
	first := deps.Stringsdeps.Split(imp, "/")[0]
	return !deps.Stringsdeps.Contains(first, ".")
}
