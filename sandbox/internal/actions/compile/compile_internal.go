package compile

import (
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps"
	"github.com/MateusMoutinhoOrg/Agnos/sandbox/deps/rundeps"
)

// target is one cross-compilation target: the GOOS/GOARCH pair the Go
// toolchain is told to build for and the release/ filename its binary is
// written to.
type target struct {
	GOOS   string
	GOARCH string
	Output string
}

// targets maps every --target name `agnos compile` accepts onto its
// GOOS/GOARCH pair and release/ output filename. "all" is not a key here: it
// expands to every entry, in targetOrder. Adding a row to the table in
// docs/Tutorials/Build.md means adding it here too.
var targets = map[string]target{
	"linux86":    {GOOS: "linux", GOARCH: "amd64", Output: "linux86.out"},
	"linuxarm64": {GOOS: "linux", GOARCH: "arm64", Output: "linuxarm64.out"},
	"linuxi32":   {GOOS: "linux", GOARCH: "386", Output: "linuxi32.out"},
	"mac86":      {GOOS: "darwin", GOARCH: "amd64", Output: "mac86.bin"},
	"macarm64":   {GOOS: "darwin", GOARCH: "arm64", Output: "macarm64.bin"},
	"windows86":  {GOOS: "windows", GOARCH: "amd64", Output: "windows86.exe"},
	"windowsi32": {GOOS: "windows", GOARCH: "386", Output: "windowsi32.exe"},
}

// targetOrder is the deterministic order the targets are built in, and the
// order "all" expands to.
var targetOrder = []string{
	"linux86", "linuxarm64", "linuxi32", "mac86", "macarm64", "windows86", "windowsi32",
}

// CompileInternal cross-compiles ./cmd/main once per resolved target name into
// release/<file>, with CGO disabled so no C toolchain is needed. `go build -o`
// creates the release/ directory itself. The names must already be validated
// (see resolveTargets).
func CompileInternal(deps *deps.Deps, path string, names []string) error {
	for _, name := range names {
		t := targets[name]
		output := "release/" + t.Output

		deps.Std.Log("compile %s: go build -o %s ./cmd/main (GOOS=%s GOARCH=%s) \n", name, output, t.GOOS, t.GOARCH)

		result, err := deps.Rundeps.Run(rundeps.RunProps{
			Dir:     path,
			Program: "go",
			Args:    []string{"build", "-o", output, "./cmd/main"},
			Env:     []string{"CGO_ENABLED=0", "GOOS=" + t.GOOS, "GOARCH=" + t.GOARCH},
		})
		if err != nil {
			return deps.Std.Errorf("compile %s: could not run go build: %w", name, err)
		}
		if result.ExitCode != 0 {
			return deps.Std.Errorf("compile %s: go build failed:\n%s", name, result.Output)
		}
	}
	return nil
}

// resolveTargets turns the raw --target values into the ordered list of target
// names to build. "all" (in any position) expands to every target; otherwise
// each value must be a known target name, duplicates are dropped, and the
// written order is kept. An empty request is a usage error.
func resolveTargets(deps *deps.Deps, requested []string) ([]string, error) {
	if len(requested) == 0 {
		return nil, deps.Std.Errorf("compile: at least one --target is required (or --target all)")
	}

	for _, name := range requested {
		if name == "all" {
			return targetOrder, nil
		}
	}

	seen := map[string]bool{}
	resolved := []string{}
	for _, name := range requested {
		if _, ok := targets[name]; !ok {
			return nil, deps.Std.Errorf("compile: unknown target %q (accepted: linux86, linuxarm64, linuxi32, mac86, macarm64, windows86, windowsi32, all)", name)
		}
		if seen[name] {
			continue
		}
		seen[name] = true
		resolved = append(resolved, name)
	}
	return resolved, nil
}
