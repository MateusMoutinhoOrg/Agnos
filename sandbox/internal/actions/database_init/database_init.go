package database_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	buildAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/build"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
)

// KeepModule is the agnos repo that stores the records, installed as a remote
// dep rather than from the catalog: it is a repo, not a contract agnos carries.
// The version is pinned here so every project this init touches installs the
// same one, and `agnos set-dep database --version <v>` moves it afterwards.
const KeepModule = "github.com/MateusMoutinhoOrg/Keep@v0.0.7"

// DatabaseDep is the name the copied contract lands under. The generated code
// names it, so it is the one spelling database-init may install it as.
const DatabaseDep = "database"

// databaseDeps are the contracts the database layer calls into beyond the
// store itself: the output channels, and the text conversion the generated
// filtrage is written through.
var databaseDeps = []string{"std", "stringsdeps"}

// InstallDeps installs the store and the contracts the generated code calls
// into. It is exported for the same reason server-init exports its own: the
// internal half renders assets and writes nothing to go.mod, so a layer that
// composes DatabaseInitInternal into its own transaction still has to install
// this set first.
func InstallDeps(sandbox *api.Sandbox, path string) error {
	for _, dep := range databaseDeps {
		if err := addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}

	return addDepAction.AddDep(sandbox, api.AddDepProps{Path: path, Dep: KeepModule, As: DatabaseDep})
}

// DatabaseInit installs the store the database layer is built over and turns
// the database mechanic on, then runs build as a follow-up step, which renders
// the group.
func DatabaseInit(sandbox *api.Sandbox, path string) error {
	if err := InstallDeps(sandbox, path); err != nil {
		return err
	}

	io := smartio.New(sandbox, path, sandbox.Config.ProjectName)
	if err := DatabaseInitInternal(sandbox, io, path); err != nil {
		return err
	}
	if err := io.Persist(); err != nil {
		return err
	}
	return buildAction.Build(sandbox, api.BuildProps{Path: path, Runtime: api.RuntimeGo})
}
