package assets

// The project's embedded assets: the text, and later the images and item
// templates, the library serves through the injected embed contract
// (deps.Deps.EmbedDeps). Keeping them in files instead of in Go string
// constants is what lets the sandbox hold no display text of its own — it
// asks for an asset by path and prints whatever comes back.
//
// This package exists for one reason: a //go:embed directive can only reach
// files inside its own package directory, so the directive has to live next
// to the assets themselves. It holds no behavior and no state beyond the
// embedded filesystem, and only code outside the sandbox may import it — the
// standard adapter does, in adapters/standard/embed.go, and wraps it into the
// embeddeps.Lib contract the sandbox reads through.

import "embed"

// Files is every asset shipped with the project, compiled into the binary, so
// an installed `agnos-cli` finds its help screen and its version with no files on
// disk next to it. Paths inside it are slash-separated and rooted here:
// "version.txt", "cli/usage.txt", "cli/messages/<name>.txt".
//
// The `all:` prefix on the cli directory embeds everything under it, at any
// depth, so adding an asset there needs no change to this directive. A new
// asset **outside** an already-embedded directory does need one — add a
// pattern for it here, or it will not exist at runtime.
//
//go:embed version.txt
//go:embed all:cli
var Files embed.FS
