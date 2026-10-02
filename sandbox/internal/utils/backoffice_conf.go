package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// BackofficeTree is the asset tree backoffice-init writes from: every file
// under it lands at the path it holds inside it, once — routes, commands, the
// database declaration, the hand-written packages, the pages. None of it is
// rendered again by a build, so it is not an asset group (see
// RenderTemplateTree).
const BackofficeTree = "templates/backoffice"

// BackofficeRawDir is the part of BackofficeTree copied byte for byte: the
// html pages are text/template sources the project renders at runtime, so
// rendering them at install would eat their actions.
const BackofficeRawDir = "assets"

// BackofficeDatabase is the database backoffice-init declares. It is the
// backoffice's alone, so backoffice-purge removes it whole and the project's
// own databases are never touched.
const BackofficeDatabase = "backofficedb"

// BackofficeDirs are the directories the backoffice owns whole: its packages,
// its database package and its pages. Everything else it wrote is either a
// route or a command — removed by name, wherever it was moved — or one file.
var BackofficeDirs = []string{
	"sandbox/internal/server/backoffice",
	DatabasesDir + "/" + BackofficeDatabase,
	"assets/backoffice",
}

// BackofficeUnits is every unit backoffice-init declares — each route and
// command of BackofficeTree, and BackofficeDatabase — under both spellings a
// unit is named by (add_backoffice_user, add-backoffice-user). A count of what
// the project declared itself leaves them out, the way it leaves out what any
// other init scaffolds.
func BackofficeUnits(sandbox *api.Sandbox) map[string]bool {
	units := map[string]bool{BackofficeDatabase: true}

	files, err := sandbox.Deps.Embeddeps.ListFilesRecursively(BackofficeTree)
	if err != nil {
		return units
	}
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		last := parts[len(parts)-1]
		if len(parts) < 2 || (last != RouteConfFile && last != CommandConfFile) {
			continue
		}
		name := parts[len(parts)-2]
		units[name] = true
		units[sandbox.Deps.Stringsdeps.ReplaceAll(name, "_", "-")] = true
	}
	return units
}
