package goimportsdeps

import (
	"go/parser"
	"go/token"
	"strconv"

	"{{.Module}}/sandbox/deps"
	goimportsdeps "{{.Module}}/sandbox/deps/goimportsdeps"
)

// Bind fills deps.Deps.Goimportsdeps, providing the capability to parse the
// import paths declared by a Go source file.
func Bind(deps *deps.Deps) {
	deps.Goimportsdeps = goimportsdeps.Lib{
		GetImports: getImports,
	}
}

// getImports fills goimportsdeps.Lib.GetImports, parsing only the import
// section of the given Go source and unquoting each import path.
func getImports(content string) ([]string, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}

	imports := make([]string, 0, len(file.Imports))
	for _, spec := range file.Imports {
		path, err := strconv.Unquote(spec.Path.Value)
		if err != nil {
			return nil, err
		}
		imports = append(imports, path)
	}

	return imports, nil
}
