package front

import (
	"{{.Module}}/sandbox/api"
	opinionatedagnosfront "{{.Module}}/sandbox/deps/OpinionatedAgnosFront"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/routeprops"
)

// spaFallback answers a path that names no file, and has no extension, with
// the index.html of assets/front instead of the 404 page. Turn it on for a
// single-page app whose router owns the url (React Router, Vue Router, ...):
// /users/42 is then the app's to draw rather than a 404. A missing /app.js
// still gets the 404, so a broken asset link never becomes the app's html.
const spaFallback = false

// Handle answers GET /<rest...> with the file of assets/front
// the path names, read out of the binary through Deps.OpinionatedAgnosFront.Resolve: the path
// itself, then <path>.html, then <path>/index.html. Whatever lands in that
// tree — hand-written html or the dist/ of any bundler — is served as it is.
//
// This route is declared with the highest priority number of the project, so
// every api route runs first and this one is the fallback in front of the 404.
// A path that names no file is answered with opinionatedagnosfront.NotFound, the
// formatted 404.html of assets/front, under a 404 status; only when that
// file is gone too is the path declined — nil without answering — so the chain
// goes on and HandleNotFound answers it. The lib's SafePath, which Resolve
// runs first, is what keeps the caller's path inside assets/front; it lives
// in the OpinionatedAgnosFront lib, so the check is never yours to keep.
func Handle(sandbox *api.Sandbox, props *routeprops.RouteProps, input *Input, response *serverdeps.Response) error {
	relative, content, ok := sandbox.Deps.OpinionatedAgnosFront.Resolve(sandbox.Deps.EmbedDeps, input.Rest)
	if !ok && spaFallback && sandbox.Deps.OpinionatedAgnosFront.ExtensionOf(input.Rest) == "" {
		relative, content, ok = sandbox.Deps.OpinionatedAgnosFront.Resolve(sandbox.Deps.EmbedDeps, "")
	}
	if !ok {
		relative, content, ok = sandbox.Deps.OpinionatedAgnosFront.Resolve(sandbox.Deps.EmbedDeps, opinionatedagnosfront.NotFound)
	}
	if !ok {
		return nil
	}

	status := api.StatusOK
	if relative == opinionatedagnosfront.NotFound {
		status = api.StatusNotFound
	}

	response.SetHeader("Content-Type", sandbox.Deps.OpinionatedAgnosFront.ContentTypeOf(relative))
	response.SetHeader("Cache-Control", opinionatedagnosfront.RevalidateCache)
	response.SetStatus(status)
	response.Write(content)

	return nil
}
