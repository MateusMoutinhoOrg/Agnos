# Templating

## Description
Index of the documentation for people using Agnos-Cli as a **template**: turning this repository into their own CLI, or converting an existing library to the same dependency-injection structure. Contributing to Agnos-Cli itself is indexed by [Development.md](/docs/Index/Development.md).

Pick **one** of the two workflows and follow it end to end. Both are phased step lists, and both take each file's fate from the same per-file action table.

---

## Tutorials

| Doc | Description |
| --- | --- |
| [ForkTemplate.md](/docs/Tutorials/ForkTemplate.md) | **Start here for a new library**: use this repo as a GitHub template |
| [AdaptExistingLib.md](/docs/Tutorials/AdaptExistingLib.md) | **Start here for an existing library**: convert it to this DI structure |
| [RenameModule.md](/docs/Tutorials/RenameModule.md) | Rename the Go module path and update every internal import |

---

## References

| Doc | Description |
| --- | --- |
| [TemplateFileActions.md](/docs/References/TemplateFileActions.md) | The per-file action both workflows follow: copy, create, rewrite, or delete |

The `bootstrap/` tree — a library built from this template embedding another one — is explained in [Bootstrap.md](/docs/References/Bootstrap.md).
