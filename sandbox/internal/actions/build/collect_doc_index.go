package build

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/themesconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CollectDocIndex groups the doc tree by theme for the documentation index of
// README.md, which is the single entry point to the docs: one section per
// theme of themes.yaml, in that file's order, each listing the first-level
// docs that name it. A theme no doc names renders no section at all, so an
// unused theme costs nothing and is not an error.
//
// Links are written relative to the repository root, where README.md lives.
func CollectDocIndex(docs []utils.Doc, themes []themesconf.Theme) []map[string]any {
	var index []map[string]any

	for _, theme := range themes {
		grouped := docsOfTheme(docs, theme.Id)
		if len(grouped) == 0 {
			continue
		}

		index = append(index, map[string]any{
			"Name":        theme.Name,
			"Description": theme.Description,
			"Docs":        docRows("", grouped),
		})
	}

	return index
}

// MergeDocs adds the docs a build generates to the ones already on disk,
// keeping the disk copy when both describe the same directory, and returns the
// result in index order. The two agree on content — the generated doc is
// rendered from the same template that wrote the file — so either copy will
// do; taking the disk one keeps the tree the source of truth whenever it has
// something to say.
func MergeDocs(on_disk []utils.Doc, generated []utils.Doc) []utils.Doc {
	merged := make([]utils.Doc, 0, len(on_disk)+len(generated))
	merged = append(merged, on_disk...)

	for _, doc := range generated {
		if hasDocDir(merged, doc.Dir) {
			continue
		}
		merged = append(merged, doc)
	}

	utils.SortDocs(merged)
	return merged
}

// hasDocDir reports whether docs already carries the doc directory.
func hasDocDir(docs []utils.Doc, dir string) bool {
	for _, doc := range docs {
		if doc.Dir == dir {
			return true
		}
	}
	return false
}
