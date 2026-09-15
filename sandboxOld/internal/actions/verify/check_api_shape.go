package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/apishape"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// apiDir is the contract half of a repo: the package a consumer copies into
// its own sandbox/deps/ when it installs this repo as a dep.
const apiDir = "sandbox/api"

// CheckApiShape enforces that sandbox/api/ stays a shape a consumer can copy
// and convert. Every agnos repo is installable by construction, so the check
// runs everywhere and is not opt-in: a violation caught here is caught by the
// author, and one left for the consumer's install is caught by whoever did not
// write it. The rule itself is in sandbox/internal/apishape.
func CheckApiShape(sandbox *api.Sandbox, io *smartio.SmartIO) []string {
	if !io.IsDir(apiDir) {
		return nil
	}

	var sources []apishape.Source
	for _, file := range goFilesUnder(sandbox, io, apiDir) {
		content, err := io.ReadFile(file)
		if err != nil {
			return []string{file + " could not be read"}
		}
		sources = append(sources, apishape.Source{Name: file, Content: string(content)})
	}

	api, err := apishape.New(sandbox, sources)
	if err != nil {
		return []string{err.Error()}
	}

	return apishape.Violations(sandbox, api)
}
