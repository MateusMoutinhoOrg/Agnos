# Handle Assets

## Description
Covers adding an asset — a file under [/assets/](/assets/) compiled into the binary and read at runtime through the injected [`EmbedDeps`](/docs/References/PublicApi/embeddeps.Lib.md) contract — and reading or listing one from library code. An asset is a payload better kept as a file than as a Go constant: a template, a long-form document, an image. Short display text is **not** an asset; it lives in [sandbox/config/](/sandbox/config/) and is covered by [HandleCliCommands.md](/docs/Tutorials/HandleCliCommands.md).

### Rules
- An asset is reached by path through `l.Deps.EmbedDeps`, never by importing the `assets` package — that package may only be imported from outside the sandbox, by an adapter.
- Every file under `assets/` is embedded by the single `//go:embed all:*` directive in `assets/asset.go`, wherever in the tree it sits. There is no pattern to keep in sync, and no way to add an asset the binary then cannot find.
- Paths are slash-separated and relative to the asset root, whatever separator the host operating system uses. The root itself is addressed as `"."`.
- The contract is **read-only**: assets ship with the program, and nothing in the library ever writes one back. Writing a file at runtime is [`IoLib`](/docs/References/PublicApi/iodeps.Lib.md), not this.
- A missing asset is a **packaging** mistake, not a user mistake: it surfaces at runtime as the error from `ReadFile`, never at build time. Report it by path rather than printing nothing.
- Adding, renaming, or deleting an asset changes the project structure — update [Structure.md](/docs/References/Structure.md#assets) in the same commit.

> [!NOTE]
> The asset tree is **empty** in this repository. The tracker's display text is short and fixed, so it lives in `sandbox/config` as compile-time constants, where the compiler checks every reference. What ships here is the mechanic — contract, adapter, and directive — wired end to end and ready for a library derived from this template. Nothing under `sandbox/` currently reads an asset.

---

## Add Asset

### Workflow
1. Write the file under [assets/](/assets/), in a directory named after what the files in it are:
   ```bash
   mkdir -p assets/templates
   cat > assets/templates/report.tmpl <<'EOF'
   Report for {{.Category}}
   ========================
   EOF
   ```
   Nothing else has to be embedded: `//go:embed all:*` in `assets/asset.go` already covers the whole directory.
2. Name the path in [sandbox/config/cli.go](/sandbox/config/cli.go), beside the other constants, so call sites reference the constant rather than a literal:
   ```go
   // ReportTemplateAsset is the report layout rendered by the `report` command.
   ReportTemplateAsset = "templates/report.tmpl"
   ```
3. Read it where the library needs it, through the injected contract — see [Retrieve Asset in Runtime](#retrieve-asset-in-runtime) below.
4. Add the file to the `/assets/` table in [Structure.md](/docs/References/Structure.md#assets).
5. Build and run the command that reads it — a path typo surfaces at runtime as a `ReadFile` error, never at build time:
   ```bash
   go build ./... && AGNOS_DATA=./scratch go run ./cmd/main report
   ```

---

## ListAssets in Runtime

The library lists embedded assets through `l.Deps.EmbedDeps.ListFiles` or `l.Deps.EmbedDeps.ListFilesRecursively`. Both return paths **relative to the directory asked for**, in lexical order, and never report directories.

1. To list the files sitting directly inside a directory, use `ListFiles`:
   ```go
   // Names are relative to "templates": "report.tmpl", not
   // "templates/report.tmpl". Nested directories are not descended into.
   names, err := l.Deps.EmbedDeps.ListFiles("templates")
   if err != nil {
       return Failure(l, config.AssetUnreadable, "templates")
   }
   ```

2. To list every file at or below a directory, use `ListFilesRecursively`:
   ```go
   // Paths keep the directories below the root asked for:
   // "templates/report.tmpl". The root itself is ".".
   paths, err := l.Deps.EmbedDeps.ListFilesRecursively(".")
   if err != nil {
       return Failure(l, config.AssetUnreadable, ".")
   }
   ```

---

## Retrieve Asset in Runtime

The library reads one asset's content through `l.Deps.EmbedDeps.ReadFile`, passing the slash-separated path relative to the asset root.

1. Read the bytes and report the failure by path, so a packaging mistake is legible:
   ```go
   content, err := l.Deps.EmbedDeps.ReadFile(config.ReportTemplateAsset)
   if err != nil {
       // The asset is missing from the binary — report which one, and exit
       // ExitFailure: the command line was fine, the packaging was not.
       return Failure(l, config.AssetUnreadable, config.ReportTemplateAsset)
   }
   l.Deps.Printf("%s", string(content))
   ```
2. Add any message the step above prints to [sandbox/config/cli.go](/sandbox/config/cli.go), following [HandleCliCommands.md](/docs/Tutorials/HandleCliCommands.md):
   ```go
   AssetUnreadable = `missing asset "%s" — the binary was built without it`
   ```
