package add_doc

import (
	"strings"

	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/docpropsconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/parsables/themesconf"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// AddDocInternal writes the two files of a new doc — doc.md and props.yaml —
// after enforcing what `verify` would later reject: an unknown theme, a
// first-level doc with no theme, a sub-doc with one, and a sub-doc whose
// parent is not a doc. It refuses to overwrite an existing doc.
func AddDocInternal(deps *deps.Deps, io *smartio.SmartIO, props api.DocProps) error {
	if err := utils.ValidateDocName(deps, props.Name); err != nil {
		return err
	}

	description := strings.TrimSpace(props.Description)
	if description == "" {
		return deps.Std.Errorf("add-doc requires --description")
	}

	segments := utils.DocSegments(props.Name)
	dir := utils.DocDir(props.Name)
	parent := utils.DocParentDir(props.Name)

	if io.IsDir(dir) {
		return deps.Std.Errorf("doc %s already exists", dir)
	}

	themes, err := checkThemes(deps, io, props.Themes, len(segments) == 1)
	if err != nil {
		return err
	}

	if len(segments) > 1 && !io.IsFile(parent+"/"+utils.DocPropsFile) {
		return deps.Std.Errorf("%s is not a doc: create it first, a sub-doc lives inside one", parent)
	}

	deps.Std.Log("add-doc creating %s \n", dir)

	conf := docpropsconf.NewEmpty(deps)
	conf.Name = utils.DocTitle(props.Name)
	conf.Description = description
	for _, theme := range themes {
		conf.AddTheme(theme)
	}

	if err := io.WriteFile(dir+"/"+utils.DocPropsFile, []byte(conf.Render())); err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Name":        conf.Name,
		"Description": conf.Description,
	}
	doc, err := deps.Embeddeps.RenderTemplate("templates/doc_doc.md", vars)
	if err != nil {
		return err
	}
	return io.WriteFile(dir+"/"+utils.DocFile, doc)
}

// checkThemes normalizes the requested theme ids and enforces where they may
// appear: a first-level doc names at least one theme declared in themes.yaml,
// a sub-doc names none — it is listed by its parent's Index.md.
func checkThemes(deps *deps.Deps, io *smartio.SmartIO, requested []string, first_level bool) ([]string, error) {
	var themes []string
	for _, theme := range requested {
		theme = strings.TrimSpace(theme)
		if theme != "" {
			themes = append(themes, theme)
		}
	}

	if !first_level {
		if len(themes) > 0 {
			return nil, deps.Std.Errorf("--theme belongs to a first-level doc only: a sub-doc is listed by its parent's %s", utils.DocIndexFile)
		}
		return themes, nil
	}

	themes_conf, err := utils.LoadThemesConf(deps, io)
	if err != nil {
		return nil, err
	}

	if len(themes) == 0 {
		return nil, deps.Std.Errorf("add-doc requires at least one --theme for a first-level doc (declared in themes.yaml: %s)",
			strings.Join(themeIds(themes_conf.Themes), ", "))
	}

	for _, theme := range themes {
		if !hasTheme(themes_conf.Themes, theme) {
			return nil, deps.Std.Errorf("unknown theme %q: themes.yaml declares %s",
				theme, strings.Join(themeIds(themes_conf.Themes), ", "))
		}
	}

	return themes, nil
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

// themeIds lists the declared theme ids, for the error message naming what a
// doc may choose from.
func themeIds(themes []themesconf.Theme) []string {
	ids := make([]string, 0, len(themes))
	for _, theme := range themes {
		ids = append(ids, theme.Id)
	}
	return ids
}
