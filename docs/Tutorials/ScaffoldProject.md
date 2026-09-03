# Scaffold a Project

## Description
Covers taking an empty directory to a compiling Go CLI with one command of your own, using nothing but `agnos` commands: `start` writes the configuration and the sandbox skeleton, `deps-init` and `dep-install` bring the dependency-injection layer in, `cli-init` installs the command-line layer, and `add-command` declares a command. Installing `agnos` is [InstallCli.md](/docs/Tutorials/InstallCli.md); growing the command afterwards is [ShapeCommands.md](/docs/Tutorials/ShapeCommands.md); writing its logic is [WriteCommandHandler.md](/docs/Tutorials/WriteCommandHandler.md).

### Rules
- Every command below takes the project directory through `--path` and defaults to the current directory. A stray positional (`agnos build .`) is a usage error, not an ignored argument — see [Commands.md](/docs/References/Commands.md#exit-codes).
- Every step that adds something ends by running `build`, which renders the generated files and then runs `go mod tidy` and `go build` over the project. A step that prints a Go compile error did not finish: fix the project before going on.
- The Go toolchain must be on your `PATH`: `agnos` asks it for its version when writing `go.mod` and hands every build to it.

---

## Workflow
1. Create the directory and scaffold the project into it. `--project-name` is required, and `--module` is required whenever the directory has no `go.mod` yet:
   ```bash
   mkdir my-tool && cd my-tool
   agnos start --project-name my-tool --module github.com/you/my-tool
   ```
   The tree now holds `AgnosConfig/` (the four configuration files), `go.mod`, and the sandbox skeleton — `sandbox/new.go`, `sandbox/api/sandbox.go`, `sandbox/internal/config/config.go` — and already compiles. Every file `agnos` writes is listed in [GeneratedFiles.md](/docs/References/GeneratedFiles.md).
2. Install the dependency-injection layer. `deps-init` creates `sandbox/deps/` and `adapters/`; the follow-up build then renders the `Deps` struct and the `standard` adapter over whatever those directories hold:
   ```bash
   agnos deps-init
   ```
3. Pick the capabilities the project needs from the installable deps and install them one by one. Each dep drops a contract under `sandbox/deps/<dep>/` and its implementation under `adapters/libs/`, then rebuilds so `Deps` and `standard.New` grow one field and one `Bind` call:
   ```bash
   agnos dep-list
   agnos dep-install iodeps
   ```
   Every dep and what it brings is listed in [DepList.md](/docs/References/DepList.md); managing them afterwards is [ManageDeps.md](/docs/Tutorials/ManageDeps.md).
4. Install the command-line layer. `cli-init` installs the two deps the layer needs (`std` and `argvdeps`), renders `cmd/main/main.go`, the `Cli` contract, and the generated `help` and `version` commands, then rebuilds:
   ```bash
   agnos cli-init
   ```
   The project is now a runnable CLI — its binary name is the lowercased project name:
   ```bash
   go run ./cmd/main help
   go run ./cmd/main version
   ```
5. Declare your first command. Both `--help` and `--category` are required; the name must normalize to `[a-z][a-z0-9-]*`, and `help` is refused because `agnos` generates that command:
   ```bash
   agnos add-command greet --help "Say hello" --category Demo
   ```
   This writes `sandbox/internal/commands/greet/entries.yaml` (the declaration) and `handler.go` (a stub), and the rebuild generates `entries.go` and a dispatch arm for it:
   ```bash
   go run ./cmd/main greet     # prints "greet called"
   go run ./cmd/main help greet
   ```
6. Give the command a flag and an argument from the command line — never by editing YAML — following [ShapeCommands.md](/docs/Tutorials/ShapeCommands.md):
   ```bash
   agnos add-flag name --command greet --identifier --name --identifier -n --type string --required --description "who to greet"
   agnos add-arg times --command greet --type int --min 1 --default 1 --description "how many times"
   ```
7. Open `sandbox/internal/commands/greet/handler.go` and replace the stub with the command's logic, following [WriteCommandHandler.md](/docs/Tutorials/WriteCommandHandler.md). `entries.Name` and `entries.Times` are already typed and validated by the time the handler runs.
8. Regenerate and check the whole project whenever you have edited it by hand:
   ```bash
   agnos build
   ```
   `build` runs the schema check first — see [RegenerateProject.md](/docs/Tutorials/RegenerateProject.md) — then re-renders every generated file and compiles the result.
9. Add `--quiet` (or `-q`) to any of the commands above when the progress lines are noise, for instance inside a script. The result on standard output and any error on standard error are never silenced — see [Commands.md](/docs/References/Commands.md#output-channels).
