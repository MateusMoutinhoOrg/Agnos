# Handle Assets

## Description
Covers adding, editing, and reading an asset — a file under [/assets/](/assets/) that the library reads through the injected [`EmbedDeps`](/docs/References/PublicApi/embeddeps.Lib.md) contract instead of holding as a Go string. Every word the command-line interface prints is one. Why the mechanic exists is explained in [EmbeddedAssets.md](/docs/References/EmbeddedAssets.md); adding the command that prints a new message is [HandleCliCommands.md](/docs/Tutorials/HandleCliCommands.md).

### Rules
- Display text is **never** written in `sandbox/`. A string literal in the sandbox is a format skeleton (`"%s\n"`) or a diagnostic, never a sentence shown to a user.
- An asset is reached by path through `l.Deps.EmbedDeps`, never by importing the `assets` package — that package may only be imported from outside the sandbox, by an adapter.
- Every file under `assets/` is embedded by the single `//go:embed all:*` directive in `assets/asset.go`, wherever in the tree it sits. There is no pattern to keep in sync, and no way to add an asset the binary then cannot find.
- A message file is a `Printf` format read through `message`: it carries the verbs, and the quotes, its values are rendered with. It must not end in a newline the caller also adds — `message` trims the trailing one.
- Adding, renaming, or deleting an asset changes the project structure — update [Structure.md](/docs/References/Structure.md) in the same commit.

---

## Add Asset

### Workflow
1. Write the file under [assets/](/assets/), named after what it says — one printable line per file under `assets/messages/`:
   ```bash
   echo 'largest transaction: %s' > assets/messages/largest-transaction.txt
   ```
   Nothing else has to be embedded: `//go:embed all:*` in `assets/asset.go` already covers the whole directory.
2. Name it in `sandbox/config/cli.go`, beside the other message constants, so call sites reference the constant rather than the path:
   ```go
   // LargestTransactionMessage reports the biggest single transaction.
   LargestTransactionMessage = "largest-transaction"
   ```
3. Add the file to the `/assets/` table in [Structure.md](/docs/References/Structure.md#assets).
4. Build and run the command that prints it — a path typo surfaces at runtime as `agnos-cli: missing asset …`, never at build time:
   ```bash
   go build ./... && AGNOS_DATA=./scratch go run ./cmd/main largest
   ```

---

## ListAssets in Runtime

The library lists embedded assets through `l.Deps.EmbedDeps.ListFiles` or `l.Deps.EmbedDeps.ListFilesRecursively`.

1. To list files directly inside a directory, use `ListFiles`:
   ```go
   names, err := l.Deps.EmbedDeps.ListFiles("messages")
   if err != nil {
       // Handle error (e.g., missing directory)
   }
   // names will contain "largest-transaction.txt", etc.
   ```

2. To list all files at or below a directory, use `ListFilesRecursively`:
   ```go
   paths, err := l.Deps.EmbedDeps.ListFilesRecursively(".")
   if err != nil {
       // Handle error
   }
   // paths will contain "messages/largest-transaction.txt", "usages.txt", etc.
   ```

---

## Retrieve Asset in Runtime

The library retrieves the contents of an embedded asset using `l.Deps.EmbedDeps.ReadFile`.

1. Pass the slash-separated path relative to the asset tree root:
   ```go
   content, err := l.Deps.EmbedDeps.ReadFile("version.txt")
   if err != nil {
       // Handle error (e.g., asset not found)
       return err
   }
   l.Deps.Printf("%s\n", string(content))
   ```
