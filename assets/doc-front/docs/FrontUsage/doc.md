# FrontUsage

The front layer serves every file of `assets/frontend/`, embedded in the binary, over http.
A **page** is a file of that tree and nothing else: no route, no declaration, no template
syntax. Hand-written html and the `dist/` of any bundler (Vite, React, Svelte, Astro...) are
served the same way. Dynamic data comes from api routes the page's js calls, or that a
`<form>` posts to.

## Bring it up

```bash
{{.GeneratorName}} front-init                    # frontio, the frontend route, assets/frontend/{index,404}.html
{{.GeneratorName}} add-page about --title About  # assets/frontend/about.html, answered on /about
{{.GeneratorName}} add-page blog/post            # assets/frontend/blog/post.html, answered on /blog/post
{{.GeneratorName}} remove-page about             # deletes the html
{{.GeneratorName}} front-purge                   # drops the layer, keeps assets/frontend/
```

Then `{{.Name}} start-server` serves them. `front-init` runs `server-init` first when the
project has no server layer. Any file put in `assets/frontend/` by hand is served the same
way as one `add-page` wrote: `add-page` is a scaffold, not a declaration.

## How a path is answered

The `frontend` route matches `GET /<rest...>` at priority `1000`, so every api route runs
before it. It reads, in order, the first one of these that exists:

| Request | Tried |
|---|---|
| `/` | `index.html` |
| `/<p>` | `<p>`, `<p>.html`, `<p>/index.html` |

A path that names no file is answered `404` with `404.html`, the formatted page `front-init`
writes — restyle it by editing it. Delete it and such a path is declined instead, so the chain
goes on and `handle_not_found.go` answers the `404`. Every file is sent with the `Content-Type` of its extension (`frontio.ContentTypeOf`,
unknown ones as `application/octet-stream`) and `Cache-Control: no-cache`.

## A form

A plain `<form method="POST" action="/login">` reaches a route whose body is `type: form`, with
no script: the browser sends `application/x-www-form-urlencoded` and navigates to what the route
answers. Each input's `name` is a property of the route's `form-schema`
([RouteYaml](../RouteYaml/doc.md#form-schema)), bound onto `Entries.Body` already typed.

```bash
{{.GeneratorName}} add-route login --method POST --trigger /login
{{.GeneratorName}} set-body login --type form --required
{{.GeneratorName}} add-body-field username --route login --required
{{.GeneratorName}} add-body-field password --route login --required
```

A request that fails the schema is answered by `handle_bad_request.go` before the handler runs.
Leave `enctype` alone: `multipart/form-data` is not read.

## A bundler's build

Point the bundler's output at `assets/frontend/` (Vite: `build.outDir`, `emptyOutDir: true`),
build it, then build {{.Name}}: the binary embeds whatever is there. Links stay relative to
`/`, because the tree is served from the root.

For a single-page app whose router owns the url, set `spaFallback = true` in
`sandbox/internal/routeslist/frontend/InternalPureHandler.go`: a path with no extension that
names no file is then answered with `index.html`. A missing `/app.js` still gets the `404`.

## Generated vs yours

| Path | Written by | Rewrite |
|---|---|---|
| `sandbox/internal/generated/frontio/frontio.go` | `build` | always |
| `docs/FrontUsage/` | `build` | always |
| `sandbox/internal/routeslist/frontend/{route.yaml,InternalPureHandler.go}` | `front-init` | once |
| `sandbox/internal/routeslist/frontend/{new.go,entries.go}` | `build` | always |
| `assets/frontend/index.html` | `front-init` | once, kept if already there |
| `assets/frontend/404.html` | `front-init` | once, kept if already there |
| `assets/frontend/<page>.html` | `add-page` | once, refused if already there |

The frontend route is a route like any other: every editor of its `route.yaml` works on it, and
`{{.GeneratorName}} remove-route frontend && {{.GeneratorName}} front-init` scaffolds it again.
`frontio.SafePath` is the only thing between a caller's path and the rest of the embedded asset
tree; it is generated, so the handler never has to keep a copy of it.
