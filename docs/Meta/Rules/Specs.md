# Rules Specification

## Description
Defines the required shape of the contribution rules document at `docs/References/RULES.md` — the binding list of actions that force companion files to be updated in the same commit.

### Rules
- `RULES.md` must open by pointing to [Specs.md](/docs/Specs.md) as the authority every file must be shaped by.
- Every rule must be stated as a **trigger → required update**: "when you do X, update Y".
- Rules must be grouped by the kind of change (file changes, documentation changes, sample changes, dependency changes, …), one `##` section each.
- Each rule must link, with relative paths, to the file it requires updating.

## Structure
1. **Title** (H1): `# Contribution Rules`.
2. **Intro**: one line stating the rules are binding and deferring the shape of each file to `docs/References/Specs.md`.
3. **One `##` section per change type**, each separated by `---`, describing the trigger and the file(s) that must be kept in sync.

> **Note**: For a concrete example, refer to [sample.md](/docs/Meta/Rules/sample.md).
