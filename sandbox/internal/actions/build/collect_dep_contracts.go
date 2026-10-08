package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// CollectDepsApi parses every file of every sandbox/deps sub-contract through
// the Go parser dep and returns one rich data map per contract directory, for
// the {{range .DepsApi}} loop in the generated docs/PublicApi/doc.md. A
// contract may be split over several files; each one becomes an entry of the
// directory's Files list, so the doc keeps the source's own grouping.
func CollectDepsApi(sandbox *api.Sandbox, io *smartio.SmartIO) ([]map[string]any, error) {

	var contracts []map[string]any
	for _, dir := range collectLibDirs(sandbox, io, "sandbox/deps") {

		var files []map[string]any
		var symbols []string
		for _, file := range goFilesOf(sandbox, io, "sandbox/deps/"+dir["Name"]) {
			parsed, err := parseGoFile(sandbox, io, file)
			if err != nil {
				return nil, err
			}
			data := fileData(sandbox, file, parsed)
			files = append(files, data)
			symbols = append(symbols, declaredSymbols(data)...)
		}

		contracts = append(contracts, map[string]any{
			"Name":    dir["Name"],
			"Title":   dir["Title"],
			"Files":   files,
			"Page":    DepsApiPageOf(dir["Name"]),
			"Symbols": identifierList(sandbox, symbols),
		})
	}

	return contracts, nil
}
