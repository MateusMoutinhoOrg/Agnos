package assets

// The project's embedded assets: the files the library serves through the
// injected embed contract (deps.Deps.EmbedDeps). An asset is a payload better
// kept as a file than as a Go constant — a template, a long-form document, an
// image — reached by path at runtime and shipped inside the binary.
//
// The tree is empty in this repository, and deliberately so. The tracker's
// display text is short and fixed, so it lives in sandbox/config as
// compile-time constants, where the compiler checks every reference. What
// remains here is the mechanic itself, wired end to end and ready for a
// library derived from this template to drop files into. Adding one is
// docs/Tutorials/HandleAssets.md; nothing about the directive below has to
// change when you do.
//
// This package exists for one reason: a //go:embed directive can only reach
// files inside its own package directory, so the directive has to live next
// to the assets themselves. It holds no behavior and no state beyond the
// embedded filesystem, and only code outside the sandbox may import it — the
// standard adapter does, in adapters/standard/embed.go, and wraps it into the
// embeddeps.Lib contract the sandbox reads through.

import "embed"

// Files is every asset shipped with the project, compiled into the binary, so
// an installed `agnos-cli` carries them with no files on disk next to it.
// Paths inside it are slash-separated and rooted here: an asset written to
// assets/templates/report.tmpl is read as "templates/report.tmpl".
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
