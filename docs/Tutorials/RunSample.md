# Run a Library Sample

## Description
Covers running the executable Go samples in [libraryExamples/](../../libraryExamples/) to see the library in action. To run the shell scripts that drive the built CLI instead, follow [RunCliExample.md](/docs/Tutorials/RunCliExample.md).

---

## Workflow
1. Browse the [libraryExamples/](../../libraryExamples/) directory and pick a sample (e.g., `TrackSpendSample/`).
2. Run it from the project root with the Go toolchain:
   ```bash
   go run ./libraryExamples/TrackSpendSample/TrackSpendSample.go
   ```
3. Pass arguments after the file when the sample takes them — `MainCallSample` runs the whole CLI, so it takes the same command line the installed binary does:
   ```bash
   go run ./libraryExamples/MainCallSample/MainCallSample.go category list
   ```
