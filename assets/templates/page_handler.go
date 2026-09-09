package {{.Package}}

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/pageio"
	"{{.Module}}/sandbox/internal/routeio"
)

// pageAsset is the embedded template this route renders, as deps.Embeddeps
// spells a path: slash-separated and relative to the root of the assets
// package, so "assets/frontend/pages/{{.Identifier}}.html" on disk.
const pageAsset = "frontend/pages/{{.Identifier}}.html"

// pageVars is what the template is rendered against — every {{"{{ .Field }}"}}
// of the html is one exported field here, so a new variable in the page is a
// new field, checked by the compiler rather than discovered at request time.
type pageVars struct {
	Title   string
	Message string
}

// RouteHandler answers {{.Method}} {{.Trigger}} with the rendered page. It
// renders through pageio.Render rather than deps.Embeddeps.RenderTemplate,
// which is what puts the asset helpers — staticref, dirref, cssref, jsref,
// inline, include — in reach of the page, so the html names a directory
// instead of hardcoding every link it needs.
//
// A template that fails to read or to render is a packaging or authoring
// mistake, never the caller's — a helper pointed at an asset that is not there
// is one of them — so it answers 500 and says nothing about the asset tree.
func RouteHandler(deps *deps.Deps, entries *Entries, response serverdeps.Response) int {
	content, err := pageio.Render(deps, pageAsset, pageVars{
		Title:   {{printf "%q" .Title}},
		Message: "this page is rendered from assets/frontend/pages/{{.Identifier}}.html",
	})
	if err != nil {
		return routeio.WriteError(deps, response, api.StatusFailure, "page",
			"could not render the {{.Identifier}} page")
	}

	response.SetHeader("Content-Type", "text/html; charset=utf-8")
	response.SetStatus(api.StatusOk)
	response.Write(content)

	return api.StatusOk
}
