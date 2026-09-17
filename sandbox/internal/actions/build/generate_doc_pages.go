package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// docPageExt is what every page beside a doc.md is named with.
const docPageExt = ".md"

// removeStaleDocPages drops every page of a doc directory that the build which
// writes them did not produce, so a unit that is gone leaves no page nothing
// links to. The doc's own doc.md and the Index.md a sub-doc would give it are
// never pages, so they are kept.
//
// written holds the project-relative path of every page just rendered, which is
// how a page this build produced is told from one an earlier build left.
func removeStaleDocPages(sandbox *api.Sandbox, io *smartio.SmartIO, dir string, written map[string]bool) {
	for _, path := range io.ListFiles(dir) {
		file := lastSegmentOf(sandbox, path)

		if file == utils.DocFile || file == utils.DocIndexFile {
			continue
		}
		if !sandbox.Deps.Stringsdeps.HasSuffix(file, docPageExt) {
			continue
		}
		if written[path] {
			continue
		}

		io.RemoveDir(path)
	}
}
