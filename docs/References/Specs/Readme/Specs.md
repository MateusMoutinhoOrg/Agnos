# Readme Specification

## Description
Defines the required structure and layout for the project's root `README.md` file.

### Rules
- The `README.md` must strictly follow the section order defined below.
- Every table must follow the formatting defined in [GeneralDoc](/docs/References/Specs/GeneralDoc/Specs.md).

#### Theme-Based Doc Index
- The README is a **pointer, not a catalogue**: it links to each theme index under `docs/Index/`, and to nothing else inside `docs/`. Each index then lists the pages of its own theme — see the [Index](/docs/References/Specs/Index/Specs.md) specification.
- Themes are the pages of `docs/Index/`: `CliUsage.md`, `LibUsage.md`, `Development.md`, `Templating.md`.
- The Doc Index is a single table with `Theme` and `Description` columns:

```markdown
| Theme | Description |
| --- | --- |
| [<Theme Name>](/docs/Index/<Theme>.md) | Who the theme is for and what it covers |
```

- **Theme ordering follows user need**: the CLI comes first, because that is how most readers meet the project; then the library behind it; then contributing to it; and finally adapting it as a template. A new reader must hit what they need without scrolling past maintainer-only content.
- **Link text and link target must match**, and every target is a page of `docs/Index/`.
- Descriptions are one line, 50–100 characters, and specific — say who the theme is for, never a generic phrase reused across rows.
- Adding, renaming, or deleting a **theme** requires updating this table in the same commit. Adding, renaming, or deleting a page inside a theme does **not** touch the README — it touches that theme's index.

## Structure

1. **Title**: The project's name (H1).
2. **Headers/Badges**: Links to relevant external resources.
3. **Short Description**: A brief, single-sentence summary of the project.
4. **Overview**: A high-level explanation of the project's design and purpose.
5. **Doc Index**: The theme table described above, optionally followed by one line pointing newcomers at the quick starts.
6. **License Ref**: A reference link to the project's license file.

> **Note**: For a concrete example, refer to [sample.md](/docs/References/Specs/Readme/sample.md).
