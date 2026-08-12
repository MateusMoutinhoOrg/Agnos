# Templating

## Description
Index of the documentation for people using Agnos-Cli as a **template**: turning this repository into their own CLI, or converting an existing library to the same dependency-injection structure. Contributing to Agnos-Cli itself is indexed by [Development.md](/docs/Index/Development.md).

Pick **one** of the two workflows and follow it end to end. Both are phased step lists, and both take each file's fate from the same per-file action table.

---

## Tutorials

- [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md)
  - **description:** **Start here for a new library**: use this repo as a GitHub template
  - [Phase 1 — Create the repository](/docs/Tutorials/ForkTemplate.md#phase-1--create-the-repository)
  - [Phase 2 — Rewrite the contracts](/docs/Tutorials/ForkTemplate.md#phase-2--rewrite-the-contracts)
  - [Phase 3 — Create the implementation](/docs/Tutorials/ForkTemplate.md#phase-3--create-the-implementation)
  - [Phase 4 — Rewrite the documentation](/docs/Tutorials/ForkTemplate.md#phase-4--rewrite-the-documentation)
  - [Phase 5 — Verify](/docs/Tutorials/ForkTemplate.md#phase-5--verify)
- [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md)
  - **description:** **Start here for an existing library**: convert it to this DI structure
  - [Phase 1 — Bring the structure in](/docs/Tutorials/AdaptExistingLib.md#phase-1--bring-the-structure-in)
  - [Phase 2 — Rewrite the contracts](/docs/Tutorials/AdaptExistingLib.md#phase-2--rewrite-the-contracts)
  - [Phase 3 — Move the code into the sandbox](/docs/Tutorials/AdaptExistingLib.md#phase-3--move-the-code-into-the-sandbox)
  - [Phase 4 — Rewrite the documentation](/docs/Tutorials/AdaptExistingLib.md#phase-4--rewrite-the-documentation)
  - [Phase 5 — Verify](/docs/Tutorials/AdaptExistingLib.md#phase-5--verify)
- [RenameModule.md](/docs/Tutorials/RenameModule.md)
  - **description:** Rename the Go module path and update every internal import

---

## References

- [TemplateFileActions.md](/docs/References/TemplateFileActions.md)
  - **description:** The per-file action both workflows follow: copy, create, rewrite, or delete
  - [Copy](/docs/References/TemplateFileActions.md#copy)
  - [Rewrite](/docs/References/TemplateFileActions.md#rewrite)
  - [Create](/docs/References/TemplateFileActions.md#create)
  - [Delete](/docs/References/TemplateFileActions.md#delete)

The `bootstrap/` tree — a library built from this template embedding another one — is explained in [Bootstrap.md](/docs/References/Bootstrap.md).
