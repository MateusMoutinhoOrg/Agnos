package backoffice_init

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	addDepAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/add_dep"
	databaseInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/database_init"
	frontInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/front_init"
	serverInitAction "github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/actions/server_init"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// backofficeDeps are the catalog contracts the backoffice calls into beyond
// the ones its layers install. serverdeps is among them although server-init
// installs it: installing it again renders the catalog's current copy, the
// one answering GetClientIp, over an older one.
var backofficeDeps = []string{
	"embeddeps",
	"envdeps",
	"hashdeps",
	"jwtdeps",
	"passworddeps",
	"randdeps",
	"ratelimitdeps",
	"serializables",
	"serverdeps",
	"sortdeps",
	"std",
	"stringsdeps",
	"timedeps",
}

// gitignoreFile is the project's .gitignore, where the backoffice's store is
// kept out of version control.
const gitignoreFile = ".gitignore"

// BackofficeInitInternal turns on, on this same open SmartIO, every layer the
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
func BackofficeInitInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("backoffice-init started with path %s \n", path)

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
		"Name":          project_conf.Name,
		"GeneratorName": sandbox.Deps.Stringsdeps.ToLower(sandbox.Config.ProjectName),
		"SecretEnv":     utils.SecretEnvName(sandbox, project_conf.Name),
	}
	if _, err := utils.RenderTemplateTree(sandbox, io, utils.BackofficeTree, utils.BackofficeRawDir, vars); err != nil {
		return err
	}

	if err := ignoreStore(sandbox, io); err != nil {
		return err
	}

	sandbox.Deps.Std.Log("backoffice-init: start the server with %s set to a random secret of at least 32 characters (openssl rand -hex 32), then add the first user with `%s add-backoffice-user --role root`\n", vars["SecretEnv"], project_conf.Name)

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxBackoffice, true)
}

// installLayers turns on each layer the backoffice requires that is off, in
// the order they require each other: front-init brings a server along itself,
// so the server is turned on here only to be explicit about it.
func installLayers(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	layers := []struct {
		name string
		init func(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error
	}{
		{utils.ExtensionSandboxServer, serverInitAction.ServerInitInternal},
		{utils.ExtensionSandboxFront, frontInitAction.FrontInitInternal},
		{utils.ExtensionSandboxDatabase, databaseInitAction.DatabaseInitInternal},
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
func refuseNameClashes(sandbox *api.Sandbox, io *smartio.SmartIO) error {
	files, err := sandbox.Deps.Embeddeps.ListFilesRecursively(utils.BackofficeTree)
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
			return sandbox.Deps.Std.Errorf("%s is already declared by %s; the backoffice declares its own at %s — rename the project's first", name, found, dir)
		}
	}
	return nil
}

// unitOf reads one path of BackofficeTree as a unit declaration: the
// directory, the unit's name and the declaration file when it is a route.yaml
// or a command.yaml, and "" as the file otherwise.
func unitOf(sandbox *api.Sandbox, file string) (string, string, string) {
	parts := sandbox.Deps.Stringsdeps.Split(file, "/")
	last := parts[len(parts)-1]
	if last != utils.RouteConfFile && last != utils.CommandConfFile {
		return "", "", ""
	}
	return sandbox.Deps.Stringsdeps.Join(parts[:len(parts)-1], "/"), parts[len(parts)-2], last
}

// ignoreStore appends the backoffice's store — the directory its database
// writes to, in the directory the server runs from — to .gitignore, once. The
// users and their password hashes live there, so it is never committed.
func ignoreStore(sandbox *api.Sandbox, io *smartio.SmartIO) error {
	entry := "/" + utils.BackofficeDatabase

	content, err := io.ReadFile(gitignoreFile)
	if err != nil {
		content = nil
	}
	for _, line := range sandbox.Deps.Stringsdeps.Split(string(content), "\n") {
		if sandbox.Deps.Stringsdeps.TrimSpace(line) == entry {
			return nil
		}
	}

	text := string(content)
	if text != "" && !sandbox.Deps.Stringsdeps.HasSuffix(text, "\n") {
		text += "\n"
	}
	return io.WriteFileOverwrite(gitignoreFile, []byte(text+entry+"\n"))
}
