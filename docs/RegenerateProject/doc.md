# Regenerate and Check a Project

## Description
Covers the two commands that keep a scaffolded project honest: `verify`, which checks the tree against the harness schema and writes nothing, and `build`, which runs that check, re-renders every generated file, and hands the result to the Go toolchain. What is checked is explained in [BuildPipeline](/docs/BuildPipeline/doc.md#the-verify-gate); which files are rewritten is listed in [GeneratedFiles](/docs/GeneratedFiles/doc.md).

### Rules
- Every other `agnos` command that changes a project ends by running `build` for you, so these two are for after **hand** edits — a new `handler.go`, a contract you wrote yourself, a moved file.
- `build` refuses to render over a tree that fails `verify`. `--unsafe` skips the gate for a mid-refactor state you intend to fix; it never skips the runtime.
- A `build` or `verify` that reports success has been accepted by `go build`, unless you passed `--runtime none`.

---

## Workflow
1. Check the schema without touching a file. Every violation is listed at once, prefixed by the path that caused it, and the command exits `1`:
   ```bash
   agnos verify
   ```
   ```text
   verify found 1 violation(s):
     - sandbox/ contains unexpected directory extra (allowed: api, binds, deps, internal)
   ```
2. Regenerate the project. `build` runs `verify`, re-renders the `all` asset group, the `deps` group when `sandbox/deps/` exists and the `cli` group when `sandbox/internal/cli/` exists, persists, then runs `go mod tidy` and `go build` over `./cmd/...`, `./sandbox/...` and `./adapters/...` — whichever of those exist:
   ```bash
   agnos build
   ```
3. Render only, when the toolchain is not the point — no Go installed on this machine, or hand-written code you know does not compile yet:
   ```bash
   agnos build --runtime none
   ```
4. Skip the schema gate while restructuring, then run `verify` before committing:
   ```bash
   agnos build --unsafe
   agnos verify
   ```
5. Drive both from a script by their exit codes — `0` accepted, `1` a violation or a compile error, `2` a wrong command line — and keep standard output clean with `--quiet`:
   ```bash
   agnos build -q || echo "build failed with $?"
   ```
6. Point either command at another directory with `--path`. Every read and write is scoped to that directory; no generation ever escapes it:
   ```bash
   agnos build --path ./my-tool
   ```
