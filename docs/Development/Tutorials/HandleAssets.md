# Handle Assets

## Description
Covers adding, editing, and reading an asset — a file under [/assets/](/assets/) that the library reads through the injected [`EmbedDeps`](/docs/LibUsage/References/PublicApi/embeddeps.Lib.md) contract instead of holding as a Go string. Every word the command-line interface prints is one. Why the mechanic exists is explained in [EmbeddedAssets.md](/docs/LibUsage/References/EmbeddedAssets.md); adding the command that prints a new message is [HandleCliCommands.md](/docs/Development/Tutorials/HandleCliCommands.md).

### Rules
- Display text is **never** written in `sandbox/`. A string literal in the sandbox is a format skeleton (`"%s\n"`) or a diagnostic, never a sentence shown to a user.
- An asset is reached by path through `l.Deps.EmbedDeps`, never by importing the `assets` package — that package may only be imported from outside the sandbox, by an adapter.
- Assets under `assets/cli/` are picked up by the `//go:embed all:cli` directive automatically. An asset added **outside** an already-embedded path needs its own pattern in `assets/assets.go`, or it will not exist at runtime — and nothing will fail at build time.
- A message file is a `Printf` format read through `message`: it carries the verbs, and the quotes, its values are rendered with. It must not end in a newline the caller also adds — `message` trims the trailing one.
- Adding, renaming, or deleting an asset changes the project structure — update [Structure.md](/docs/Development/References/Structure.md) in the same commit, per [RULES.md](/docs/Development/References/RULES.md#file-changes).

---

## Add an Asset

### Workflow
1. Write the file under [assets/](/assets/), named after what it says, in the directory its reader lives in — interface text goes under `assets/cli/`, one printable line per file under `assets/cli/messages/`:
   ```bash
   echo 'largest transaction: %s' > assets/cli/messages/largest-transaction.txt
   ```
2. Check the file is inside a pattern already embedded by `assets/assets.go`. It is, for anything under `assets/cli/`; for a new top-level file or directory, add its pattern:
   ```go
   //go:embed version.txt
   //go:embed all:cli
   //go:embed all:templates // ← a new asset directory
   var Files embed.FS
   ```
3. Name it in `sandbox/internal/cli/cli.go`, beside the other message constants, so call sites reference the constant rather than the path:
   ```go
   // LargestTransactionMessage reports the biggest single transaction.
   LargestTransactionMessage = "largest-transaction"
   ```
4. Read it where it is printed. `message` returns one line, trimmed of its trailing newline; `asset` returns a whole file verbatim:
   ```go
   l.Deps.Printf(message(l, LargestTransactionMessage)+"\n", biggest.String())
   ```
5. Add the file to the `/assets/` table in [Structure.md](/docs/Development/References/Structure.md#assets).
6. Build and run the command that prints it — a path typo surfaces at runtime as `agnos-cli: missing asset …`, never at build time:
   ```bash
   go build ./... && AGNOS_DATA=./scratch go run ./cmd/main largest
   ```

### Edit an Existing Asset

Rewording the interface, or translating it, touches no Go at all: edit the file and rebuild, since the assets are compiled into the binary.

```bash
$EDITOR assets/cli/usage.txt          # the help screen
$EDITOR assets/version.txt            # the version `agnos-cli version` reports
go build ./... && go run ./cmd/main --help
```

A release bump is `assets/version.txt` plus the `@v0.0.3` install tag pinned in the tutorials that document installing — see the version note in [Structure.md](/docs/Development/References/Structure.md#cmd).
