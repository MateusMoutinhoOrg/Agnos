package verify

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/themesconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// CheckDocs enforces the documentation tree: every doc directory declares
// itself in a props.yaml, every theme it names exists in themes.yaml, every
// first-level doc belongs to at least one theme, and no sub-doc claims a theme
// (a sub-doc is listed only by its parent's Index.md). A theme no doc names is
// not a violation: it simply renders no section in README.md's index.
//
// A project with no docs/ directory has nothing to check.
func CheckDocs(deps *deps.Deps, io *smartio.SmartIO) []string {
	var violations []string

	if !io.IsDir(utils.DocsDir) {
		return violations
	}

	themes_conf, err := utils.LoadThemesConf(deps, io)
	if err != nil {
		return append(violations, err.Error())
	}

	docs, err := utils.CollectDocTree(deps, io)
	if err != nil {
		return append(violations, err.Error())
	}

	for _, doc := range docs {
		if len(doc.Themes) == 0 {
			violations = append(violations, doc.Path+"/props.yaml declares no themes"+
				" (a first-level doc must name at least one theme of themes.yaml)")
		}

		for _, theme := range doc.Themes {
			if !hasTheme(themes_conf.Themes, theme) {
				violations = append(violations, doc.Path+"/props.yaml names unknown theme "+
					theme+" (it is not declared in themes.yaml)")
				continue
			}
		}

		violations = append(violations, checkSubdocs(doc.Subdocs)...)
	}

	return violations
}

// checkSubdocs enforces, at every depth, that a sub-doc carries no themes: it
// is listed by its parent's Index.md, never by README.md's index.
func checkSubdocs(docs []utils.Doc) []string {
	var violations []string

	for _, doc := range docs {
		if len(doc.Themes) > 0 {
			violations = append(violations, doc.Path+"/props.yaml declares themes"+
				" (themes belong to first-level docs only)")
		}
		violations = append(violations, checkSubdocs(doc.Subdocs)...)
	}

	return violations
}

// hasTheme reports whether themes.yaml declares the theme id.
func hasTheme(themes []themesconf.Theme, id string) bool {
	for _, theme := range themes {
		if theme.Id == id {
			return true
		}
	}
	return false
}
