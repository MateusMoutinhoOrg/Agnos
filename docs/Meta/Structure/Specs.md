# Structure Specification

## Description
Defines the required shape of `docs/Structure.md`. Structure.md describes the project's **schema** — the *kinds* of files a project is built from — not an exhaustive listing of every concrete file.

### Rules
- Structure.md documents **schema slots**, not individual files. One row represents a kind of file (e.g. "an adapter", "a lib function"), not a specific instance.
- Each top-level directory of the project gets its own `##` section, separated by `---`.
- Directory contents are described with Markdown tables (see [GeneralDoc](/docs/Meta/GeneralDoc/Specs.md#use-file-tables-for-directory-descriptions)).
- Structure.md must **not** link to individual specifications. A slot governed by a specification carries a **Spec** column naming it; the reader resolves that name through [Specs.md](/docs/Specs.md).
- The document must let a reader construct the project from scratch — including creating `sandbox/`, `sandbox/contracts/`, `sandbox/internal/`, `adapters/`, `cmd/`, `examples/cliExamples/`, `examples/libraryExamples/`, and `docs/` — using the schema it maps and the specifications it names.
- The intro must state the one-way dependency flow between the code trees: `adapters/`, `cmd/`, and `examples/libraryExamples/` depend on `sandbox/`, and `sandbox/` depends on nothing outside itself.
- Slots with no meaningful shape to specify (e.g. `LICENSE`, `go.mod`) leave the Spec cell empty.

## Structure
1. **Title** (H1): `# Project Structure`.
2. **Intro**: one short paragraph explaining that the document is a schema map and that named specs are resolved through `docs/Specs.md`.
3. **One `##` section per top-level directory** (`Root`, `/sandbox/`, `/adapters/`, `/cmd/`, `/examples/cliExamples/`, `/examples/libraryExamples/`, `/docs/`), each with a `File | Description | Spec` table (or nested subsection tables), separated by `---`. Directories nested inside a top-level directory (e.g. `/sandbox/contracts/`, `/sandbox/internal/`) are documented as subsections within their parent's `##` section.

> **Note**: For a concrete example, refer to [sample.md](/docs/Meta/Structure/sample.md).
