# Readme Specification

## Description
Defines the required structure and layout for the project's root `README.md` file.

### Rules
- The `README.md` must strictly follow the section order defined below.
- Every table must follow the formatting defined in [GeneralDoc](/docs/References/Meta/GeneralDoc/Specs.md).

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

- **Theme ordering follows user need**: themes for *using* the library come first (initialization, API usage, samples), then themes for *extending* it (lib functions, dependencies, adapters), and finally themes for *maintaining* the repo (docs management, template adaptation, project rules). A new reader must hit what they need without scrolling past maintainer-only content.
- **A doc may appear under multiple themes** — duplicate rows across themes are expected and encouraged.
- **Every top-level doc under `docs/` must appear in at least one theme.** No orphans. Files indexed by an index doc (specifications under `Meta/`, API detail pages under `PublicApi/`) are covered through that index doc's row.
- **Link text and link target must match** and point to the real file location (a Tutorial always links into `/docs/Tutorials/`, never another directory).
- Descriptions are one line, 50–100 characters, and specific — say what the reader gets, never a generic phrase reused across rows.
- Creating, renaming, or deleting a `.md` file requires updating its Doc Index row(s) in the same commit — in **every** theme table that lists it.

## Structure

1. **Title**: The project's name (H1).
2. **Headers/Badges**: Links to relevant external resources.
3. **Short Description**: A brief, single-sentence summary of the project.
4. **Overview**: A high-level explanation of the project's design and purpose.
5. **Quick Start**: Step-by-step instructions to install and run a basic example.
6. **Must Read callout**: The required-reading table (Rules, Structure, Specs).
7. **Theme sections**, ordered usage → extension → maintenance, each with its description and Doc Index table:
   - *Library Usage*, *Samples* (usage) — the Samples theme also carries the table of runnable examples.
   - *Extending the Library*, *Dependency Management* (extension).
   - *Documentation Management*, *Template Adaptation*, *Project Rules & Structure* (maintenance).
   - Themes may be added, merged, or renamed when content fits better another way, as long as the ordering rule holds.
8. **License Ref**: A reference link to the project's license file.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Meta/Readme/sample.md).
