# FrontUsage

The front layer serves every file of `assets/frontend/`, embedded in the binary, over http.
A **page** is a file of that tree and nothing else: no route, no declaration, no template
syntax. Hand-written html and the `dist/` of any bundler (Vite, React, Svelte, Astro...) are
served the same way. Dynamic data comes from api routes the page's js calls.

## Bring it up

```bash
{{.GeneratorName}} front-init                    # frontio, the frontend route, assets/frontend/index.html
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

A path that names no file is declined, so the chain goes on and `handle_not_found.go` answers
`404`. Every file is sent with the `Content-Type` of its extension (`frontio.ContentTypeOf`,
unknown ones as `application/octet-stream`) and `Cache-Control: no-cache`.

## A bundler's build

Point the bundler's output at `assets/frontend/` (Vite: `build.outDir`, `emptyOutDir: true`),
build it, then build {{.Name}}: the binary embeds whatever is there. Links stay relative to
`/`, because the tree is served from the root.

For a single-page app whose router owns the url, set `spaFallback = true` in
`sandbox/internal/routeslist/frontend/InternalPureHandler.go`: a path with no extension that
names no file is then answered with `index.html`. A missing `/app.js` still declines.

## Generated vs yours

| Path | Written by | Rewrite |
|---|---|---|
| `sandbox/internal/frontio/frontio.go` | `build` | always |
| `docs/FrontUsage/` | `build` | always |
| `sandbox/internal/routeslist/frontend/{route.yaml,InternalPureHandler.go}` | `front-init` | once |
| `sandbox/internal/routeslist/frontend/{new.go,entries.go}` | `build` | always |
| `assets/frontend/index.html` | `front-init` | once, kept if already there |
| `assets/frontend/<page>.html` | `add-page` | once, refused if already there |

The frontend route is a route like any other: every editor of its `route.yaml` works on it, and
`{{.GeneratorName}} remove-route frontend && {{.GeneratorName}} front-init` scaffolds it again.
`frontio.SafePath` is the only thing between a caller's path and the rest of the embedded asset
tree; it is generated, so the handler never has to keep a copy of it.
