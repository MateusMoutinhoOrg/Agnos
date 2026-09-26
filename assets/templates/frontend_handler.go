package frontend

import (
	"{{.Module}}/sandbox/api"
	"{{.Module}}/sandbox/deps/serverdeps"
	"{{.Module}}/sandbox/internal/frontio"
)

// spaFallback answers a path that names no file, and has no extension, with
// the index.html of assets/frontend instead of declining. Turn it on for a
// single-page app whose router owns the url (React Router, Vue Router, ...):
// /users/42 is then the app's to draw rather than a 404. A missing /app.js
// still declines, so a broken asset link stays a 404 and never becomes html.
const spaFallback = false

// InternalPureHandler answers GET /<rest...> with the file of assets/frontend
// the path names, read out of the binary through frontio.Resolve: the path
// itself, then <path>.html, then <path>/index.html. Whatever lands in that
// tree — hand-written html or the dist/ of any bundler — is served as it is.
//
// This route is declared with the highest priority number of the project, so
// every api route runs first and this one is the fallback in front of the 404.
// A path that names no file is declined — nil without answering — so the chain
// goes on and HandleNotFound answers it. frontio.SafePath is what keeps the
// caller's path inside assets/frontend; it is generated and rewritten by every
// build, so the check is never yours to keep.
func InternalPureHandler(sandbox *api.Sandbox, route *api.Route, entries *Entries, response *serverdeps.Response) error {
	relative, content, ok := frontio.Resolve(sandbox, entries.Rest)
	if !ok && spaFallback && frontio.ExtensionOf(sandbox, entries.Rest) == "" {
		relative, content, ok = frontio.Resolve(sandbox, "")
	}
	if !ok {
		return nil
	}

	response.SetHeader("Content-Type", frontio.ContentTypeOf(sandbox, relative))
	response.SetHeader("Cache-Control", frontio.RevalidateCache)
	response.SetStatus(api.StatusOk)
	response.Write(content)

	return nil
}
