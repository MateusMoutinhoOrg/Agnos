# Structure Specification

## Description
Defines the required shape of `docs/References/Structure.md`. Structure.md describes the project's **schema** — the *kinds* of files a project is built from — not an exhaustive listing of every concrete file.

### Rules
- Structure.md documents **schema slots**, not individual files. One row represents a kind of file (e.g. "a command's handler", "an adapter lib"), not a specific instance.
- Each top-level directory of the project gets its own `##` section, separated by `---`.
- Directory contents are described with Markdown tables (see [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md#use-file-tables-for-directory-descriptions)).
- Structure.md must **not** link to individual specifications. A slot governed by a specification carries a **Spec** column naming it; the reader resolves that name through [Specs.md](/docs/References/Specs.md).
- A slot written by `agnos build` is marked *generated* in its description and carries no Spec: its shape is its template's.
- The document must let a reader construct the project from scratch — including creating `AgnosConfig/`, `sandbox/`, `adapters/`, `assets/`, `cmd/`, and `docs/` — using the schema it maps and the specifications it names, and must say which slots a generated project shares with this repository.
- The intro must state the one-way dependency flow between the code trees: `adapters/` and `cmd/` depend on `sandbox/`, and `sandbox/` depends on nothing outside itself.
- Slots with no meaningful shape to specify (e.g. `LICENSE`, `go.mod`) leave the Spec cell empty.

## Structure
1. **Title** (H1): `# Project Structure`.
2. **Intro**: one short paragraph explaining that the document is a schema map and that named specs are resolved through `docs/References/Specs.md`, followed by the dependency-flow diagram.
3. **One `##` section per top-level directory** (`Root`, `/AgnosConfig/`, `/sandbox/`, `/adapters/`, `/assets/`, `/cmd/`, `/release/`, `/old/`, `/docs/`), each with a `File | Description | Spec` table (or nested subsection tables), separated by `---`. Directories nested inside a top-level directory are documented as subsections within their parent's `##` section.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/Structure/sample.md).
