# ExplanationDocs Specification

## Description
Defines the required shape of the `doc.md` of an **explanation** doc — a page explaining a mechanic. An explanation page describes how a mechanic, feature, or ability of the project works, building understanding rather than listing items or prescribing steps.

### Rules
- Every page must comply with [GeneralDoc](/docs/Specs/GeneralDoc/doc.md).
- Each `##` section must explain a single aspect of the topic.
- Code examples are illustrative: complete and runnable where possible, with inline comments highlighting what the section explains. Generated code is shown as generated, marked with a comment saying so.
- Explanation pages must not contain step-by-step procedures — link to the relevant tutorial instead.
- Every new doc declares itself in its `props.yaml` — see [DocProps](/docs/Specs/DocProps/doc.md); the theme indexes are generated from it.

## Structure
1. **Title** (H1): the mechanic or feature being explained.
2. **`## Description`**: one short paragraph on what the page explains.
3. **One `##` section per aspect**, separated by `---`, each combining prose with illustrative code examples where useful.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/ExplanationDocs/sample.md).
