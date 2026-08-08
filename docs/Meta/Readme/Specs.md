# Readme Specification

## Description
Defines the required structure and layout for the project's root `README.md` file.

### Rules
- The `README.md` must strictly follow the section order defined below.
- Every table must follow the formatting defined in [GeneralDoc](/docs/Meta/GeneralDoc/Specs.md).

#### Theme-Based Doc Index
- Documents are indexed by **theme** — what the reader wants to accomplish — not by Diátaxis category. The category (Tutorial, Reference, Explanation) appears as a `Type` column inside each theme table; the category directories and file locations are unaffected.
- Each theme is a `##` section with a one-sentence audience/scope description followed by a **Doc Index table** with `Doc`, `Description`, and `Type` columns:

```markdown
## <Theme Name>

One sentence saying who this theme is for / what it covers.

| Doc | Description | Type |
| --- | --- | --- |
| [/docs/...](/docs/...) | What the reader accomplishes with it | Tutorial / Reference / Explanation |
```

- **Theme ordering follows user need**: the CLI comes first, because that is how most readers meet the project; then the library behind it; then themes for *extending* it (sandbox items, dependencies, adapters); and finally themes for *maintaining* the repo (docs management, template adaptation, project rules). A new reader must hit what they need without scrolling past maintainer-only content.
- **Rows inside a Doc Index table are ordered by reader need**: the docs a reader must hit first (installing, first usage, core operations) come before references, then customization/advanced topics, then samples, with niche material last. A row must never sit above one the reader needs earlier.
- **A doc may appear under multiple themes** — duplicate rows across themes are expected and encouraged.
- **Every top-level doc under `docs/` must appear in at least one theme.** No orphans. Files indexed by an index doc (specifications under `Meta/`, API detail pages under `PublicApi/`) are covered through that index doc's row.
- **Link text and link target must match** and point to the real file location (a Tutorial always links into `/docs/`, never another directory).
- Descriptions are one line, 50–100 characters, and specific — say what the reader gets, never a generic phrase reused across rows.
- Creating, renaming, or deleting a `.md` file requires updating its Doc Index row(s) in the same commit — in **every** theme table that lists it.

## Structure

1. **Title**: The project's name (H1).
2. **Headers/Badges**: Links to relevant external resources.
3. **Short Description**: A brief, single-sentence summary of the project.
4. **Overview**: A high-level explanation of the project's design and purpose.
5. **Quick Start CLI**: A single copy-pasteable command installing the CLI globally, followed by the first commands to run with it. This comes before the library quick start: the CLI is the project's main face, the library behind it is the background feature.
6. **Quick Start Library**: Installing the module and the smallest Go program using it.
7. **Must Read callout**: The required-reading table (Rules, Structure, Specs).
8. **Theme sections**, in this order, each with its description and Doc Index table:
   - *CLI Usage* — installing, running, and extending the command-line interface.
   - *CLI Examples* — carries the table of scripts in `examples/cliExamples/`.
   - *Library Usage* — consuming the same behavior from Go code.
   - *Library Examples* — carries the table of programs in `examples/libraryExamples/`.
   - *Sandbox Management* — adding lib functions and objects inside the sandbox.
   - *Dependency Management* — the `Deps` contract and the adapters filling it.
   - *Documentation Management*, *Template Adaptation*, *Project Rules & Structure* (maintenance).
   - Themes may be added, merged, or renamed when content fits better another way, as long as the ordering rule holds.
9. **License Ref**: A reference link to the project's license file.

> **Note**: For a concrete example, refer to [sample.md](/docs/Meta/Readme/sample.md).
