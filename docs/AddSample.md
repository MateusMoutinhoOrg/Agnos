# Add a Library Sample

## Description
Covers creating a runnable Go sample in [libraryExamples/](../../libraryExamples/) that demonstrates a library feature. The shell-script counterpart, demonstrating a feature against the built CLI, is covered by [AddCliExample.md](/docs/AddCliExample.md).

### Rules
- Creating a sample requires updating the Library Examples section of the [README.md](/README.md) and [Structure.md](/docs/Structure.md).
- A sample must be self-contained and runnable with a single `go run` command.
- The sample file must follow its specification — locate it in [Specs.md](/docs/Specs.md).

---

## Workflow
1. Create a directory inside [libraryExamples/](../../libraryExamples/) named after the feature being demonstrated (e.g., `libraryExamples/NewFeatureSample/`).
2. Inside it, create the sample file with the same name as the directory (e.g., `NewFeatureSample.go`).
3. Write a runnable `package main` program that builds deps through an adapter, injects them into the lib, and uses the feature. Comment the key parts.
4. If the sample needs setup instructions, add a `README.md` in the sample's directory.
5. Add the sample to the Library Examples section of the [README.md](/README.md).
6. Register the new directory and file in [Structure.md](/docs/Structure.md).
7. Verify the sample runs, following [RunSample.md](/docs/RunSample.md).
