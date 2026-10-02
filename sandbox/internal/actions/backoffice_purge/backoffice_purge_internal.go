package backoffice_purge

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/smartio"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/utils"
)

// BackofficePurgeInternal removes from the target project every file
// backoffice-init wrote, plus what the build generated beside them, then drops
// any directory the removal left empty and turns the mechanic off.
//
// What it removes is read off BackofficeTree itself, so it is exactly what
// init installs: a route or a command goes whole — its generated new.go and
// entries.go with it — by name, from whatever folder it was moved to; the
// backoffice's packages, database and pages go whole (BackofficeDirs); every
// other file goes alone. Nothing the project wrote is touched: RouteProps and
// api.Config lose their backoffice part with its file, and start-server never
// had one.
//
// The layers the backoffice stood on — server, front, database — are left in
// place, and so are the deps it installed: other code may use them. So is the
// store on disk, ./backofficedb: it holds the users, and removing data is not
// a purge's to do.
func BackofficePurgeInternal(sandbox *api.Sandbox, io *smartio.SmartIO, path string) error {
	sandbox.Deps.Std.Log("backoffice-purge started with path %s \n", path)

	files, err := sandbox.Deps.Embeddeps.ListFilesRecursively(utils.BackofficeTree)
	if err != nil {
		return err
	}

	doc_files, err := utils.ExtensionFiles(sandbox, utils.ExtensionSandboxBackoffice)
	if err != nil {
		return err
	}

	dirs := append([]string{}, utils.BackofficeDirs...)
	var removed []string
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		name := parts[len(parts)-1]
		unit := ""
		if len(parts) > 1 {
			unit = parts[len(parts)-2]
		}
		switch name {
		case utils.RouteConfFile:
			dirs = append(dirs, utils.RouteDir(sandbox, io, unit))
		case utils.CommandConfFile:
			dirs = append(dirs, utils.CommandDir(sandbox, io, unit))
		default:
			removed = append(removed, file)
		}
	}

	for _, file := range append(removed, doc_files...) {
		if io.IsFile(file) {
			io.RemoveDir(file)
		}
	}

	for _, dir := range dirs {
		if !io.IsDir(dir) {
			continue
		}
		for _, entry := range io.ListAllRecursively(dir) {
			io.RemoveDir(entry)
		}
		io.RemoveDir(dir)
		removed = append(removed, dir+"/")
	}

	for _, dir := range ancestorDirs(sandbox, append(removed, doc_files...)) {
		if io.IsDir(dir) && len(io.ListAll(dir)) == 0 {
			io.RemoveDir(dir)
		}
	}

	sandbox.Deps.Std.Log("backoffice-purge kept ./%s, the store holding the users: remove it by hand if you mean to\n", utils.BackofficeDatabase)

	return utils.SetExtension(sandbox, io, utils.ExtensionSandboxBackoffice, false)
}

// ancestorDirs returns every directory that contains one of the given files,
// deepest first, so an emptied child is removed before its parent is tested.
func ancestorDirs(sandbox *api.Sandbox, files []string) []string {
	seen := map[string]bool{}
	var dirs []string
	for _, file := range files {
		parts := sandbox.Deps.Stringsdeps.Split(file, "/")
		for i := 1; i < len(parts); i++ {
			dir := sandbox.Deps.Stringsdeps.Join(parts[:i], "/")
			if dir != "" && !seen[dir] {
				seen[dir] = true
				dirs = append(dirs, dir)
			}
		}
	}
	sandbox.Deps.Sortdeps.Slice(dirs, func(i int, j int) bool {
		return sandbox.Deps.Stringsdeps.Count(dirs[i], "/") > sandbox.Deps.Stringsdeps.Count(dirs[j], "/")
	})
	return dirs
}
