package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// A page is a route with an html asset beside it, so it needs no declaration
// of its own: sandbox/internal/routes/<page>/route.yaml is the declaration and
// assets/frontend/pages/<page>.html is the content. The helpers below are the
// one place that pairing is spelled, so add-page, remove-page, remove-route
// and front-purge all read it the same way.

// FrontendDir is the root of the asset tree the front layer reads. Everything
// under it is the project's own content: no build writes there, and front-purge
// leaves it whole.
const FrontendDir = "assets/frontend"

// PagesDir is the directory of that tree holding one html template per
// declared page.
const PagesDir = FrontendDir + "/pages"

// StaticAssetsDir is the directory the static route serves, and the one the
// pageio helpers build their links against.
const StaticAssetsDir = FrontendDir + "/static"

// StaticRouteName is the route front-init writes to serve StaticAssetsDir. It
// is a page's one reserved name: a page of that name would take the mount over
// from the route serving it.
const StaticRouteName = "static"

// PageAsset is the project-relative path of a page's html template. The name
// is normalized the way a route name is, so both the identifier ("about-us")
// and the package directory it becomes ("about_us") name the same file.
func PageAsset(sandbox *api.Sandbox, name string) string {
	return PagesDir + "/" + RouteIdentifier(sandbox, name) + ".html"
}

// IsPage reports whether a route name has an html template beside it — the
// single test that tells a page from any other route.
func IsPage(sandbox *api.Sandbox, io *smartio.SmartIO, name string) bool {
	return io.IsFile(PageAsset(sandbox, name))
}
