## Doc Refactor

### New Doc Tree
- docs/Index/<theme-id>.md
- docs/<doc-name>/doc.md
- docs/<doc-name>/props.yaml
- docs/<doc-name>/Index.md
- docs/<doc-name>/<sub-doc-name>/ (same shape, recursive)

`Index` is a reserved name at the first level. Any other file inside a doc dir is an asset and is ignored by the collector.

### docs/Index/<theme-id>.md
#### automated: yes
#### description:
One file per theme in themes.yaml, listing every first-level doc whose props.yaml names that theme id, sorted by `order` then name. The generator owns `docs/Index/` entirely: stale indexes are deleted.
#### model:
```
# <Theme Name>
<Theme Description>

| Doc | Description |
| --- | --- |
| [<Doc Name>](/docs/<doc-name>/doc.md) | <doc description> |
```

### docs/<doc-name>/doc.md
#### automated: no
#### description:
The documentation the user wrote. Fixed filename so the path is derivable from the dir name alone.

### docs/<doc-name>/props.yaml
#### automated: no
#### description:
Doc metadata, same shape as themes.yaml. `themes` is required on first-level docs and forbidden on sub-docs. `order` is optional (default: alphabetical).
#### sample:
```yaml
name: Public Api
description: Every exported symbol of the sandbox and adapters
themes:
  - lib-usage
  - development
order: 3
```

### docs/<doc-name>/Index.md
#### automated: yes
#### description:
Generated only when the doc dir has sub-doc dirs. Same table model as the theme index, listing the direct sub-docs. `doc.md` links to it by hand once.

### docs/<doc-name>/<sub-doc-name>/
#### automated: no
#### description:
A sub-doc is a doc: `doc.md` + `props.yaml` + optional sub-docs, to any depth. It appears only in its parent's `Index.md`, never in a theme index.

### Verify rules
- props.yaml missing or unparsable → error
- theme id not in themes.yaml → error
- first-level doc with no themes → error
- theme with no docs → error
- `themes` set on a sub-doc → error

### Task:
- Refactor all the documentation to match this pattern
- Adapt the code to generate `docs/Index/*.md` and every `docs/**/Index.md` on `agnos build`, and add the verify rules above
- Update `assets/all/README.md` to link `/docs/Index/{{ .Id }}.md`
