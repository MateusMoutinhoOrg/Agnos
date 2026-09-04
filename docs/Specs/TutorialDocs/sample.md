# Add a Schema Check

## Description
Covers adding one rule to the `verify` gate — a `check_*.go` file under `sandbox/internal/actions/verify/`. What the gate enforces today is [BuildPipeline](/docs/BuildPipeline/doc.md#the-verify-gate).

### Rules
- A check reads through the open `*smartio.SmartIO` and writes nothing.
- A check returns `[]string`, one entry per violation, each naming the path that caused it.
- Adding a check requires updating the rule table of [BuildPipeline](/docs/BuildPipeline/doc.md#the-verify-gate).

---

## Workflow
1. Create `sandbox/internal/actions/verify/check_<rule>.go` with one exported function returning the violations:
   ```go
   // CheckCmd enforces that cmd/ holds only directories, one per executable.
   func CheckCmd(deps *deps.Deps, io *smartio.SmartIO) []string {
       var violations []string
       for _, file := range io.ListFiles("cmd") {
           violations = append(violations, "cmd/ contains loose file "+lastSegment(file))
       }
       return violations
   }
   ```
2. Append its result in `verify_internal.go`, beside the existing checks.
3. Rebuild with the bootstrap binary and run `verify` over this tree — it must pass.
4. Register the rule in [BuildPipeline](/docs/BuildPipeline/doc.md#the-verify-gate).
