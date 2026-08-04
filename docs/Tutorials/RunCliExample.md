# Run a CLI Example

## Description
Covers running the shell scripts in [cliExamples/](/cliExamples/), each of which builds the CLI and drives it the way a user would. Writing a new one is covered by [AddCliExample.md](/docs/Tutorials/AddCliExample.md); the Go programs that use the library from code instead are covered by [RunSample.md](/docs/Tutorials/RunSample.md).

### Rules
- A CLI example needs nothing installed: it builds the binary itself.
- It writes to a scratch directory it removes on exit, so running one never touches the records in your home directory.

---

## Workflow
1. Browse [cliExamples/](/cliExamples/) and pick a script — each is named after the goal it demonstrates, so `ManageCategories.sh` is a good starting point.
2. Run it from the project root:
   ```bash
   bash ./cliExamples/ManageCategories.sh
   ```
3. Read the transcript alongside the script: each `== …` line in the output is the comment above the commands that produced what follows it.
4. Run the rest in order to see the whole interface:
   ```bash
   for script in ./cliExamples/*.sh; do bash "$script"; done
   ```
5. Try the same commands against your own budget once you have installed the binary, following [InstallCli.md](/docs/Tutorials/InstallCli.md) and [UseCli.md](/docs/Tutorials/UseCli.md).
