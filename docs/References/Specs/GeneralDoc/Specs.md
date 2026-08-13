# GeneralDoc Specification

## Description
Defines the baseline shape every `.md` file in this project must have, whatever its kind. It is the common ground the other documentation specifications build on: a README, a tutorial, or a reference page must satisfy **this** spec plus the one for its own kind, located through [Specs.md](/docs/References/Specs.md).

### Rules

#### Topic-Driven Structure
- Organize the document around **topics** and **subtopics** expressed as Markdown headings.
- Each section covers a single concern.
- Separate top-level sections with horizontal rules (`---`).

#### Heading Hierarchy

| Level | Usage |
|-------|-------|
| `#` (H1) | Document title — **one per file** |
| `##` (H2) | Major sections (e.g. `Description`, `Workflow`, a top-level directory) |
| `###` (H3) | Subsections (e.g. a subdirectory, a grouped set of steps) |
| `####` (H4) | Minor details within a subsection |

- Never skip heading levels (e.g. never jump from `##` to `####`).

#### Be Concise
- Write short, direct sentences.
- Avoid filler words, redundant explanations, and unnecessary qualifiers.

| ❌ Avoid | ✅ Prefer |
|----------|----------|
| "This file is responsible for the implementation of the authentication logic." | "Implements the authentication logic." |
| "In order to be able to run the project…" | "To run the project…" |
| "It should be noted that this function returns an error." | "Returns an error on failure." |

#### Cross-Reference Between Documents
- Always use **relative paths** (e.g. `../../docs/References/PublicApi.md`), never absolute filesystem paths.
- Link to the **most specific section** possible using anchors (e.g. `Structure.md#adapters`).
- When referencing content explained elsewhere, **link** to it instead of duplicating it.

#### Avoid Duplication
- Information needed in several places is written **once** and linked from everywhere else.
- Duplicated content drifts out of sync over time.

#### Use File Tables for Directory Descriptions
- Describe the contents of a directory with a Markdown table using `File` and `Description` columns.

```markdown
### `/config/`

| File | Description |
|------|-------------|
| `config.go` | Loads and validates environment configuration |
| `defaults.go` | Defines default values for all settings |
```

#### Code Examples
- Always specify the language in fenced code blocks (e.g. ` ```go `, ` ```bash `).
- Prefer **runnable** snippets over fragments.
- Add inline comments highlighting the important parts.
- Do not include unrelated boilerplate that distracts from the point being demonstrated.

#### Consistent Terminology
- Use the same term for the same concept throughout all documentation.
- Define project-specific terms on first use when they are not obvious.

| Concept | Preferred Term |
|---------|---------------|
| The object that fills the dependency contract | **adapter** |
| The struct declaring the injectable behaviors as function fields | **Deps** |
| A contract struct the library hands back to callers | **output** |
| An output struct the library creates and propagates `Deps` into | **library object** |
| A function filling one field of a struct contract, taking a pointer to the struct that carries the state | **factory** |
| The struct a factory takes a pointer to — an output struct in `/sandbox/`, the adapter struct in `/adapters/` | **carrier** |
| The `New` function that builds one carrier, runs every factory of its package over it, and returns the filled contract struct | **constructor** |
| A runnable example in `/examples/libraryExamples/` or `/examples/cliExamples/` | **sample** |
| The `api.Lib.Sandboxmain` field and the dispatch behind it | **command-line interface** |
| A word the interface dispatches on, such as `balance` | **command** |
| A single-goal, step-by-step guide in `/docs/` | **tutorial** |
| The description of how a file must be shaped | **specification** |

#### Keep Documentation in Sync
- Documentation must reflect the current state of the code.
- When a change affects documentation, update every impacted file in the **same commit**.

## Structure
1. **Title** (H1): the subject of the document, one per file.
2. **`## Description`**: one short paragraph stating what the document covers.
3. **Body sections** (`##`), separated by `---`, each covering a single topic and nesting subtopics with `###`/`####`.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/GeneralDoc/sample.md).
