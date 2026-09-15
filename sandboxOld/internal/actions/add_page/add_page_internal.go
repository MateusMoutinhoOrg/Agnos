package add_page

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addRouteAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_route"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// pageioDir is the front layer's shared package, and the same directory
// `build` reads hasFront from: without it there is nothing for the scaffolded
// handler to render through.
const pageioDir = "sandbox/internal/pageio"

// pageMethod and pageCategory are fixed: a page is read by a browser
// navigating to it, and it is listed in docs/Routes under a heading of its own
// so the pages of a project read as a set rather than scattered among its api
// routes.
const (
	pageMethod   = "GET"
	pageCategory = "Pages"
)

// AddPageInternal writes the two halves of a new page — the route package that
// answers it and the html template it renders — on one open transaction.
//
// The route half goes through AddRouteInternal rather than a second copy of
// the same rendering, so a page's route.yaml is the one add-route writes and
// every route editor keeps working on it. That call is also what refuses an
// existing route: io.WriteFile does not overwrite.
//
// An existing html template, on the other hand, is kept. That is the way back
// from a front-purge, which drops a page's route and leaves its content: this
// is the one file of a page that holds work nobody can regenerate.
func AddPageInternal(sandbox *api.Sandbox, io *smartio.SmartIO, props api.PageProps) error {
	if err := utils.ValidateRouteName(sandbox, props.Name); err != nil {
		return err
	}

	identifier := utils.RouteIdentifier(sandbox, props.Name)

	if identifier == utils.StaticRouteName {
		return sandbox.Deps.Std.Errorf("%q is the route serving the static assets and cannot be a page", utils.StaticRouteName)
	}
	if !io.IsDir(pageioDir) {
		return sandbox.Deps.Std.Errorf("the project has no front layer: run front-init before declaring a page")
	}

	trigger := sandbox.Deps.Stringsdeps.TrimSpace(props.Trigger)
	if trigger == "" {
		trigger = "/" + identifier
	}

	help := sandbox.Deps.Stringsdeps.TrimSpace(props.Help)
	if help == "" {
		help = "Renders the " + identifier + " page from the embedded html template"
	}

	if err := addRouteAction.AddRouteInternal(sandbox, io, props.Name, pageMethod, trigger, help, pageCategory); err != nil {
		return err
	}

	return writePage(sandbox, io, props, identifier, trigger)
}

// writePage overwrites the stub handler add-route just left with one that
// renders the page, and writes the html template beside it.
func writePage(sandbox *api.Sandbox, io *smartio.SmartIO, props api.PageProps, identifier string, trigger string) error {
	segment, err := utils.RouteIdentifierSegment(sandbox, trigger)
	if err != nil {
		return err
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}

	title := sandbox.Deps.Stringsdeps.TrimSpace(props.Title)
	if title == "" {
		title = identifier
	}

	vars := map[string]interface{}{
		"Identifier": identifier,
		"Package":    utils.RoutePackage(sandbox, props.Name),
		"Module":     module_conf.Module,
		"Method":     pageMethod,
		"Trigger":    segment,
		"Title":      title,
	}

	handler := utils.RouteDir(sandbox, props.Name) + "/handler.go"
	if err := utils.RenderTemplateToDest(sandbox, io, "templates/page_handler.go", vars, handler); err != nil {
		return err
	}

	page := utils.PageAsset(sandbox, props.Name)
	if io.IsFile(page) {
		sandbox.Deps.Std.Log("add-page: %s already exists, keeping it \n", page)
		return nil
	}

	sandbox.Deps.Std.Log("add-page creating %s \n", page)

	content, err := sandbox.Deps.Embeddeps.RenderTemplate("templates/page_html.html", vars)
	if err != nil {
		return err
	}
	return io.WriteFile(page, content)
}
