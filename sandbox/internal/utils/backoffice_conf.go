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
const BackofficeDatabase = "backoffice_db"

// BackofficeStore is the directory the backoffice's database writes to, under
// the directory the server runs from: its key-prefix, kept apart from the
// package so a search for one never lands in the other.
const BackofficeStore = "data/backofficedb"

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
	units := map[string]bool{BackofficeDatabase: true, sandbox.Deps.StringsDeps.ReplaceAll(BackofficeDatabase, "_", "-"): true}

	files, err := sandbox.Deps.EmbedDeps.ListFilesRecursively(BackofficeTree)
	if err != nil {
		return units
	}
	for _, file := range files {
		parts := sandbox.Deps.StringsDeps.Split(file, "/")
		last := parts[len(parts)-1]
		if len(parts) < 2 || (last != RouteConfFile && last != CommandConfFile) {
			continue
		}
		name := parts[len(parts)-2]
		units[name] = true
		units[sandbox.Deps.StringsDeps.ReplaceAll(name, "_", "-")] = true
	}
	return units
}

// SecretEnvName is the environment variable a project named name reads its
// backoffice session secret from: the name upper-cased, every byte but a
// letter or a digit turned into "_", then "_BACKOFFICE_SECRET" —
// MEUSITE_BACKOFFICE_SECRET for meusite. The backoffice's SecretEnv spells the
// same rule at runtime from api.Config.ProjectName; this is the copy the docs
// render from.
func SecretEnvName(sandbox *api.Sandbox, name string) string {
	upper := []byte(sandbox.Deps.StringsDeps.ToUpper(name))
	for i, char := range upper {
		if !(char >= 'A' && char <= 'Z') && !(char >= '0' && char <= '9') {
			upper[i] = '_'
		}
	}
	return string(upper) + "_BACKOFFICE_SECRET"
}
