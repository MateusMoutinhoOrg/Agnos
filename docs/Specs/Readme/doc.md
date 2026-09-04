# Readme Specification

## Description
Defines the required structure and layout for the project's root `README.md` file.

### Rules
- `README.md` is generated: `agnos build` renders `<ProjectName>Config/docs/ReadmeHeader.md` (itself a Go `text/template` with the full build `vars` in scope) into it on every build. Edit `ReadmeHeader.md`, never `README.md`.
- The `README.md` must strictly follow the section order defined below.
- Every table must follow the formatting defined in [GeneralDoc](/docs/Specs/GeneralDoc/doc.md).

#### Theme-Based Doc Index
- The README is a **pointer, not a catalogue**: it links to each theme index under `docs/Index/`, and to nothing else inside `docs/`. Each index then lists the pages of its own theme — see the [Index](/docs/Specs/Index/doc.md) specification.
- Themes are the entries of `<ProjectName>Config/themes.yaml`; the README's table is generated from them, one row per entry, linking `/docs/Index/<theme-id>.md`.
- The Doc Index is a single table with `Theme` and `Description` columns:

```markdown
| Theme | Description |
| --- | --- |
| `[<Theme Name>](/docs/Index/<theme-id>.md)` | Who the theme is for and what it covers |
```

- **Theme ordering follows user need**: the CLI comes first, because that is how most readers meet the project; then the project it generated, which is what they work in next; then the library behind the CLI; and finally contributing to it. A new reader must hit what they need without scrolling past maintainer-only content.
- **Link text and link target must match**, and every target is a page of `docs/Index/`.
- Descriptions are one line, 50–100 characters, and specific — say who the theme is for, never a generic phrase reused across rows.
- Adding, renaming, or deleting a **theme** is an edit to `themes.yaml` plus a build; adding, renaming, or deleting a doc inside a theme does **not** touch the README — it touches that theme's generated index.

## Structure

1. **Title**: The project's name (H1).
2. **Headers/Badges**: Links to relevant external resources.
3. **Short Description**: A brief, single-sentence summary of the project.
4. **Overview**: A high-level explanation of the project's design and purpose.
5. **Doc Index**: The theme table described above, optionally followed by one line pointing newcomers at the quick starts.
6. **License Ref**: A reference link to the project's license file.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/Readme/sample.md).
