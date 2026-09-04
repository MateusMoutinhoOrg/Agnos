package build

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/themesconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// GenerateDocIndexes writes every generated page of the documentation tree:
// one docs/Index/<theme-id>.md per theme of themes.yaml, listing the
// first-level docs that name that theme, and one <doc>/Index.md per doc that
// has sub-docs, listing its direct sub-docs.
//
// docs/Index/ belongs to the generator end to end: it is removed and rewritten
// on every build, so an index of a theme that no longer exists never survives.
// A project with no docs/ directory generates nothing.
func GenerateDocIndexes(deps *deps.Deps, io *smartio.SmartIO, docs []utils.Doc, themes []themesconf.Theme) error {
	if !io.IsDir(utils.DocsDir) {
		return nil
	}

	index_dir := utils.DocsDir + "/" + utils.DocsIndexName
	io.RemoveDir(index_dir)
	io.CreateDir(index_dir)

	for _, theme := range themes {
		vars := map[string]any{
			"Name":        theme.Name,
			"Description": theme.Description,
			"Docs":        docRows(index_dir, docsOfTheme(docs, theme.Id)),
		}

		dest := index_dir + "/" + theme.Id + ".md"
		if err := utils.RenderTemplateToDest(deps, io, "templates/theme_index.md", vars, dest); err != nil {
			return err
		}
	}

	return generateSubdocIndexes(deps, io, docs)
}

// generateSubdocIndexes walks the tree and writes <doc>/Index.md for every doc
// that has sub-docs, at any depth.
func generateSubdocIndexes(deps *deps.Deps, io *smartio.SmartIO, docs []utils.Doc) error {
	for _, doc := range docs {
		if len(doc.Subdocs) == 0 {
			continue
		}

		vars := map[string]any{
			"Name":        doc.Name,
			"Description": doc.Description,
			"Docs":        docRows(doc.Path, doc.Subdocs),
		}

		dest := doc.Path + "/" + utils.DocIndexFile
		if err := utils.RenderTemplateToDest(deps, io, "templates/doc_index.md", vars, dest); err != nil {
			return err
		}

		if err := generateSubdocIndexes(deps, io, doc.Subdocs); err != nil {
			return err
		}
	}
	return nil
}

// docsOfTheme keeps the docs whose props.yaml names the theme id.
func docsOfTheme(docs []utils.Doc, id string) []utils.Doc {
	var kept []utils.Doc
	for _, doc := range docs {
		for _, theme := range doc.Themes {
			if theme == id {
				kept = append(kept, doc)
				break
			}
		}
	}
	return kept
}

// docRows turns docs into the rows an index template ranges over: the name,
// the description, and the link to the doc's own doc.md, written relative to
// the directory the index file itself lives in.
func docRows(from_dir string, docs []utils.Doc) []map[string]any {
	rows := make([]map[string]any, 0, len(docs))
	for _, doc := range docs {
		rows = append(rows, map[string]any{
			"Name":        doc.Name,
			"Description": doc.Description,
			"Link":        relativeLink(from_dir, doc.Path+"/"+utils.DocFile),
		})
	}
	return rows
}

// relativeLink writes target as a path relative to from_dir. Both are
// project-relative ("docs/Index", "docs/Install/doc.md"), and the result is
// what a Markdown link must hold: GitHub resolves a "/"-prefixed link against
// the site root rather than the repository, so no generated link starts with
// one.
func relativeLink(from_dir string, target string) string {
	from := splitPath(from_dir)
	to := splitPath(target)

	common := 0
	for common < len(from) && common < len(to)-1 && from[common] == to[common] {
		common++
	}

	parts := make([]string, 0, len(from)-common+len(to)-common)
	for range from[common:] {
		parts = append(parts, "..")
	}
	parts = append(parts, to[common:]...)

	return strings.Join(parts, "/")
}

// splitPath breaks a project-relative path into its segments, dropping the
// empty ones a leading, trailing or doubled "/" would produce.
func splitPath(path string) []string {
	parts := make([]string, 0, strings.Count(path, "/")+1)
	for _, part := range strings.Split(path, "/") {
		if part != "" && part != "." {
			parts = append(parts, part)
		}
	}
	return parts
}
