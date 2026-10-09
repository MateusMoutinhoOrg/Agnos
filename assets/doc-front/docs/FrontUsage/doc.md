# FrontUsage

The front layer serves every file of `assets/front/`, embedded in the binary, over http.
A **page** is a file of that tree and nothing else: no route, no declaration, no template
syntax. Hand-written html and the `dist/` of any bundler (Vite, React, Svelte, Astro...) are
served the same way. Dynamic data comes from api routes the page's js calls, or that a
`<form>` posts to.

## Bring it up

```bash
{{.GeneratorName}} front-init                    # the OpinionatedAgnosFront lib, the front route, assets/front/{index,404}.html
{{.GeneratorName}} add-page about --title About  # assets/front/about.html, answered on /about
{{.GeneratorName}} add-page blog/post            # assets/front/blog/post.html, answered on /blog/post
{{.GeneratorName}} remove-page about             # deletes the html
{{.GeneratorName}} front-purge                   # drops the layer, keeps assets/front/
```

Then `{{.ProjectName}} start-server` serves them. `front-init` runs `server-init` first when the
project has no server layer. Any file put in `assets/front/` by hand is served the same
way as one `add-page` wrote: `add-page` is a scaffold, not a declaration.

## How a path is answered

The `front` route matches `GET /<rest...>` at priority `1000`, so every api route runs
before it. It reads, in order, the first one of these that exists:

| Request | Tried |
|---|---|
| `/` | `index.html` |
| `/<p>` | `<p>`, `<p>.html`, `<p>/index.html` |

A path that names no file is answered `404` with `404.html`, the formatted page `front-init`
writes — restyle it by editing it. Delete it and such a path is declined instead, so the chain
goes on and `handle_not_found.go` answers the `404`. Every file is sent with the `Content-Type` of its extension (`OpinionatedAgnosFront.ContentTypeOf`,
unknown ones as `application/octet-stream`) and `Cache-Control: no-cache`.

## A form

A plain `<form method="POST" action="/login">` reaches a route whose body is `type: form`, with
no script: the browser sends `application/x-www-form-urlencoded` and navigates to what the route
answers. Each input's `name` is a property of the route's `form-schema`
([RouteYaml](../RouteYaml/doc.md#form-schema)), bound onto `Input.Body` already typed.

```bash
{{.GeneratorName}} add-route login --method POST --trigger /login
{{.GeneratorName}} set-body login --type form --required
{{.GeneratorName}} add-body-field username --route login --required
{{.GeneratorName}} add-body-field password --route login --required
```

A request that fails the schema is answered by `handle_bad_request.go` before the handler runs.
Leave `enctype` alone: `multipart/form-data` is not read.

## A bundler's build

Point the bundler's output at `assets/front/` (Vite: `build.outDir`, `emptyOutDir: true`),
build it, then build {{.ProjectName}}: the binary embeds whatever is there. Links stay relative to
`/`, because the tree is served from the root.

For a single-page app whose router owns the url, set `spaFallback = true` in
`sandbox/internal/routes/front/handler.go`: a path with no extension that
names no file is then answered with `index.html`. A missing `/app.js` still gets the `404`.

## Generated vs yours

| Path | Written by | Rewrite |
|---|---|---|
| `sandbox/deps/OpinionatedAgnosFront/`, `adapters/impls/OpinionatedAgnosFront/` | `front-init` | once, like any dep |
| `docs/FrontUsage/` | `build` | always |
| `sandbox/internal/routes/front/{route.yaml,handler.go}` | `front-init` | once |
| `sandbox/internal/routes/front/generated.{new,input}.go` | `build` | always |
| `assets/front/index.html` | `front-init` | once, kept if already there |
| `assets/front/404.html` | `front-init` | once, kept if already there |
| `assets/front/<page>.html` | `add-page` | once, refused if already there |

The front route is a route like any other: every editor of its `route.yaml` works on it, and
`{{.GeneratorName}} remove-route front && {{.GeneratorName}} front-init` scaffolds it again.
`OpinionatedAgnosFront.SafePath`, which `Resolve` runs first, is the only thing between a caller's
path and the rest of the embedded asset tree; it lives in the lib, so the handler never has to
keep a copy of it.
