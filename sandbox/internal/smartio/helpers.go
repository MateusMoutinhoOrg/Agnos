package smartio

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/api"
)

// rootedPath joins io.Root onto a project-relative path, right before the
// value is handed to the real filesystem. It is idempotent: a path that is
// already under Root is returned unchanged, so a listing result fed back into
// another SmartIO call is never prefixed twice.
func rootedPath(sandbox *api.Sandbox, io *SmartIO, path string) string {
	if io.Root == "" {
		return path
	}
	if path == io.Root || io.sandbox.Deps.Stringsdeps.HasPrefix(path, io.Root+"/") {
		return path
	}
	return io.Root + "/" + path
}

// unrootedPath is the inverse of rootedPath: it strips io.Root back off a
// path the filesystem returned, so callers only ever see project-relative
// paths.
func unrootedPath(sandbox *api.Sandbox, io *SmartIO, path string) string {
	if io.Root == "" {
		return path
	}
	if path == io.Root {
		return ""
	}
	if io.sandbox.Deps.Stringsdeps.HasPrefix(path, io.Root+"/") {
		return path[len(io.Root)+1:]
	}
	return path
}

// unrootedPaths maps unrootedPath over a slice.
func unrootedPaths(sandbox *api.Sandbox, io *SmartIO, paths []string) []string {
	out := make([]string, 0, len(paths))
	for _, p := range paths {
		out = append(out, unrootedPath(sandbox, io, p))
	}
	return out
}

func processInputPath(sandbox *api.Sandbox, io *SmartIO, path string) (string, error) {
	p := io.Replacers.Format(path)
	if io.Ignore.IsIgnorable(p) {
		return p, sandbox.Deps.Std.Errorf("path %q is ignorable", p)
	}
	return p, nil
}

func filterIgnored(io *SmartIO, paths []string) []string {
	var result []string
	for _, p := range paths {
		if !io.Ignore.IsIgnorable(p) {
			result = append(result, p)
		}
	}
	return result
}

// isPendingRemoval checks if a path (or any of its parents) has been
// scheduled for removal in the current transaction.
func isPendingRemoval(sandbox *api.Sandbox, io *SmartIO, path string) bool {
	for _, removed := range io.PendingRemoveDirs {
		if path == removed || io.sandbox.Deps.Stringsdeps.HasPrefix(path, removed+"/") {
			return true
		}
	}
	return false
}

// isPendingCreate checks if a path has been scheduled for creation
// in the current transaction.
func isPendingCreate(io *SmartIO, path string) bool {
	for _, created := range io.PendingCreateDirs {
		if path == created {
			return true
		}
	}
	return false
}

// filterPendingRemoved removes entries that are under a pending removal directory.
func filterPendingRemoved(sandbox *api.Sandbox, io *SmartIO, paths []string) []string {
	var result []string
	for _, p := range paths {
		if !isPendingRemoval(sandbox, io, p) {
			result = append(result, p)
		}
	}
	return result
}
