# Template File Actions

## Description
Lists every file and directory of this template and the action it takes when the template is forked into a new library or an existing library is adapted to it. Each file falls into exactly one action: **Copy**, **Create**, **Rewrite**, or **Delete**. The workflows using this list are [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) and [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md).

---

## Copy

Taken as-is from the template. They describe the structure itself, not the library, so they carry over unchanged. Adapting them is allowed but never required.

Copying these files carries over the template's **generic** guides and specifications only. The new library must still **create** its own case-specific tutorials and reference pages — see [Create](#create).

| Path | Description |
|------|-------------|
| `docs/References/Meta/*` | The specifications every file of the new library must be shaped by |
| `docs/References/RULES.md` | The binding contribution rules |
| `docs/References/Specs.md` | The index locating each specification |
| `docs/References/TemplateFileActions.md` | This page |
| `docs/Tutorials/*` | The workflow guides |
| `docs/Explanations/*` | The explanations of the library's mechanics and features |
| `sandbox/new.go` | The `New` constructor storing `Deps` on `api.Lib` and running the internal factories over it |

---

## Create

Written from scratch for the library being built or adapted. Nothing of the template's content survives here — the example files occupying these paths are removed by **[Delete](#delete)**. Every created file must be shaped by the specification in its row.

| Path | Description | Specification |
|------|-------------|---------------|
| `adapters/<name>/<name>.go` | One adapter per opinionated implementation of the `Deps` contract | [Adapters](/docs/References/Meta/Adapters/Specs.md) |
| `sandbox/internal/lib/*` | The lib's field factories and the `New` constructor running them all, reaching every dependency through `l.Deps` | [LibFunctions](/docs/References/Meta/LibFunctions/Specs.md) |
| `sandbox/internal/<object>/*` | One package per object the library hands back: its field factories and the `New` constructor running them all | [LibObjects](/docs/References/Meta/LibObjects/Specs.md) |
| `docs/References/PublicApi/*` | One detail page per public API entry | [ReferenceDocs](/docs/References/Meta/ReferenceDocs/Specs.md) |
| `docs/References/<Name>.md` | Any reference page the new library needs beyond the public API index | [ReferenceDocs](/docs/References/Meta/ReferenceDocs/Specs.md) |
| `docs/Tutorials/<Goal>.md` | One tutorial per workflow specific to the new library — the template tutorials carried over by **[Copy](#copy)** do **not** fulfill this | [TutorialDocs](/docs/References/Meta/TutorialDocs/Specs.md) |
| `examples/<example>/<example>.go` | One runnable sample per demonstrated use case | [Examples](/docs/References/Meta/Examples/Specs.md) |

---

## Rewrite

Kept in place, with their content replaced by the new library's. The file keeps its path and its shape; only what it declares or documents changes. Every rewritten file must be shaped by the specification in its row.

| Path | Rewrite with | Specification |
|------|--------------|---------------|
| `README.md` | The new library's overview, quick start, badges, Doc Index, and Samples section | [Readme](/docs/References/Meta/Readme/Specs.md) |
| `sandbox/contracts/deps/deps.go` | The `Deps` function fields the new library requires | [Deps](/docs/References/Meta/Deps/Specs.md) |
| `sandbox/contracts/api/api.go` | The `Lib` struct and one struct per object the new library hands back | [Outputs](/docs/References/Meta/Outputs/Specs.md) |
| `adapters/standard/standard.go` | The default adapter, filling the new `Deps` contract | [Adapters](/docs/References/Meta/Adapters/Specs.md) |
| `docs/References/PublicApi.md` | The index of the new public API entries | [ReferenceDocs](/docs/References/Meta/ReferenceDocs/Specs.md) |
| `docs/References/Structure.md` | The layout of the new library | [Structure](/docs/References/Meta/Structure/Specs.md) |

---

## Delete

The template's example content. Removed once the new library's own files exist.

| Path | Description |
|------|-------------|
| `adapters/*` — except `adapters/standard/` | The example alternative adapters |
| `sandbox/internal/*` | The example lib factories and object packages |
| `docs/References/PublicApi/*` | The example API detail pages |
| `examples/*` | The example samples |
