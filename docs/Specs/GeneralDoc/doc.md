# GeneralDoc Specification

## Description
Defines the baseline shape every `.md` file in this project must have, whatever its kind. It is the common ground the other documentation specifications build on: a README, a tutorial, or a reference page must satisfy **this** spec plus the one for its own kind, located through [Specs](/docs/Specs/doc.md).

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
| "This file is responsible for the implementation of the dispatch logic." | "Implements the dispatch logic." |
| "In order to be able to run the project…" | "To run the project…" |
| "It should be noted that this function returns an error." | "Returns an error on failure." |

#### Cross-Reference Between Documents
- Always use **repository-rooted paths** (e.g. `/docs/PublicApi/doc.md`), never absolute filesystem paths.
- Link to the **most specific section** possible using anchors (e.g. `Commands.md#exit-codes`).
- When referencing content explained elsewhere, **link** to it instead of duplicating it.
- Never link into `old/`: it is reference material, not documentation.

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
- Always specify the language in fenced code blocks (e.g. ` ```go `, ` ```bash `, ` ```yaml `).
- Prefer **runnable** snippets over fragments.
- Add inline comments highlighting the important parts.
- Do not include unrelated boilerplate that distracts from the point being demonstrated.
- Commands a reader types are shown against an **installed** `agnos`; the repository's own bootstrap binary is shown as `./release/bootstrap.bin` when the distinction matters.

#### Consistent Terminology
- Use the same term for the same concept throughout all documentation.
- Define project-specific terms on first use when they are not obvious.

| Concept | Preferred Term |
|---------|---------------|
| The closed tree under `sandbox/` that reaches nothing outside itself | **sandbox** |
| The struct of sub-contract structs the sandbox receives, `deps.Deps` | **Deps** |
| One `sandbox/deps/<x>/` package restating an outside library's api | **sub-contract** |
| One `adapters/libs/<x>/` package filling one `Deps` field | **adapter lib** |
| One `adapters/availables/<name>/` package assembling libs into a `Deps` | **adapter** (or **assembly**) |
| The uniform `Bind(deps *deps.Deps)` an adapter lib exports, and the `<X>Bind` a `sandbox/binds/` file exports | **binder** |
| A struct of function fields the sandbox hands back (`api.Sandbox`, `api.Actions`, `api.Cli`) | **contract** |
| One reusable operation under `sandbox/internal/actions/<name>/` | **action** |
| One package under `sandbox/internal/commands/<name>/` — `entries.yaml`, `entries.go`, `handler.go` | **command** |
| The hand-written `entries.yaml` | **declaration** |
| The generated `dispatch<Name>` reading argv into `Entries` | **dispatch** |
| The `CommandHandler` function in `handler.go` | **handler** |
| A directory under `assets/<group>/` rendered as one unit | **asset group** |
| A directory under `assets/deplist/<dep>/`, installable by `dep-install` | **dep** |
| A `build/collect_*.go` function turning a directory listing into template data | **collector** |
| The `go` or `none` toolchain step after persist | **runtime** |
| A project written by `agnos start` | **generated project** |
| The transactional filesystem rooted at `--path` | **SmartIO** |
| A single-goal, step-by-step guide in `/docs/` | **tutorial** |
| The description of how a file must be shaped | **specification** |

#### Keep Documentation in Sync
- Documentation must reflect the current state of the code.
- When a change affects documentation, update every impacted file in the **same commit**.
- `CLAUDE.md` states the same patterns as rules for the coding agent; a pattern changed here changes there in the same commit.

## Structure
1. **Title** (H1): the subject of the document, one per file.
2. **`## Description`**: one short paragraph stating what the document covers.
3. **Body sections** (`##`), separated by `---`, each covering a single topic and nesting subtopics with `###`/`####`.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/GeneralDoc/sample.md).
