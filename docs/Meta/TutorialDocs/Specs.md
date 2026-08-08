# TutorialDocs Specification

## Description
Defines the required shape of a **Tutorial** page — any `.md` file under `docs/Tutorials/`. A tutorial is a workflow guide for a single goal (e.g. `AddSample.md`, `AddDependency.md`), built from actionable numbered steps.

### Rules
- Every page must comply with [GeneralDoc](/docs/Meta/GeneralDoc/Specs.md).
- **One goal per page**: a tutorial has exactly one `## Workflow`. A guide covering more than one goal must be split into one page per goal, cross-linked by their steps.
- The title names that single goal as an action (e.g. `Add a Sample`, not `Handle Samples`).
- A tutorial page must contain a `## Description` and a `## Workflow` with **actionable**, numbered steps.
- An optional `### Rules` section states constraints specific to that tutorial.
- Steps prescribe actions, not descriptions; use fenced code blocks when a step involves writing or running code.
- **Full Code at the end**: when the workflow builds a program across multiple steps (code snippets that the reader assembles into one file), the page must end with a `## Full Code` section containing the complete, copy-pasteable program the steps produce. A tutorial whose code already appears whole in a single step is exempt.
- Background explanations belong in `docs/Explanations/` — link to them instead of embedding them.
- Every new page must be registered in the [README.md](/README.md) Doc Index.

## Structure
1. **Title** (H1): the goal, phrased as an action.
2. **`## Description`**: one short paragraph on what the tutorial covers, linking to the neighboring tutorials it does *not* cover.
3. **`### Rules`** *(optional)*: constraints for this tutorial.
4. **`---`**: horizontal rule separating the header from the workflow.
5. **`## Workflow`**: numbered, actionable steps, with fenced code blocks where a step involves code, and links to other tutorials for any step that is itself a separate goal.
6. **`## Full Code`** *(required when the workflow assembles a program across steps)*: the complete resulting code in a single fenced block, ready to copy and run — see [TrackTransactions.md](/docs/TrackTransactions.md) for a concrete example.

> **Note**: For a concrete example, refer to [sample.md](/docs/Meta/TutorialDocs/sample.md).
