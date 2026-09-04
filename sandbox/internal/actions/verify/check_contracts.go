package verify

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// contractDirs are the trees docs/PublicApi is generated from: sandbox/api is
// the surface sandbox.New returns, sandbox/deps the contracts a caller fills
// or replaces.
var contractDirs = []string{"sandbox/api", "sandbox/deps"}

// CheckContracts runs the rules docs/PublicApi depends on and returns one
// string per violation, in a stable order. The page is rendered from these
// files by the Go parser dep, so every one of them has to parse, and every
// exported declaration has to carry the doc comment that becomes its
// description.
func CheckContracts(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	for _, dir := range contractDirs {
		if !io.IsDir(dir) {
			continue
		}
		for _, file := range goFilesUnder(io, dir) {
			violations = append(violations, checkContractFile(deps, io, file)...)
		}
	}

	return violations
}

// checkContractFile parses one contract file and reports every exported
// top-level declaration it finds without a doc comment. Struct fields and
// interface methods are not required to be commented: their comment fills a
// description column when it is there, and the type's own comment covers them
// when it is not.
func checkContractFile(deps *deps.Deps, io *smartio.SmartIO, file string) []string {
	content, err := io.ReadFile(file)
	if err != nil {
		return []string{file + " could not be read"}
	}

	parsed, err := deps.Goimportsdeps.Parse(string(content))
	if err != nil {
		return []string{file + " is not parsable Go: " + err.Error()}
	}

	var violations []string

	for _, entry := range parsed.Types {
		if entry.Exported && strings.TrimSpace(entry.Doc) == "" {
			violations = append(violations, undocumented(file, "type", entry.Name))
		}
	}
	for _, entry := range parsed.Functions {
		if entry.Exported && strings.TrimSpace(entry.Doc) == "" {
			violations = append(violations, undocumented(file, "function", entry.Name))
		}
	}
	for _, entry := range parsed.Constants {
		if entry.Exported && strings.TrimSpace(entry.Doc) == "" {
			violations = append(violations, undocumented(file, "const", entry.Name))
		}
	}
	for _, entry := range parsed.Variables {
		if entry.Exported && strings.TrimSpace(entry.Doc) == "" {
			violations = append(violations, undocumented(file, "var", entry.Name))
		}
	}

	return violations
}

// undocumented phrases the one violation this file reports.
func undocumented(file string, kind string, name string) string {
	return file + " declares " + kind + " " + name +
		" without a doc comment; docs/PublicApi is generated from those comments"
}
