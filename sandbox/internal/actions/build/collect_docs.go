package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectDocs walks docs/ and returns its first-level docs, each carrying its
// sub-docs recursively, sorted the way every index lists them. It feeds
// GenerateDocIndexes, which writes docs/Index/<theme-id>.md and one Index.md
// per doc that has sub-docs. A project with no docs/ directory yields an empty
// slice.
func CollectDocs(deps *deps.Deps, io *smartio.SmartIO) ([]utils.Doc, error) {
	return utils.CollectDocTree(deps, io)
}
