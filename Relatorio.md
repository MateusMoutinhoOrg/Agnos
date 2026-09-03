# Agnos Framework: Developer Experience (DX) and Usability Report

## 1. Executive Summary

Agnos is a highly promising CLI builder framework that successfully tackles one of the most painful aspects of Command Line Interface (CLI) development: boilerplate code. By enforcing a declarative approach using YAML configurations, Agnos allows developers (and AI agents) to focus strictly on the business logic of their commands.

However, a framework's success relies heavily on its Developer Experience (DX), predictability, and out-of-the-box stability. During the evaluation of Agnos, several critical bugs and usability pain points were identified. 

This report provides a high-level, conceptual analysis of these issues. It focuses on how these bugs impact the user experience and outlines conceptual recommendations for how the framework's behavior and features should evolve in future versions, without delving into internal codebase mechanics.

---

## 2. Strengths and Core Value Proposition

Before addressing the areas for improvement, it is important to recognize the conceptual strengths of Agnos that make it a compelling tool.

### 2.1. Declarative Design Paradigm
The decision to separate the CLI interface configuration (`entries.yaml`) from the business logic (`handler.go`) is conceptually brilliant. It forces developers into a clean architecture. By looking at a single YAML file, a developer can instantly understand a command's purpose, expected arguments, flags, and types.

### 2.2. Boilerplate Elimination
By abstracting the flag parsing and argument binding processes, Agnos significantly lowers the barrier to entry for building robust CLIs. Developers do not need to learn complex setup functions; they simply define their schema and write their logic.

### 2.3. AI-Friendly Architecture
For Large Language Models (LLMs) assisting developers, Agnos is an ideal target. AI agents thrive in strictly structured environments. The predictability of the directory structure and the simplicity of YAML generation mean that LLMs can scaffold Agnos commands with an incredibly low error rate.

---

## 3. Identified Bugs and Usability Issues

Despite its strengths, the current iteration of Agnos contains bugs and structural decisions that severely disrupt the developer workflow. 

### 3.1. Bug: Inconsistent Directory Creation in `agnos start`

**The Scenario:**
When a developer attempts to initialize a new project in a specific directory from their current workspace by running:
`agnos start -p project_name --path ./target_directory`

**The Expected Behavior:**
All project files, configurations, and internal folders (`sandbox/`, `AgnosConfig/`, `cmd/`, `adapters/`, and `go.mod`) should be contained entirely within the `./target_directory`.

**The Actual Behavior (Bug):**
The framework scatters the files. The `go.mod` file is correctly placed inside the `./target_directory`, but the core Agnos folders (`sandbox/`, `AgnosConfig/`, `cmd/`, `adapters/`) are erroneously generated in the Current Working Directory (the parent folder). 

**The Impact on Developer Experience (DX):**
This is a highly frustrating bug. It pollutes the user's workspace, leading to fragmented projects that fail to build. The developer is forced to manually clean up the scattered folders, navigate into the target directory, and run the command again without the `--path` flag. This breaks trust in the tooling right at the first step of adoption.

### 3.2. Usability Issue: Non-Deterministic Positional Arguments

**The Scenario:**
In CLI tools, the order of positional arguments is sacred. If a user runs a subtraction command `sub 10 5`, the tool must deterministically map `10` to the first argument and `5` to the second argument. 

In Agnos, arguments are defined in `entries.yaml` using a YAML Map (a key-value dictionary):
```yaml
args:
  first_number:
    type: string
  second_number:
    type: string
```

**The Actual Behavior:**
YAML maps are inherently unordered structures. When Agnos processes this configuration, there is no guarantee whether `first_number` or `second_number` will be parsed first. Consequently, the CLI might bind `5` to `first_number` and `10` to `second_number`, completely inverting the logic.

**The Impact on Developer Experience (DX):**
This unpredictability makes positional arguments entirely unreliable. Developers are forced to rely on "hacks," such as naming their arguments alphabetically (`arg1`, `arg2`), hoping the underlying parser iterates them in alphabetical order. This breaks the declarative promise of the framework and introduces subtle, hard-to-track bugs in mathematical or data-processing commands.

### 3.3. Usability Issue: Manual Type Casting (Loss of YAML Type Value)

**The Scenario:**
Agnos allows developers to declare the data type of an argument directly in the YAML file:
```yaml
args:
  age:
    type: integer
```

**The Expected Behavior:**
Because the framework knows the argument is an `integer`, the developer expects the framework to hand them an integer when their command logic executes.

**The Actual Behavior:**
Agnos ignores the declared YAML type when passing the data to the command logic. Everything is passed as a generic string. 

**The Impact on Developer Experience (DX):**
The developer is forced to write repetitive, manual data conversion code (e.g., parsing a string to an integer, checking for conversion errors, handling invalid formats) for every single numeric or boolean argument. This defeats the purpose of declaring the type in the YAML file. The YAML `type` field becomes purely cosmetic, offering no real runtime utility.

### 3.4. Usability Issue: Lack of Scaffolding CLI Commands

**The Scenario:**
Agnos enforces a strict directory structure. To add a single new command to the CLI, a developer must manually create a folder under `sandbox/internal/commands/`, create an `entries.yaml` file, remember the exact schema for the YAML file, create a `handler.go` file, and remember the exact function signature required by the framework.

**The Expected Behavior:**
A modern CLI builder should provide tooling to automate its own boilerplate. A developer should be able to run a simple command to scaffold a new feature.

**The Impact on Developer Experience (DX):**
The manual process is tedious, error-prone, and slow. Developers spend more time creating folders and copying boilerplate templates than they do writing their actual command logic.

---

## 4. Proposed Conceptual Changes and Enhancements

To address the issues outlined above and evolve Agnos into a mature, production-grade framework, the following conceptual changes should be prioritized.

### 4.1. Fix Scoping in Initialization Commands
The `--path` flag must act as an absolute boundary. The framework must guarantee that no file or folder generation ever escapes the directory specified by `--path`. If a target directory is specified, Agnos must treat that directory as the absolute root for all generation processes.

### 4.2. Shift to Array-Based Argument Configuration
To solve the non-deterministic argument ordering, Agnos must stop using YAML maps for positional arguments. Instead, the schema should require developers to define arguments as a List/Array. 

**Conceptual Fix:**
```yaml
args:
  - name: first_number
    type: string
  - name: second_number
    type: string
```
Because lists are ordered by definition, the framework will explicitly know that `first_number` is always index 0, and `second_number` is always index 1. This guarantees 100% predictability.

### 4.3. Implement Native Type Binding and Validation
The framework must honor the types declared in the YAML file. If a developer defines an argument as an `integer` or `float`, Agnos itself should handle the string-to-number conversion. 

Furthermore, Agnos should handle validation failures automatically. If a user runs a command and passes the word "apple" to a field defined as an `integer`, Agnos should intercept this, print a clean "Invalid Type" error to the console, and halt execution. The command logic should only ever be executed if all data inputs perfectly match the declared YAML schema.

### 4.4. Introduce a Robust Command Scaffolding Tool
Agnos should introduce an `add command` utility (e.g., `agnos add command my_feature`). 
When executed, this tool should:
1. Automatically create the necessary folder structure.
2. Generate a standard `entries.yaml` with placeholder values.
3. Generate a standard `handler.go` ready for business logic.

This single feature would exponentially speed up the development lifecycle and drastically improve the Developer Experience.

### 4.5. Introduce Advanced YAML Validation Features (Future Roadmap)
Looking ahead, Agnos should expand its declarative power by allowing developers to define complex validation rules directly in the YAML file. 
For example:
- `required: true` (Already present but should be strictly enforced)
- `min_value: 0` / `max_value: 100` (For numbers)
- `regex: "^[a-z]+$"` (For strings)

By offloading these validations to the framework, developers would never have to write "If X is less than 0, return error" inside their command logic. The framework would guarantee that the data is clean and valid before the command ever runs.

---

## 5. Conclusion

Agnos is built on a very solid philosophical foundation. The separation of interface configuration and business logic is a pattern that yields highly maintainable code. 

However, the current bugs related to project generation (`--path` scoping), structural flaws in YAML parsing (non-deterministic maps), and a lack of built-in conveniences (automatic type casting and scaffolding tools) hold the framework back from its true potential.

By addressing these high-level DX issues, Agnos can transition from a promising concept into a highly competitive, beloved tool for building Command Line Interfaces.
