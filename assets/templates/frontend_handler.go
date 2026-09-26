package frontend

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/generated/frontio"
)

// spaFallback answers a path that names no file, and has no extension, with
// the index.html of assets/frontend instead of the 404 page. Turn it on for a
// single-page app whose router owns the url (React Router, Vue Router, ...):
// /users/42 is then the app's to draw rather than a 404. A missing /app.js
// still gets the 404, so a broken asset link never becomes the app's html.
const spaFallback = false

// InternalPureHandler answers GET /<rest...> with the file of assets/frontend
// the path names, read out of the binary through frontio.Resolve: the path
// itself, then <path>.html, then <path>/index.html. Whatever lands in that
// tree — hand-written html or the dist/ of any bundler — is served as it is.
//
// This route is declared with the highest priority number of the project, so
// every api route runs first and this one is the fallback in front of the 404.
// A path that names no file is answered with frontio.NotFound, the formatted
// 404.html of assets/frontend, under a 404 status; only when that file is
// gone too is the path declined — nil without answering — so the chain goes on
// and HandleNotFound answers it. frontio.SafePath is what keeps the caller's
// path inside assets/frontend; it is generated and rewritten by every build,
// so the check is never yours to keep.
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	relative, content, ok := frontio.Resolve(sandbox, entries.Rest)
	if !ok && spaFallback && frontio.ExtensionOf(sandbox, entries.Rest) == "" {
		relative, content, ok = frontio.Resolve(sandbox, "")
	}
	if !ok {
		relative, content, ok = frontio.Resolve(sandbox, frontio.NotFound)
	}
	if !ok {
		return nil
	}

	status := api.StatusOk
	if relative == frontio.NotFound {
		status = api.StatusNotFound
	}

	response.SetHeader("Content-Type", frontio.ContentTypeOf(sandbox, relative))
	response.SetHeader("Cache-Control", frontio.RevalidateCache)
	response.SetStatus(status)
	response.Write(content)

	return nil
}
