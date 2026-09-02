package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos-Cli/sandbox/deps"
)

// GetImports parses the given Go source file content and returns every import
// path it declares, in source order. Unparsable content yields an empty slice.
func GetImports(deps *deps.Deps, content string) []string {

	imports, err := deps.Goimportsdeps.GetImports(content)
	if err != nil {
		return nil
	}

	return imports
}
