# Handle Library Samples

## Description
Covers creating and running executable Go samples in [examples/libraryExamples/](/examples/libraryExamples/) that demonstrate library features — typically after adding one through [HandleLibElements.md](/docs/Development/Tutorials/HandleLibElements.md). The shell-script counterparts, demonstrating features against the built CLI, are covered by [HandleCliExamples.md](/docs/Development/Tutorials/HandleCliExamples.md).

---

## Run a Library Sample

### Workflow
1. Browse the [examples/libraryExamples/](/examples/libraryExamples/) directory and pick a sample (e.g., `TrackSpendSample/`).
2. Run it from the project root with the Go toolchain:
   ```bash
   go run ./examples/libraryExamples/TrackSpendSample/TrackSpendSample.go
   ```
3. Pass arguments after the file when the sample takes them — `MainCallSample` runs the whole CLI, so it takes the same command line the installed binary does:
   ```bash
   go run ./examples/libraryExamples/MainCallSample/MainCallSample.go category list
   ```

---

## Add a Library Sample

### Rules
- Creating a sample requires updating [ApiSamplesList.md](/docs/LibUsage/References/ApiSamplesList.md) and [Structure.md](/docs/Development/References/Structure.md).
- A sample must be self-contained and runnable with a single `go run` command.
- The sample file must follow its specification — locate it in [Specs.md](/docs/Development/References/Specs.md).

### Workflow
1. Create a directory inside [examples/libraryExamples/](/examples/libraryExamples/) named after the feature being demonstrated (e.g., `examples/libraryExamples/NewFeatureSample/`).
2. Inside it, create the sample file with the same name as the directory (e.g., `NewFeatureSample.go`).
3. Write a runnable `package main` program that builds deps through an adapter, injects them into the lib, and uses the feature. Comment the key parts.
4. If the sample needs setup instructions, add a `README.md` in the sample's directory.
5. Add the sample to [ApiSamplesList.md](/docs/LibUsage/References/ApiSamplesList.md).
6. Register the new directory and file in [Structure.md](/docs/Development/References/Structure.md).
7. Verify the sample runs.
