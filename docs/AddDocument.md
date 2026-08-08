# Add a Document

## Description
Covers creating a new `.md` file in [docs/](../) — or updating an existing one — and registering it across the project.

### Rules
- Every `.md` file must comply with the specifications that govern it — locate them in [Specs.md](/docs/Specs.md).
- A file governed by a specification must reproduce the shape that specification requires.
- Creating a `.md` file requires updating the [README.md](/README.md) Doc Index and [Structure.md](/docs/Structure.md) in the same commit.

---

## Workflow
1. Identify the category the document belongs to:
   - **Reference** — listable material: structures, rules, specifications, API indexes → `docs/References/`.
   - **Explanation** — how a mechanic, feature, or ability works → `docs/Explanations/`.
   - **Tutorial** — a workflow guide for a single goal → `docs/Tutorials/`.
2. Check [Specs.md](/docs/Specs.md) for the specifications matching the file, and read them before writing.
3. Create the `.md` file in the chosen directory with a descriptive, topic-based name (e.g., `AddSample.md`, `PublicApi.md`).
4. Write the content following those specifications, paying special attention to:
   - **Topic-driven structure** — one concern per section.
   - **Conciseness** — short, direct sentences.
   - **Heading hierarchy** — never skip heading levels.
5. Add cross-references using **relative paths**, and add the reverse link in every document that should point back to this one.
6. Add an entry to the [README.md](/README.md) Doc Index, in every **theme** table where the document fits — a doc may appear under multiple themes.
7. Register the file in [Structure.md](/docs/Structure.md).
