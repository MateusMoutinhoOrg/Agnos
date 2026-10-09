package utils

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/internal/stagedfs"
)

// UnitDir is one unit found in a tree — a command under CommandsDir, a route
// under RoutesDir: Name is its package, the directory's own name, and Dir the
// project-relative directory holding it.
type UnitDir struct {
	Name string
	Dir  string
}

// FindUnitDirs walks root at every depth and returns every directory holding
// marker — command.yaml, route.yaml — in the order the walk meets them. The
// marker is the whole of what makes a directory a unit: one without it is a
// folder grouping units, and the units under it are found all the same, so
// sandbox/internal/routes/backoffice/backoffice_login/route.yaml is the route login. A
// marker at root itself names no unit.
func FindUnitDirs(sandbox *api.Sandbox, io *stagedfs.StagedFS, root string, marker string) []UnitDir {
	units := []UnitDir{}
	for _, dir := range io.ListDirsRecursively(root) {
		name := LastSegment(sandbox, dir)
		if name == "" {
			continue
		}
		if _, err := io.ReadFile(dir + "/" + marker); err != nil {
			continue
		}
		units = append(units, UnitDir{Name: name, Dir: dir})
	}
	return units
}

// FindUnitDir is the directory of the unit named pkg under root, and whether
// one holds marker. Two units sharing a name is a violation verify reports;
// the first the walk meets answers until it is fixed.
func FindUnitDir(sandbox *api.Sandbox, io *stagedfs.StagedFS, root string, marker string, pkg string) (string, bool) {
	return FindUnitDirIn(FindUnitDirs(sandbox, io, root, marker), pkg)
}

// FindUnitDirIn is FindUnitDir over units already listed.
func FindUnitDirIn(units []UnitDir, pkg string) (string, bool) {
	for _, unit := range units {
		if unit.Name == pkg {
			return unit.Dir, true
		}
	}
	return "", false
}

// DuplicateUnitNames maps every name more than one unit of units carries to
// the directories carrying it. A unit's name is its Go package and the one
// every verb addresses it by, so it is unique across the whole tree, whatever
// folder it sits in.
func DuplicateUnitNames(units []UnitDir) map[string][]string {
	dirs := map[string][]string{}
	for _, unit := range units {
		dirs[unit.Name] = append(dirs[unit.Name], unit.Dir)
	}
	duplicates := map[string][]string{}
	for name, list := range dirs {
		if len(list) > 1 {
			duplicates[name] = list
		}
	}
	return duplicates
}

// CheckUniqueUnitNames refuses a tree where two units share a name: the
// generated dispatch imports each under its name, so the build would not
// compile.
func CheckUniqueUnitNames(sandbox *api.Sandbox, kind string, units []UnitDir) error {
	seen := map[string]string{}
	for _, unit := range units {
		if other, taken := seen[unit.Name]; taken {
			return sandbox.Deps.StdDeps.Errorf("%s %s is declared twice, in %s and in %s: a %s name is unique across every folder",
				kind, unit.Name, other, unit.Dir, kind)
		}
		seen[unit.Name] = unit.Dir
	}
	return nil
}

// UnitGroup normalizes a --dir typed on the command line into the folder a
// new unit lands in under its tree: each segment spelled the way a package
// directory is ("Admin Area/v1" -> "admin_area/v1"), no empty, "." or ".."
// segment, "" for the tree's top — which "/" and "." spell too.
func UnitGroup(sandbox *api.Sandbox, raw string) (string, error) {
	strs := sandbox.Deps.StringsDeps
	raw = strs.Trim(strs.TrimSpace(raw), "/")
	if raw == "" || raw == "." {
		return "", nil
	}
	segments := []string{}
	for _, segment := range strs.Split(raw, "/") {
		identifier := CommandName(sandbox, segment)
		if identifier == "" || identifier == "." || identifier == ".." {
			return "", sandbox.Deps.StdDeps.Errorf("invalid --dir %q: every folder needs a name", raw)
		}
		for _, letter := range identifier {
			valid := (letter >= 'a' && letter <= 'z') || (letter >= '0' && letter <= '9') || letter == '-'
			if !valid {
				return "", sandbox.Deps.StdDeps.Errorf("invalid --dir %q: a folder name holds only letters, digits, spaces, dashes and underscores", raw)
			}
		}
		segments = append(segments, strs.ReplaceAll(identifier, "-", "_"))
	}
	return strs.Join(segments, "/"), nil
}

// UnitDirIn is the directory a unit named pkg lands in under root, inside
// group ("" for root itself).
func UnitDirIn(root string, group string, pkg string) string {
	if group == "" {
		return root + "/" + pkg
	}
	return root + "/" + group + "/" + pkg
}

// UnitGroupOf is the folder a unit's directory sits in under root, "" when it
// sits at root itself.
func UnitGroupOf(sandbox *api.Sandbox, root string, dir string) string {
	relative := sandbox.Deps.StringsDeps.TrimPrefix(dir, root+"/")
	index := sandbox.Deps.StringsDeps.LastIndex(relative, "/")
	if index < 0 {
		return ""
	}
	return relative[:index]
}

// HoldsOtherUnit reports whether a unit directory holds another unit below
// it, so removing it whole would take that one too.
func HoldsOtherUnit(sandbox *api.Sandbox, io *stagedfs.StagedFS, dir string, marker string) bool {
	return len(FindUnitDirs(sandbox, io, dir, marker)) > 0
}

// PruneEmptyGroups removes, from the folder dir sat in up to root (never root
// itself), every folder left with nothing in it once dir is gone: a folder
// holds units, and one holding none is noise.
func PruneEmptyGroups(sandbox *api.Sandbox, io *stagedfs.StagedFS, root string, dir string) {
	strs := sandbox.Deps.StringsDeps
	for {
		index := strs.LastIndex(dir, "/")
		if index < 0 {
			return
		}
		dir = dir[:index]
		if dir == root || !strs.HasPrefix(dir, root+"/") {
			return
		}
		if len(io.ListAll(dir)) > 0 {
			return
		}
		io.RemoveDir(dir)
	}
}
