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
// "version.txt", "usages.txt", "messages/<name>.txt".
//
// The single pattern below takes the whole directory: `*` matches every entry
// next to this file, the `all:` prefix descends into each directory it
// matches and keeps the names other patterns would skip. Adding an asset
// anywhere under /assets/ therefore needs no change here — put the file in
// the tree and it exists at runtime.
//
// This file matches its own pattern and is embedded along with the assets.
// That costs a few hundred bytes in the binary and puts "asset.go" in a
// recursive listing of the asset root; it is the price of a directive that
// never has to be edited.
//
//go:embed all:*
var Files embed.FS
