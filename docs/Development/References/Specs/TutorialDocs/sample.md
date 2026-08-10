# Add a Sample

## Description
Covers creating a runnable sample in [examples/libraryExamples/](/examples/libraryExamples/) that demonstrates a library feature.

### Rules
- Creating a sample requires updating the Samples section of the [README.md](/README.md) and [Structure.md](/docs/Development/References/Structure.md).
- A sample must be self-contained and runnable with a single `go run` command.
- The sample file must follow its specification — locate it in [Specs.md](/docs/Development/References/Specs.md).

---

## Workflow
1. Create a directory inside [examples/libraryExamples/](/examples/libraryExamples/) named after the feature being demonstrated (e.g., `examples/libraryExamples/NewFeatureSample/`).
2. Inside it, create the sample file with the same name as the directory (e.g., `NewFeatureSample.go`).
3. Write a runnable `package main` program that builds deps through an adapter, injects them into the lib, and uses the feature. Comment the key parts.
4. If the sample needs setup instructions, add a `README.md` in the sample's directory.
5. Add the sample to the Samples section of the [README.md](/README.md).
6. Register the new directory and file in [Structure.md](/docs/Development/References/Structure.md).
7. Verify the sample runs, following [HandleSamples.md](/docs/Development/Protocols/HandleSamples.md).
