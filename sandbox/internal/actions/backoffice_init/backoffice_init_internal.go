package backoffice_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// backofficeDeps are the catalog contracts the backoffice calls into beyond
// the ones its layers install. serverdeps is among them although server-init
// installs it: installing it again renders the catalog's current copy, the
// one answering GetClientIp, over an older one. archivedeps and iodeps are the
// backups': the zip a snapshot travels as, and the --database folder a
// snapshot is read from and a restore writes to.
var backofficeDeps = []string{
	"archivedeps",
	"embeddeps",
	"envdeps",
	"hashdeps",
	"iodeps",
	"jwtdeps",
	"passworddeps",
	"randdeps",
	"ratelimitdeps",
	"serializabledeps",
	"serverdeps",
	"sortdeps",
	"stddeps",
	"stringsdeps",
	"timedeps",
}

// gitignoreFile is the project's .gitignore, where the backoffice's store is
// kept out of version control.
const gitignoreFile = ".gitignore"

// BackofficeInitInternal turns on, on this same open StagedFS, every layer the
// backoffice stands on that is off — server, front, database — installs the
// catalog deps it calls into, and writes BackofficeTree into the project.
//
// Every file it writes is the project's from that moment: a second
// backoffice-init keeps each one already there, so an edit to a page or a
// handler is never undone. Nothing the project wrote is edited either — the
// backoffice's part of RouteProps and of api.Config are files of their own
// (routeprops/backoffice.go, api/backofficeconfig.go) the generated
// aggregates embed, and its reading of the secret is a middleware in front of
// start-server rather than a change to start-server's handler.
func BackofficeInitInternal(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	sandbox.Deps.StdDeps.Logf("backoffice-init started with path %s \n", path)

	if err := installLayers(sandbox, io, path); err != nil {
		return err
	}

	for _, dep := range backofficeDeps {
		if err := addDepAction.AddDepInternal(sandbox, io, api.AddDepProps{Path: path, Dep: dep}); err != nil {
			return err
		}
	}

	if err := refuseNameClashes(sandbox, io); err != nil {
		return err
	}

	module_conf, err := utils.LoadModuleConf(sandbox, io)
	if err != nil {
		return err
	}
	project_conf, err := utils.LoadProjectConf(sandbox, io)
	if err != nil {
		return err
	}

	vars := map[string]interface{}{
		"Module":        module_conf.Module,
		"ProjectName":   project_conf.ProjectName,
		"GeneratorName": sandbox.Deps.StringsDeps.ToLower(sandbox.Config.ProjectName),
		"SecretEnv":     utils.SecretEnvName(sandbox, project_conf.ProjectName),
	}
	if _, err := utils.RenderTemplateTree(sandbox, io, utils.BackofficeTree, utils.BackofficeRawDir, vars); err != nil {
		return err
	}

	if err := ignoreStore(sandbox, io); err != nil {
		return err
	}

	sandbox.Deps.StdDeps.Logf("backoffice-init: add the first user with `%s add-backoffice-user --role root`, then start the server; set %s to a random secret of at least 32 characters (openssl rand -hex 32) so sessions survive a restart, or one is generated for each run\n", project_conf.ProjectName, vars["SecretEnv"])

	return utils.SetExtension(sandbox, io, utils.ExtensionBackoffice, true)
}

// installLayers turns on each layer the backoffice requires that is off, in
// the order they require each other: front-init brings a server along itself,
// so the server is turned on here only to be explicit about it.
func installLayers(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error {
	layers := []struct {
		name string
		init func(sandbox *api.Sandbox, io *stagedfs.StagedFS, path string) error
	}{
		{utils.ExtensionServer, serverInitAction.ServerInitInternal},
		{utils.ExtensionFront, frontInitAction.FrontInitInternal},
		{utils.ExtensionDatabase, databaseInitAction.DatabaseInitInternal},
	}

	for _, layer := range layers {
		enabled, err := utils.ExtensionEnabled(sandbox, io, layer.name)
		if err != nil {
			return err
		}
		if enabled {
			continue
		}
		if err := layer.init(sandbox, io, path); err != nil {
			return err
		}
	}
	return nil
}

// refuseNameClashes stops before anything is written when a route or a
// command of BackofficeTree is already declared somewhere else in the
// project: a route's and a command's name is unique across folders, so the
// project's would shadow the backoffice's, or the reverse.
func refuseNameClashes(sandbox *api.Sandbox, io *stagedfs.StagedFS) error {
	files, err := sandbox.Deps.EmbedDeps.ListFilesRecursively(utils.BackofficeTree)
	if err != nil {
		return err
	}

	for _, file := range files {
		dir, name, unit := unitOf(sandbox, file)
		if unit == "" {
			continue
		}
		found := utils.RouteDir(sandbox, io, name)
		if unit == utils.CommandConfFile {
			found = utils.CommandDir(sandbox, io, name)
		}
		if found != dir && io.IsFile(found+"/"+unit) {
			return sandbox.Deps.StdDeps.Errorf("%s is already declared by %s; the backoffice declares its own at %s — rename the project's first", name, found, dir)
		}
	}
	return nil
}

// unitOf reads one path of BackofficeTree as a unit declaration: the
// directory, the unit's name and the declaration file when it is a route.yaml
// or a command.yaml, and "" as the file otherwise.
func unitOf(sandbox *api.Sandbox, file string) (string, string, string) {
	parts := sandbox.Deps.StringsDeps.Split(file, "/")
	last := parts[len(parts)-1]
	if last != utils.RouteConfFile && last != utils.CommandConfFile {
		return "", "", ""
	}
	return sandbox.Deps.StringsDeps.Join(parts[:len(parts)-1], "/"), parts[len(parts)-2], last
}

// ignoreStore appends the backoffice's stores — the directories its
// databases write to while --database is left at its default — to
// .gitignore, once each. The users, their password hashes and every snapshot
// of them live there, so neither is ever committed.
func ignoreStore(sandbox *api.Sandbox, io *stagedfs.StagedFS) error {
	content, err := io.ReadFile(gitignoreFile)
	if err != nil {
		content = nil
	}
	text := string(content)

	for _, store := range []string{utils.BackofficeStore, utils.BackupStore} {
		entry := "/" + store
		if ignored(sandbox, text, entry) {
			continue
		}
		if text != "" && !sandbox.Deps.StringsDeps.HasSuffix(text, "\n") {
			text += "\n"
		}
		text += entry + "\n"
	}

	if text == string(content) {
		return nil
	}
	return io.WriteFile(gitignoreFile, []byte(text))
}

// ignored tells whether text, a .gitignore, already holds entry as a line.
func ignored(sandbox *api.Sandbox, text string, entry string) bool {
	for _, line := range sandbox.Deps.StringsDeps.Split(text, "\n") {
		if sandbox.Deps.StringsDeps.TrimSpace(line) == entry {
			return true
		}
	}
	return false
}
