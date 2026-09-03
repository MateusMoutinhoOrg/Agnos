package verify

import (
	"go/ast"
	"go/parser"
	"go/token"
	"sort"
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
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

	violations = append(violations, checkSandboxContents(io)...)
	violations = append(violations, checkSandboxImports(deps, io, module)...)
	violations = append(violations, checkSandboxApi(deps, io, module)...)
	violations = append(violations, checkSandboxDeps(deps, io, module)...)
	violations = append(violations, checkSandboxBinds(deps, io)...)

	return violations
}

// checkSandboxContents enforces that sandbox/ holds only the api, binds, deps
// and internal directories plus a loose new.go.
func checkSandboxContents(io *smartio.SmartIO) []string {
	var violations []string

	for _, dir := range io.ListDirs("sandbox") {
		name := lastSegment(dir)
		if !contains(sandboxAllowedDirs, name) {
			violations = append(violations, "sandbox/ contains unexpected directory "+name+
				" (allowed: api, binds, deps, internal)")
		}
	}

	for _, file := range io.ListFiles("sandbox") {
		name := lastSegment(file)
		if !contains(sandboxAllowedFiles, name) {
			violations = append(violations, "sandbox/ contains unexpected file "+name+
				" (allowed: new.go)")
		}
	}

	return violations
}

// checkSandboxImports enforces that no file under sandbox/ imports a
// module-internal package that lives outside sandbox/.
func checkSandboxImports(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(io, "sandbox") {
		for _, imp := range fileImports(deps, io, file) {
			if isModuleInternal(imp, module) && !isUnder(imp, module+"/sandbox") {
				violations = append(violations,
					file+" imports "+imp+" which is outside sandbox/")
			}
		}
	}

	return violations
}

// checkSandboxApi enforces that sandbox/api/* imports nothing at all beyond
// other sandbox/api packages (no stdlib, no external modules, and in
// particular not sandbox/deps — api is pure contract).
func checkSandboxApi(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(io, "sandbox/api") {
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
// standard library and other sandbox/deps packages.
func checkSandboxDeps(deps *deps.Deps, io *smartio.SmartIO, module string) []string {
	var violations []string

	for _, file := range goFilesUnder(io, "sandbox/deps") {
		for _, imp := range fileImports(deps, io, file) {
			if isStdlib(imp) {
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
		apiFiles[lastSegment(file)] = true
	}

	for _, file := range io.ListFiles("sandbox/binds") {
		name := lastSegment(file)
		if !strings.HasSuffix(name, ".go") {
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
// file that is not a function (imports are ignored).
func topLevelNonFuncDecls(deps *deps.Deps, io *smartio.SmartIO, file string) []string {
	tree := parseFile(deps, io, file)
	if tree == nil {
		return nil
	}

	var found []string
	for _, decl := range tree.Decls {
		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}
		switch gen.Tok {
		case token.IMPORT:
			continue
		case token.TYPE:
			found = append(found, "a type")
		case token.CONST:
			found = append(found, "a const")
		case token.VAR:
			found = append(found, "a var")
		}
	}
	return found
}

// fileImports returns the import paths of one Go file, sorted.
func fileImports(deps *deps.Deps, io *smartio.SmartIO, file string) []string {
	tree := parseFile(deps, io, file)
	if tree == nil {
		return nil
	}

	var imports []string
	for _, spec := range tree.Imports {
		imports = append(imports, strings.Trim(spec.Path.Value, `"`))
	}
	sort.Strings(imports)
	return imports
}

// parseFile reads file through the transaction-aware io and parses it. A file
// that cannot be read or parsed yields nil (the compiler reports those).
func parseFile(deps *deps.Deps, io *smartio.SmartIO, file string) *ast.File {
	content, err := io.ReadFile(file)
	if err != nil {
		return nil
	}
	tree, err := parser.ParseFile(token.NewFileSet(), file, content, parser.SkipObjectResolution)
	if err != nil {
		return nil
	}
	return tree
}

// goFilesUnder lists every .go file at or below dir, in listing order.
func goFilesUnder(io *smartio.SmartIO, dir string) []string {
	var files []string
	for _, file := range io.ListFilesRecursively(dir) {
		if strings.HasSuffix(file, ".go") {
			files = append(files, file)
		}
	}
	return files
}

func lastSegment(path string) string {
	parts := strings.Split(path, "/")
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

// isModuleInternal reports whether imp is a package of this module.
func isModuleInternal(imp string, module string) bool {
	return imp == module || strings.HasPrefix(imp, module+"/")
}

// isUnder reports whether imp is prefix itself or a package below it.
func isUnder(imp string, prefix string) bool {
	return imp == prefix || strings.HasPrefix(imp, prefix+"/")
}

// isStdlib reports whether imp is a standard-library package (its first path
// segment carries no dot, so it is not a domain).
func isStdlib(imp string) bool {
	first := strings.Split(imp, "/")[0]
	return !strings.Contains(first, ".")
}
