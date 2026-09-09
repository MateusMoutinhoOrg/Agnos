# FrontUsage

The front layer answers html. A **page** is a route with a template beside it:
`sandbox/internal/routes/<page>/` declares and handles it,
`assets/frontend/pages/<page>.html` is what it renders. `sandbox/internal/pageio` is the
render layer both sides go through, and it is the directory `build` reads the layer from.

## Bring it up

```bash
{{.GeneratorName}} front-init                       # installs the deps, renders pageio, writes the static route
{{.GeneratorName}} add-page home --trigger /        # a page answering GET /
{{.GeneratorName}} add-page about --title "About"   # a page answering GET /about
{{.GeneratorName}} remove-page about                # drops the route and the html
{{.GeneratorName}} front-purge                      # drops the layer, keeps assets/frontend/
```

Then `{{.Name}} start-server` serves them.

`front-init` runs `server-init` first when the project has no server layer. A page is a
route, so `docs/Routes` lists it and every route editor (`set-route`, `add-param`,
`add-header`, …) works on its `route.yaml`.

## Generated vs yours

| Path | Written by | Rewrite |
|---|---|---|
| `sandbox/internal/pageio/templates.go` | `build` | always |
| `docs/FrontUsage/` | `build` | always |
| `sandbox/internal/routes/static/{route.yaml,handler.go}` | `front-init` | once |
| `assets/frontend/static/styles/main.css`, `.../scripts/main.js` | `front-init` | once |
| `sandbox/internal/routes/<page>/{route.yaml,handler.go}` | `add-page` | once |
| `assets/frontend/pages/<page>.html` | `add-page` | once |
| `sandbox/internal/routes/<page>/entries.go` | `build` | always |

`once` files are a starting point and yours from the moment they exist. To go back to the
scaffolded one: `{{.GeneratorName}} remove-route static && {{.GeneratorName}} front-init` for
the static route, `{{.GeneratorName}} remove-page <page> && {{.GeneratorName}} add-page <page>`
for a page — that one deletes the html too. `{{.GeneratorName}} front-purge` followed by
`front-init` returns the whole layer without touching `assets/frontend/`.

Whoever edits `sandbox/internal/routes/static/handler.go` keeps `safeSegments`: it is the only
thing between a caller's `/static/<path>` and the rest of the embedded asset tree.

## Helpers

Every page is rendered through `pageio.Render`, which registers these. A path is
slash-separated and relative to `assets/frontend/static`; one that cannot be read fails the
render, so a broken link is a 500 on the page that carries it and never a 404 in the browser.

| Helper | Returns |
|---|---|
| `{{"{{ staticref \"styles/main.css\" }}"}}` | the url, stamped `?sha=` with the digest of the current bytes |
| `{{"{{ cssref \"styles/main.css\" }}"}}` | the whole `<link rel="stylesheet">` tag |
| `{{"{{ jsref \"scripts/main.js\" }}"}}` | the whole deferred `<script>` tag |
| `{{"{{ dirref \"styles\" }}"}}` | every `.css` then every `.js`/`.mjs` at or below the directory, sorted; `""` and `"."` name the root |
| `{{"{{ inline \"styles/critical.css\" }}"}}` | the file's bytes, written straight into the page |
| `{{"{{ include \"frontend/pages/nav.html\" }}"}}` | another template, rendered over the same vars (path relative to the assets root) |

`?sha=` is a cache key and never an input: the static route serves the file the path names
whatever the sha says, and answers `immutable` when the two agree, `no-cache` when they do not.
So a stale link is slow, never broken.

## Page vars

A page's handler declares `pageVars`, and every `{{"{{ .Field }}"}}` of its html is one
exported field of it — a new variable in the page is a new field, checked by the compiler
rather than discovered at request time.

```go
content, err := pageio.Render(deps, pageAsset, pageVars{Title: "Home"})
```

## The mount

`pageio.StaticMount` is generated from the `identifier` of the first segment of
`sandbox/internal/routes/static/route.yaml`, which is `/static` as scaffolded. Move it with
`{{.GeneratorName}} remove-segment` / `add-segment --identifier /assets` and the next build
moves every link the helpers write. Editing the constant instead does nothing: it is rewritten
from the declaration.
