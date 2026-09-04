# TutorialDocs Specification

## Description
Defines the required shape of the `doc.md` of a **workflow** doc — a page that guides the reader through one goal (e.g. `ScaffoldProject`, `AddAction`), built from actionable numbered steps. Which kind a doc is comes from the page, not from where it sits: every doc is a directory of `docs/`.

### Rules
- Every page must comply with [GeneralDoc](/docs/Specs/GeneralDoc/doc.md).
- **One goal per page**: a tutorial has exactly one `## Workflow`. A guide covering more than one goal must be split into one page per goal, cross-linked by their steps. A workflow with several named stages splits them into `###` subsections, which then become the page's topics in its index entry.
- The title names that single goal as an action (e.g. `Add an Action`, not `Handle Actions`).
- A tutorial page must contain a `## Description` and a `## Workflow` with **actionable**, numbered steps.
- An optional `### Rules` section states constraints specific to that tutorial.
- Steps prescribe actions, not descriptions; use fenced code blocks when a step involves writing or running code.
- **Full Code at the end**: when the workflow builds a program across multiple steps (code snippets that the reader assembles into one file), the page must end with a `## Full Code` section containing the complete, copy-pasteable program the steps produce. A tutorial whose code already appears whole in a single step is exempt.
- Background explanations belong in their own explanation page — link to them instead of embedding them.
- Every new doc declares itself in its `props.yaml` — see [DocProps](/docs/Specs/DocProps/doc.md); the theme indexes are generated from it.

## Structure
1. **Title** (H1): the goal, phrased as an action.
2. **`## Description`**: one short paragraph on what the tutorial covers, linking to the neighboring tutorials it does *not* cover.
3. **`### Rules`** *(optional)*: constraints for this tutorial.
4. **`---`**: horizontal rule separating the header from the workflow.
5. **`## Workflow`**: numbered, actionable steps, with fenced code blocks where a step involves code, and links to other tutorials for any step that is itself a separate goal.
6. **`## Full Code`** *(required when the workflow assembles a program across steps)*: the complete resulting code in a single fenced block, ready to copy and run — see [WriteCommandHandler](/docs/WriteCommandHandler/doc.md) for a concrete example.

> **Note**: For a concrete example, refer to [sample](/docs/Specs/TutorialDocs/sample.md).
